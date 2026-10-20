package resolvers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
	"github.com/mia0x75/halo/internal/worker"
)

// CancelCron 取消一个预约（工单关闭）
func (r *mutationRootResolver) CancelCron(ctx context.Context, id string) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		s := worker.NewScheduler()
		if err = s.Cancel(id); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		domainCron, cronErr := r.app.Services.Cron.GetByUUID(ctx, id)
		if cronErr != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, cronErr.Error())
			break
		}
		cron := bridge.CronFromDomain(domainCron)

		domainTicket, tkErr := r.app.Services.Ticket.GetByCronID(ctx, cron.CronID)
		if tkErr != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, tkErr.Error())
			break
		}
		ticket := bridge.TicketFromDomain(domainTicket)

		event.Fire(event.EventCronCancelled, &event.CronCancelled{
			User:   *credential.User,
			Cron:   *cron,
			Ticket: *ticket,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// Cron 查看某一个预约的详细信息
func (r *queryRootResolver) Cron(ctx context.Context, id string) (cron *models.Cron, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		domainCron, getErr := r.app.Services.Cron.GetByUUID(ctx, id)
		if getErr != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
			break
		}
		cron = bridge.CronFromDomain(domainCron)

		break
	}

	if err != nil {
		cron = nil
	}

	return
}

// Crons 分页查看预约任务列表
func (r *queryRootResolver) Crons(ctx context.Context, after *string, before *string, first *int, last *int) (*gqlapi.CronConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	from := math.MaxInt64
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return nil, err
		}
		from = i
	}
	hasPreviousPage := true
	hasNextPage := true

	if from == math.MaxInt64 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []*gqlapi.CronEdge{}

	// 列表不允许查询Content内容
	domainCrons, err := r.app.Services.Cron.List(ctx, int64(from), *first)
	if err != nil {
		return nil, err
	}
	crons := make([]*models.Cron, 0, len(domainCrons))
	for _, c := range domainCrons {
		crons = append(crons, bridge.CronFromDomain(c))
	}

	for _, cron := range crons {
		edges = append(edges, &gqlapi.CronEdge{
			Node:   cron,
			Cursor: EncodeCursor(fmt.Sprintf("%d", cron.CronID)),
		})
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return nil, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", crons[0].CronID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", crons[len(crons)-1].CronID))
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.CronConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// Tasks 返回所有正在等待执行的计划任务
func (r *queryRootResolver) Tasks(ctx context.Context) (L []*models.Cron, err error) {
	s := worker.NewScheduler()
	return s.Crons()
}

// cronJob is one schedulable system job. The command is fixed by the server:
// the API never accepts a free-form command, so scheduling a job cannot turn
// into remote command execution.
type cronJob struct {
	cmd    string
	params []string
}

// cronJobs maps the public job types to their fixed commands.
var cronJobs = map[gqlapi.CronJobEnum]cronJob{
	gqlapi.CronJobEnumMetricsRollup: {
		cmd:    "halocli",
		params: []string{"metrics", "rollup"},
	},
	gqlapi.CronJobEnumMetricsCleanup: {
		cmd:    "halocli",
		params: []string{"metrics", "cleanup", "--days=90"},
	},
}

// resolveCronJob returns the fixed command of a job type. The days argument
// only applies to the cleanup job.
func resolveCronJob(jobType gqlapi.CronJobEnum, days *int) (cronJob, error) {
	job, ok := cronJobs[jobType]
	if !ok {
		return cronJob{}, fmt.Errorf("错误代码: %s, 错误信息: 不支持的作业类型(job=%s)。", gqlapi.ReturnCodeInvalidParams, jobType)
	}

	params := append([]string{}, job.params...)
	if days != nil && *days > 0 {
		// 清理作业的保留天数由调用方决定。
		params[len(params)-1] = fmt.Sprintf("--days=%d", *days)
	}
	job.params = params
	return job, nil
}

// CreateCron 创建一个系统定时作业
func (r *mutationRootResolver) CreateCron(ctx context.Context, input gqlapi.CreateCronInput) (*models.Cron, error) {
	job, err := resolveCronJob(input.Job, input.Days)
	if err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(input.Interval)
	if err != nil || interval <= 0 {
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 执行间隔(interval=%s)无效。", gqlapi.ReturnCodeInvalidParams, input.Interval)
	}

	encoded, err := json.Marshal(job.params)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cron := &models.Cron{
		UUID:      uuid.NewString(),
		Status:    string(domain.CronStatusPending),
		Name:      input.Name,
		Cmd:       job.cmd,
		Params:    string(encoded),
		Interval:  interval.String(),
		Duration:  "0s",
		LastRun:   now.UTC().Format(time.RFC3339),
		NextRun:   now.Add(interval).UTC().Format(time.RFC3339),
		Recurrent: 1,
		Version:   1,
	}
	cron.BeforeInsert()

	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	if _, err := r.app.Services.Cron.Create(ctx, cronToDomain(cron), credential.User.UserID); err != nil {
		return nil, err
	}

	// 让调度器立刻装载新作业，不必等下一次周期装载。
	worker.Reload()
	return cron, nil
}

// RemoveCron 删除一个系统定时作业
func (r *mutationRootResolver) RemoveCron(ctx context.Context, id string) (*bool, error) {
	cron, err := r.app.Services.Cron.GetByUUID(ctx, id)
	if err != nil {
		return nil, err
	}

	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	if err := r.app.Services.Cron.Delete(ctx, cron.CronID, credential.User.UserID); err != nil {
		return nil, err
	}

	worker.Reload()
	ok := true
	return &ok, nil
}

// cronToDomain converts a cron model into its domain entity.
func cronToDomain(cron *models.Cron) *domain.Cron {
	return &domain.Cron{
		CronID:    cron.CronID,
		UUID:      cron.UUID,
		Status:    domain.CronStatus(cron.Status),
		Name:      cron.Name,
		Cmd:       cron.Cmd,
		Params:    cron.Params,
		Interval:  &cron.Interval,
		Duration:  &cron.Duration,
		LastRun:   &cron.LastRun,
		NextRun:   &cron.NextRun,
		Recurrent: &cron.Recurrent,
		Version:   cron.Version,
		UpdateAt:  &cron.UpdateAt,
		CreateAt:  cron.CreateAt,
	}
}
