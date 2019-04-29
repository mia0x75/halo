package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mia0x75/halo/crons"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// CancelCron 取消一个预约（工单关闭）
func (r *mutationRootResolver) CancelCron(ctx context.Context, id string) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		s := crons.NewScheduler()
		if err = s.Cancel(id); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		cron := &models.Cron{}
		ticket := &models.Ticket{}
		if _, err = g.Engine.Where("`uuid` = ?", id).Get(cron); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		if _, err = g.Engine.Where("`cron_id` = ?", cron.CronID).Get(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		events.Fire(events.EventCronCancelled, &events.CronCancelledArgs{
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
		cron = &models.Cron{
			UUID: id,
		}
		if err = g.Engine.Find(cron); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

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
	edges := []gqlapi.CronEdge{}
	crons := []*models.Cron{}

	var err error

	// 列表不允许查询Content内容
	if err = g.Engine.Desc("cron_id").Where("cron_id < ?", from).Limit(*first).Find(&crons); err != nil {
		return nil, err
	}

	for _, cron := range crons {
		edges = append(edges, gqlapi.CronEdge{
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
	pageInfo := gqlapi.PageInfo{
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
	s := crons.NewScheduler()
	return s.Crons()
}

type cronResolver struct{ *Resolver }
