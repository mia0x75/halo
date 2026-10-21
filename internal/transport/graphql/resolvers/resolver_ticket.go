package resolvers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/driver"
	"github.com/mia0x75/parser/format"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/crypto"
	"github.com/mia0x75/halo/internal/service"
	"github.com/mia0x75/halo/internal/tools"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
	"github.com/mia0x75/halo/internal/validate"
	"github.com/mia0x75/halo/internal/worker"
)

// 保留，不可以删除
var _ = driver.ValueExpr{}

// CreateTicket 创建一个工单
func (r *mutationRootResolver) CreateTicket(ctx context.Context, input models.CreateTicketInput) (ticket *models.Ticket, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		statements := []*models.Statement{}
		user := credential.User

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		// 查询缓存
		cluster := cache.ClustersMap.Any(func(elem *models.Cluster) bool {
			return elem.UUID == input.ClusterUUID
		})

		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, input.ClusterUUID)
			break
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, input.ClusterUUID)
			break
		}

		// 检查群集访问权：直接授权，或角色被授予该群集所属分组。
		if !transportauth.CanAccessCluster(user, credential.Roles, cluster) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break
		}

		password, _ := crypto.DecryptAES(cluster.Password, config.Config().Secret.Crypto)
		target := bridge.ClusterToDomain(cluster)
		target.Database = input.Database

		if err = statDatabase(ctx, r.app.Agent, target, string(password), input.Database); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 查询缓存
		reviewer := cache.UsersMap.Any(func(elem *models.User) bool {
			return elem.UUID == input.ReviewerUUID
		})
		if reviewer == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 审核用户(uuid=%s)不存在。", rc, input.ReviewerUUID)
			break
		}

		if reviewer.Status != gqlapi.UserStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeUserNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 审核用户(uuid=%s)状态异常。", rc, input.ReviewerUUID)
			break
		}

		// 检查审核人关联
		if !cache.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == reviewer.UserID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeNoEdgeToReviewer
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联审核用户(uuid=%s)。", rc, user.UUID, reviewer.UUID)
			break
		}

		// TODO: 暂时保留
		defer func() {
			if err := recover(); err != nil {
				rc = gqlapi.ReturnCodeUnknowError
				log.Errorf("[E] 错误代码: %s, 错误信息: %s\n", rc, tools.PanicDetail())
			}
		}()

		// 拆分语句
		p := parser.New()
		var stmts []ast.StmtNode
		stmts, _, err = p.Parse(input.Content, "", "")
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// 创建工单
		ticket = &models.Ticket{
			Subject:    input.Subject,
			Content:    input.Content,
			UserID:     user.UserID,
			ClusterID:  cluster.ClusterID,
			ReviewerID: reviewer.UserID,
			Database:   input.Database,
			Status:     gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForVld],
		}

		i := 1
		for _, node := range stmts {
			var sb strings.Builder
			_ = node.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			sql := sb.String()
			if len(sql) == 0 {
				sql = strings.TrimSpace(node.Text())
			}
			stat := models.Statement{
				Sequence:   uint16(i),
				Content:    sql,
				Type:       StatementType2Uint8(node),
				Status:     gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForVld],
				StmtNode:   node,
				Violations: &models.Violations{},
			}
			statements = append(statements, &stat)
			i++
		}

		domainTicket := bridge.TicketToDomain(ticket)
		domainStmts := make([]*domain.Statement, 0, len(statements))
		for _, st := range statements {
			domainStmts = append(domainStmts, bridge.StatementToDomain(st))
		}
		newID, createErr := r.app.Services.Ticket.Create(ctx, domainTicket, domainStmts)
		if createErr != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, createErr.Error())
			break
		}
		ticket.TicketID = uint(newID)
		ticket.UUID = domainTicket.UUID
		for i, st := range statements {
			st.TicketID = uint(newID)
			// The service assigned each persisted statement its uuid; carry it
			// back so the async validation can update the rows by uuid.
			st.UUID = domainStmts[i].UUID
		}

		notifyTicketStatusChange(user.UUID, ticket.UUID, fmt.Sprintf("新建工单(uuid=%s)成功，系统正在开始启动自动化审核。", ticket.UUID))
		go r.validation(statements, cluster, ticket)

		event.Fire(event.EventTicketCreated, &event.TicketCreated{
			User:    *user,
			Ticket:  *ticket,
			Cluster: *cluster,
		})

		// 退出for循环
		break
	}

	if err != nil {
		ticket = nil
	}

	return
}

// UpdateTicket 修改一个工单
func (r *mutationRootResolver) UpdateTicket(ctx context.Context, input models.UpdateTicketInput) (ticket *models.Ticket, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		statements := []*models.Statement{}
		user := credential.User
		cluster := cache.ClustersMap.Any(func(elem *models.Cluster) bool {
			return elem.UUID == input.ClusterUUID
		})

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, input.ClusterUUID)
			break
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, input.ClusterUUID)
			break
		}

		// 检查群集访问权：直接授权，或角色被授予该群集所属分组。
		if !transportauth.CanAccessCluster(user, credential.Roles, cluster) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break
		}

		password, _ := crypto.DecryptAES(cluster.Password, config.Config().Secret.Crypto)
		target := bridge.ClusterToDomain(cluster)
		target.Database = input.Database

		if err = statDatabase(ctx, r.app.Agent, target, string(password), input.Database); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 查询缓存
		reviewer := cache.UsersMap.Any(func(elem *models.User) bool {
			return elem.UUID == input.ReviewerUUID
		})

		if reviewer == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 审核用户(uuid=%s)不存在。", rc, input.ReviewerUUID)
			break
		}

		if reviewer.Status != gqlapi.UserStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeUserNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 审核用户(uuid=%s)的当前状态异常。", rc, input.ReviewerUUID)
			break
		}

		// 检查审核人关联
		if !cache.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == reviewer.UserID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeNoEdgeToReviewer
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联审核用户(uuid=%s)。", rc, user.UUID, reviewer.UUID)
			break
		}

		// TODO: 暂时保留
		defer func() {
			if err := recover(); err != nil {
				rc = gqlapi.ReturnCodeUnknowError
				log.Errorf("[E] 错误代码: %s, 错误信息: %s\n", rc, tools.PanicDetail())
			}
		}()

		// 拆分语句
		i := 1
		p := parser.New()
		var stmts []ast.StmtNode
		stmts, _, err = p.Parse(input.Content, "", "")
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		//获取数据库原有的ticket信息
		existing, getErr := r.app.Services.Ticket.GetByUUID(ctx, input.TicketUUID)
		if getErr != nil {
			if errors.Is(getErr, domain.ErrNotFound) {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			} else {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
			}
			break
		}
		ticket = bridge.TicketFromDomain(existing)

		// 状态已关闭、已执行成功、已执行失败的工单不可编辑
		if ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumClosed] ||
			ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumDone] ||
			ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumExecFailure] {
			// TODO： 处理rc
			err = fmt.Errorf("错误代码: %s, 错误信息: 关闭、已执行成功和执行失败的工单不可编辑。", rc)
			break
		}

		// 更新工单
		ticket.Subject = input.Subject
		ticket.Content = input.Content
		ticket.UserID = user.UserID
		ticket.ClusterID = cluster.ClusterID
		ticket.ReviewerID = reviewer.UserID
		ticket.Database = input.Database
		ticket.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForVld]

		for _, node := range stmts {
			var sb strings.Builder
			_ = node.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			sql := sb.String()
			if len(sql) == 0 {
				sql = strings.TrimSpace(node.Text())
			}
			stat := models.Statement{
				Sequence:   uint16(i),
				Type:       StatementType2Uint8(node),
				Content:    sql,
				Status:     gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForVld],
				StmtNode:   node,
				Violations: &models.Violations{},
			}
			statements = append(statements, &stat)
			i++
		}

		domainTicket := bridge.TicketToDomain(ticket)
		domainStmts := make([]*domain.Statement, 0, len(statements))
		for _, st := range statements {
			domainStmts = append(domainStmts, bridge.StatementToDomain(st))
		}
		if updErr := r.app.Services.Ticket.Update(ctx, domainTicket, domainStmts); updErr != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, updErr.Error())
			break
		}
		for i, st := range statements {
			st.TicketID = ticket.TicketID
			// Update re-creates the statements, so pick up the freshly assigned
			// uuids for the async validation writeback.
			st.UUID = domainStmts[i].UUID
		}

		notifyTicketStatusChange(user.UUID, ticket.UUID, fmt.Sprintf("更新工单(uuid=%s)成功，系统正在开始启动自动化审核。", ticket.UUID))
		go r.validation(statements, cluster, ticket)

		event.Fire(event.EventTicketUpdated, &event.TicketUpdated{
			User:    *user,
			Ticket:  *ticket,
			Cluster: *cluster,
		})

		// 退出for循环
		break
	}

	if err != nil {
		ticket = nil
	}

	return
}

// RemoveTicket 有条件的删除一个工单
func (r *mutationRootResolver) RemoveTicket(ctx context.Context, id string) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		user := credential.User

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		removed, removeErr := r.app.Services.Ticket.Remove(ctx, id, user.UserID)
		if removeErr != nil {
			switch {
			case errors.Is(removeErr, domain.ErrNotFound):
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, id)
			case errors.Is(removeErr, domain.ErrForbidden):
				rc = gqlapi.ReturnCodeForbidden
				err = fmt.Errorf("错误代码: %s, 错误信息: 只有工单(uuid=%s)的发起人可以删除工单。", rc, id)
			case errors.Is(removeErr, domain.ErrInvalidInput):
				err = fmt.Errorf("错误代码: %s, 错误信息: 关闭、已执行成功和执行失败的工单不可删除。", rc)
			default:
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, removeErr.Error())
			}
			break
		}

		cluster := cache.ClustersMap.Any(func(elem *models.Cluster) bool {
			return elem.ClusterID == removed.ClusterID
		})
		event.Fire(event.EventTicketRemoved, &event.TicketRemoved{
			User:    *user,
			Ticket:  *bridge.TicketFromDomain(removed),
			Cluster: *cluster,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// PatchTicketStatus 修改工单状态
// VLD_FAILURE -> CLOSED
// VLD_WARNING -> CLOSED
// WAITING_FOR_MRV -> CLOSED
// +
// VLD_WARNING -> MRV_FAILURE
// WAITING_FOR_MRV -> MRV_FAILURE
// VLD_WARNING -> LGTM
// WAITING_FOR_MRV -> LGTM 同时需要指定执行时间，客户端根据要求执行ExecuteTicket或者ScheduleTicket方法
func (r *mutationRootResolver) PatchTicketStatus(ctx context.Context, input models.PatchTicketStatusInput) (ok bool, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		user := credential.User

		statusValue, valid := gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnum(input.Status)]
		if !valid {
			err = fmt.Errorf("错误代码: %s, 错误信息: 无效的工单状态。", rc)
			break
		}

		t, patchErr := r.app.Services.Ticket.PatchStatus(ctx, input.TicketUUID, domain.TicketStatus(statusValue))
		if patchErr != nil {
			var te *service.TransitionError
			switch {
			case errors.Is(patchErr, domain.ErrNotFound):
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			case errors.As(patchErr, &te):
				switch te.Kind {
				case service.TransitionTerminal:
					err = fmt.Errorf("错误代码: %s, 错误信息: 关闭、已执行成功和执行失败的工单不可编辑。", rc)
				case service.TransitionClose:
					err = fmt.Errorf("错误代码: %s, 错误信息: 当前工单(uuid=%s)状态(status=%s)不允许执行关闭操作。", rc, input.TicketUUID, ticketStatusName(te.Current))
				case service.TransitionManualReview:
					err = fmt.Errorf("错误代码: %s, 错误信息: 当前工单(uuid=%s)状态(status=%s)不允许执行人工审核操作。", rc, input.TicketUUID, ticketStatusName(te.Current))
				case service.TransitionInvalidTarget:
					err = fmt.Errorf("错误代码: %s, 错误信息: 无效的工单状态修改请求。", rc)
				}
			default:
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, patchErr.Error())
			}
			break
		}
		ticket := bridge.TicketFromDomain(t)

		cluster := cache.ClustersMap.Any(func(elem *models.Cluster) bool {
			return elem.ClusterID == ticket.ClusterID
		})
		event.Fire(event.EventTicketStatusPatched, &event.TicketStatusPatched{
			User:    *user,
			Ticket:  *ticket,
			Cluster: *cluster,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// ticketStatusName maps a domain status back to its GraphQL enum name, used to
// render the same messages as the legacy endpoint.
func ticketStatusName(s domain.TicketStatus) gqlapi.TicketStatusEnum {
	for k, v := range gqlapi.TicketStatusEnumMap {
		if v == uint8(s) {
			return k
		}
	}
	return ""
}

// ExecuteTicket 执行一个工单，公用ScheduleTicket，所有错误由ScheduleTicket处理
func (r *mutationRootResolver) ExecuteTicket(ctx context.Context, id string) (ok bool, err error) {
	input := models.ScheduleTicketInput{
		TicketUUID: id,
		Schedule:   time.Now().UTC().Format("2006-01-02 15:04:05"),
	}

	if _, err = r.ScheduleTicket(ctx, input); err == nil {
		ok = true
	}

	return
}

// ScheduleTicket 预约执行一个工单
func (r *mutationRootResolver) ScheduleTicket(ctx context.Context, input models.ScheduleTicketInput) (cron *models.Cron, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		user := credential.User

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		t, getErr := r.app.Services.Ticket.GetByUUID(ctx, input.TicketUUID)
		if getErr != nil {
			if errors.Is(getErr, domain.ErrNotFound) {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			} else {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
			}
			break
		}
		ticket := bridge.TicketFromDomain(t)

		// 判断工单的状态，只有系统审核成功的状态才允许执行
		if ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumLgtm] {
			// TODO: 处理rc，处理错误信息
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)当前不可执行。", rc, ticket.UUID)
			break
		}

		// 不判断用户状态，因为这个用户是从Context中获取的，前面的代码已经做了判断
		if ticket.ReviewerID != user.UserID {
			// TODO: 处理rc，处理错误信息
			err = fmt.Errorf("错误代码: %s, 错误信息: xxxx。", rc)
			break
		}

		// 工单当前状态是LGTM才允许执行
		if ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumLgtm] {
			err = fmt.Errorf("错误代码: %s, 错误信息: xxxx。", rc)
			break
		}

		// 是否需要判断群集的状态
		cluster := cache.ClustersMap.Any(func(elem *models.Cluster) bool {
			return elem.ClusterID == ticket.ClusterID
		})
		if cluster == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)目标群集不存在。", rc, ticket.UUID)
			break
		}

		if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
			rc = gqlapi.ReturnCodeClusterNotAvailable
			err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, cluster.UUID)
			break
		}

		// 启动一个计划
		s := worker.NewScheduler()
		local, _ := time.LoadLocation("Local")
		when, _ := time.ParseInLocation("2006-01-02 15:04:05", input.Schedule, local)
		cronUUID, _ := s.RunAt(when, ticket.Subject, "execute", "-T", ticket.UUID)

		domainCron, cronErr := r.app.Services.Cron.GetByUUID(ctx, cronUUID)
		if cronErr != nil {
			cron = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, cronErr.Error())
			break
		}
		// 定时任务归属到预约人与目标群集，审计时能回答「谁预约的、跑在哪」。
		domainCron.UserID = user.UserID
		domainCron.ClusterID = ticket.ClusterID
		if updErr := r.app.Services.Cron.Update(ctx, domainCron, user.UserID); updErr != nil {
			cron = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, updErr.Error())
			break
		}
		cron = bridge.CronFromDomain(domainCron)

		if updErr := r.app.Services.Ticket.UpdateCron(ctx, input.TicketUUID, cron.CronID); updErr != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, updErr.Error())
			break
		}
		ticket.CronID = sql.NullInt64{
			Int64: int64(cron.CronID),
			Valid: true,
		}

		event.Fire(event.EventTicketScheduled, &event.TicketScheduled{
			User:    *user,
			Ticket:  *ticket,
			Cluster: *cluster,
			Cron:    *cron,
		})

		break
	}

	return
}

// Ticket 查看一个工单
func (r *queryRootResolver) Ticket(ctx context.Context, id string) (ticket *models.Ticket, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		t, getErr := r.app.Services.Ticket.GetByUUID(ctx, id)
		if getErr != nil {
			if errors.Is(getErr, domain.ErrNotFound) {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, id)
			} else {
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
			}
			break
		}
		ticket = bridge.TicketFromDomain(t)

		break
	}

	if err != nil {
		ticket = nil
	}

	return
}

// Tickets 分页查看全部工单
func (r *queryRootResolver) Tickets(ctx context.Context, after *string, before *string, first *int, last *int, status *uint8, keyword *string) (*gqlapi.TicketConnection, error) {
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
	edges := []*gqlapi.TicketEdge{}

	// 列表不允许查询Content内容
	filter := service.TicketFilter{CursorID: int64(from), Limit: *first}
	if status != nil {
		filter.Status = *status
	}
	if keyword != nil {
		filter.Keyword = *keyword
	}
	domainTickets, err := r.app.Services.Ticket.ListFiltered(ctx, filter)
	if err != nil {
		return nil, err
	}
	tickets := make([]*models.Ticket, 0, len(domainTickets))
	for _, t := range domainTickets {
		m := bridge.TicketFromDomain(t)
		m.Content = ""
		tickets = append(tickets, m)
	}

	for _, ticket := range tickets {
		edges = append(edges, &gqlapi.TicketEdge{
			Node:   ticket,
			Cursor: EncodeCursor(fmt.Sprintf("%d", ticket.TicketID)),
		})
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return nil, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", tickets[0].TicketID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", tickets[len(tickets)-1].TicketID))
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.TicketConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// TicketSearch 工单搜索，TODO: 下一个版本考虑实现
func (r *queryRootResolver) TicketSearch(ctx context.Context, search string, after *string, before *string, first *int, last *int) (*gqlapi.TicketConnection, error) {
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

	panic("not implemented")
}

type ticketResolver struct{ *Resolver }

// Cluster 工单的目标群集
func (r *ticketResolver) Cluster(ctx context.Context, obj *models.Ticket) (cluster *models.Cluster, err error) {
	rc := gqlapi.ReturnCodeOK
	cluster = cache.ClustersMap.Any(func(elem *models.Cluster) bool {
		return elem.ClusterID == obj.ClusterID
	})
	if cluster == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)依赖的群集不存在。", rc, obj.UUID)
	}
	return
}

// User 工单的发起人
func (r *ticketResolver) User(ctx context.Context, obj *models.Ticket) (user *models.User, err error) {
	rc := gqlapi.ReturnCodeOK
	user = cache.UsersMap.Any(func(elem *models.User) bool {
		return elem.UserID == obj.UserID
	})

	if user == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)的发起人不存在。", rc, obj.UUID)
	}
	return
}

// Reviewer 工单的审核人
func (r *ticketResolver) Reviewer(ctx context.Context, obj *models.Ticket) (user *models.User, err error) {
	rc := gqlapi.ReturnCodeOK
	user = cache.UsersMap.Any(func(elem *models.User) bool {
		return elem.UserID == obj.ReviewerID
	})

	if user == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)的审核人不存在。", rc, obj.UUID)
	}
	return
}

// Cron 执行预约信息
func (r *ticketResolver) Cron(ctx context.Context, obj *models.Ticket) (cron *models.Cron, err error) {
	rc := gqlapi.ReturnCodeOK
	if !obj.CronID.Valid {
		return
	}
	c, getErr := r.app.Services.Cron.Get(ctx, uint(obj.CronID.Int64))
	if getErr != nil {
		cron = nil
		rc = gqlapi.ReturnCodeUnknowError
		err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
		return
	}
	cron = bridge.CronFromDomain(c)

	return
}

// Statements 工单的分解语句，TODO: 分页未完成
func (r *ticketResolver) Statements(ctx context.Context, obj *models.Ticket, after *string, before *string, first *int, last *int) (*gqlapi.StatementConnection, error) {
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

	domainStmts, err := r.app.Services.Statement.ListByTicket(ctx, uint(obj.TicketID))
	if err != nil {
		return nil, err
	}
	stmts := make([]*models.Statement, 0, len(domainStmts))
	for _, s := range domainStmts {
		stmts = append(stmts, bridge.StatementFromDomain(s))
	}
	edges := []*gqlapi.StatementEdge{}
	for _, stmt := range stmts {
		edges = append(edges, &gqlapi.StatementEdge{
			Node:   stmt,
			Cursor: EncodeCursor(fmt.Sprintf("%d", stmt.Sequence)),
		})
	}
	// Statements 是非空字段，空列表要返回空连接而不是 nil
	if len(edges) == 0 {
		return &gqlapi.StatementConnection{
			PageInfo:   &gqlapi.PageInfo{},
			Edges:      []*gqlapi.StatementEdge{},
			TotalCount: 0,
		}, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", stmts[0].Sequence))
	endCursor := EncodeCursor(fmt.Sprintf("%d", stmts[len(stmts)-1].Sequence))
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.StatementConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// Comments 工单的审核意见
func (r *ticketResolver) Comments(ctx context.Context, obj *models.Ticket, after *string, before *string, first *int, last *int) (*gqlapi.CommentConnection, error) {
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

	domainComments, err := r.app.Services.Comment.ListByTicket(ctx, uint(obj.TicketID))
	if err != nil {
		return nil, err
	}
	comments := make([]*models.Comment, 0, len(domainComments))
	for _, c := range domainComments {
		comments = append(comments, bridge.CommentFromDomain(c))
	}
	edges := []*gqlapi.CommentEdge{}
	for _, comment := range comments {
		edges = append(edges, &gqlapi.CommentEdge{
			Node:   comment,
			Cursor: EncodeCursor(fmt.Sprintf("%d", comment.CommentID)),
		})
	}
	// 空列表返回空数组，避免调用方拿到 null
	if len(edges) == 0 {
		return &gqlapi.CommentConnection{
			PageInfo: &gqlapi.PageInfo{},
			Edges:    []*gqlapi.CommentEdge{},
		}, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", comments[0].CommentID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", comments[len(comments)-1].CommentID))
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.CommentConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// 校验工单详情并处理请求结果
func (r *mutationRootResolver) validation(stmts []*models.Statement, cluster *models.Cluster, ticket *models.Ticket) {
	validate.Run(stmts, cluster, ticket)

	for _, s := range stmts {
		s.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForMrv]
		clauses := s.Violations.Clauses()
		if len(clauses) == 0 {
			continue
		}
		// 如果有问题，至少先是警告
		s.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldWarning]

		for _, c := range clauses {
			// 如果存在严重的问题，则标记失败
			if c.Level == 1 {
				s.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldFailure]
				break
			}
		}
		s.Report = s.Violations.Marshal()
	}

	ticket.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForMrv]
	// 默认验证通过，遍历语句，如果存在有警告的语句，则标记为警告，继续。如果存在失败的语句，标记失败并退出
	for _, s := range stmts {
		switch s.Status {
		case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldWarning]:
			ticket.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldWarning]
		case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldFailure]:
			ticket.Status = gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldFailure]
		}
		if ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldFailure] {
			break
		}
	}

	// 命中明细落库：文本报告无法按级别或规则聚合，仪表盘需要明细。
	hits := collectRuleHits(ticket, stmts)

	// 更新到数据库
	err := r.app.Tx.WithTx(context.Background(), func(ctx context.Context, repos service.Repos) error {
		for _, stmt := range stmts {
			if err := repos.Statement.Update(ctx, bridge.StatementToDomain(stmt)); err != nil {
				return err
			}
		}
		// 重新校验时先清掉该工单的旧明细，避免留下两份命中。
		if _, err := repos.RuleHit.DeleteByTicket(ctx, ticket.TicketID); err != nil {
			return err
		}
		if err := repos.RuleHit.CreateBatch(ctx, hits); err != nil {
			return err
		}
		return repos.Ticket.Update(ctx, bridge.TicketToDomain(ticket))
	})
	if err != nil {
		log.Errorf("[E] An unexpected error occured during data updating, err: %s", err.Error())
	}
}

// collectRuleHits turns the clause hits of every statement into persistable
// rows. The rule name is filled in by the rule dispatcher.
func collectRuleHits(ticket *models.Ticket, stmts []*models.Statement) []*domain.RuleHit {
	if ticket == nil {
		return nil
	}
	hits := make([]*domain.RuleHit, 0)
	for _, stmt := range stmts {
		if stmt == nil || stmt.Violations == nil {
			continue
		}
		for _, clause := range stmt.Violations.Clauses() {
			hits = append(hits, &domain.RuleHit{
				UUID:     uuid.NewString(),
				TicketID: ticket.TicketID,
				Sequence: stmt.Sequence,
				RuleName: clause.Rule,
				Level:    clause.Level,
				Message:  clause.Description,
			})
		}
	}
	return hits
}

// StatementType2Uint8 generates a label for a statement.
func StatementType2Uint8(node ast.StmtNode) uint8 {
	switch node.(type) {
	case *ast.AlterTableStmt:
		return 1
	case *ast.AnalyzeTableStmt:
		return 2
	case *ast.BeginStmt:
		return 3
	case *ast.CommitStmt:
		return 4
	case *ast.CreateDatabaseStmt:
		return 5
	case *ast.CreateIndexStmt:
		return 6
	case *ast.CreateTableStmt:
		return 7
	case *ast.CreateViewStmt:
		return 8
	case *ast.CreateUserStmt:
		return 9
	case *ast.DeleteStmt:
		return 10
	case *ast.DropDatabaseStmt:
		return 11
	case *ast.DropIndexStmt:
		return 12
	case *ast.DropTableStmt:
		return 13
	case *ast.ExplainStmt:
		return 14
	case *ast.InsertStmt:
		return 15
	case *ast.LoadDataStmt:
		return 16
	case *ast.RollbackStmt:
		return 17
	case *ast.SelectStmt:
		return 18
	case *ast.SetStmt, *ast.SetPwdStmt:
		return 19
	case *ast.ShowStmt:
		return 20
	case *ast.TruncateTableStmt:
		return 21
	case *ast.UpdateStmt:
		return 22
	case *ast.GrantStmt:
		return 23
	case *ast.RevokeStmt:
		return 24
	case *ast.DeallocateStmt:
		return 25
	case *ast.ExecuteStmt:
		return 26
	case *ast.PrepareStmt:
		return 27
	case *ast.UseStmt:
		return 28
	}
	return 0
}
