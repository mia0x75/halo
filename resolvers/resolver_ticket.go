package resolvers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/core"
	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/driver"
	"github.com/mia0x75/parser/format"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/crons"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
	"github.com/mia0x75/halo/validate"
)

// 保留，不可以删除
var _ = driver.ValueExpr{}

// CreateTicket 创建一个工单
func (r *mutationRootResolver) CreateTicket(ctx context.Context, input models.CreateTicketInput) (ticket *models.Ticket, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		statements := []*models.Statement{}
		user := credential.User

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		// 查询缓存
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == input.ClusterUUID {
				return true
			}
			return false
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

		// 检查群集关联
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == cluster.ClusterID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break
		}

		passwd := func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}

		if _, err = cluster.Stat(input.Database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 查询缓存
		reviewer := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.ReviewerUUID {
				return true
			}
			return false
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
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
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
		stmts := []ast.StmtNode{}
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

		session := g.Engine.NewSession()
		session.Begin()
		defer session.Close()

		if _, err = session.Insert(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}

		i := 1
		for _, node := range stmts {
			var sb strings.Builder
			node.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			sql := sb.String()
			if len(sql) == 0 {
				sql = strings.TrimSpace(node.Text())
			}
			stat := models.Statement{
				Sequence:   uint16(i),
				Content:    sql,
				Type:       StatementType2Uint8(node),
				Status:     gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForVld],
				TicketID:   ticket.TicketID,
				StmtNode:   node,
				Violations: &models.Violations{},
			}
			statements = append(statements, &stat)
			i++
		}
		if _, err = session.Insert(&statements); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}
		if err = session.Commit(); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		notifyTicketStatusChange(user.UUID, ticket.UUID, fmt.Sprintf("新建工单(uuid=%s)成功，系统正在开始启动自动化审核。", ticket.UUID))
		go validation(statements, cluster, ticket)

		events.Fire(events.EventTicketCreated, &events.TicketCreatedArgs{
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
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		statements := []*models.Statement{}
		user := credential.User
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.UUID == input.ClusterUUID {
				return true
			}
			return false
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

		// 检查群集关联
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
				elem.AncestorID == user.UserID &&
				elem.DescendantID == cluster.ClusterID {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)没有关联群集(uuid=%s)。", rc, user.UUID, cluster.UUID)
			break
		}

		passwd := func(c *models.Cluster) []byte {
			bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
			return bs
		}

		if _, err = cluster.Stat(input.Database, passwd); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 查询缓存
		reviewer := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.ReviewerUUID {
				return true
			}
			return false
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
		if !caches.EdgesMap.Include(func(elem *models.Edge) bool {
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
		stmts := []ast.StmtNode{}
		stmts, _, err = p.Parse(input.Content, "", "")
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// 开启事务
		session := g.Engine.NewSession()
		defer session.Close()
		session.Begin()

		//获取数据库原有的ticket信息
		ticket = &models.Ticket{}
		if _, err = session.Where("uuid = ?", input.TicketUUID).Get(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}

		if ticket.TicketID == 0 {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			session.Rollback()
			break
		}

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

		// 删除原来关联语句
		stat := models.Statement{
			TicketID: ticket.TicketID,
		}
		if _, err = session.Where("ticket_id = ?", ticket.TicketID).Delete(&stat); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}
		for _, node := range stmts {
			var sb strings.Builder
			node.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
			sql := sb.String()
			if len(sql) == 0 {
				sql = strings.TrimSpace(node.Text())
			}
			stat := models.Statement{
				Sequence:   uint16(i),
				Type:       StatementType2Uint8(node),
				Content:    sql,
				Status:     gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForVld],
				TicketID:   ticket.TicketID,
				StmtNode:   node,
				Violations: &models.Violations{},
			}
			statements = append(statements, &stat)
			i++
		}
		// 更新工单
		if _, err = session.ID(ticket.TicketID).AllCols().Update(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}

		// 重新添加语句
		if _, err = session.Insert(&statements); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}

		if err = session.Commit(); err != nil {
			session.Rollback()
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		notifyTicketStatusChange(user.UUID, ticket.UUID, fmt.Sprintf("更新工单(uuid=%s)成功，系统正在开始启动自动化审核。", ticket.UUID))
		go validation(statements, cluster, ticket)

		events.Fire(events.EventTicketUpdated, &events.TicketUpdatedArgs{
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
		found := false
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		ticket := &models.Ticket{}
		user := credential.User
		ticket.UUID = id

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		session := g.Engine.NewSession()
		defer session.Close()
		session.Begin()
		//获取数据库原有的ticket信息
		if found, err = session.Where("uuid = ?", ticket.UUID).Get(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		} else if !found {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, id)
			break
		}

		if ticket.UserID != user.UserID {
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 只有工单(uuid=%s)的发起人可以删除工单。", rc, id)
			break
		}
		if ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumClosed] ||
			ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumDone] ||
			ticket.Status == gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumExecFailure] {
			// TODO: 处理rc
			err = fmt.Errorf("错误代码: %s, 错误信息: 关闭、已执行成功和执行失败的工单不可删除。", rc)
			break
		}

		// 删除原来关联语句
		stat := models.Statement{
			TicketID: ticket.TicketID,
		}
		if _, err = session.Where("ticket_id  = ?", ticket.TicketID).Delete(&stat); err != nil {
			session.Rollback()
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		// 删除工单
		if _, err = session.ID(ticket.TicketID).Delete(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			session.Rollback()
			break
		}

		if err = session.Commit(); err != nil {
			session.Rollback()
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.ClusterID == ticket.ClusterID {
				return true
			}
			return false
		})
		events.Fire(events.EventTicketRemoved, &events.TicketRemovedArgs{
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
		found := false
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User
		ticket := &models.Ticket{
			UUID: input.TicketUUID,
		}
		found, err = g.Engine.Where("`uuid` = ?", input.TicketUUID).Get(ticket)
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		if !found {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			break
		}

		var currentStatus gqlapi.TicketStatusEnum
		for k, v := range gqlapi.TicketStatusEnumMap {
			if ticket.Status == v {
				currentStatus = k
				break
			}
		}

		if currentStatus == gqlapi.TicketStatusEnumClosed ||
			currentStatus == gqlapi.TicketStatusEnumDone ||
			currentStatus == gqlapi.TicketStatusEnumExecFailure {
			// TODO: 处理rc
			// rc = g.ReturnCodeTicketClosed
			err = fmt.Errorf("错误代码: %s, 错误信息: 关闭、已执行成功和执行失败的工单不可编辑。", rc)
			break
		}

		if status, ok := gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnum(input.Status)]; !ok {
			err = fmt.Errorf("错误代码: %s, 错误信息: 无效的工单状态。", rc)
			break
		} else {
			switch status {
			case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumClosed]: // 手动关闭工单
				if currentStatus == gqlapi.TicketStatusEnumVldFailure ||
					currentStatus == gqlapi.TicketStatusEnumVldWarning ||
					currentStatus == gqlapi.TicketStatusEnumWaitingForMrv {
					ticket.Status = status
				} else {
					// TODO: 处理rc
					// rc = g.ReturnCodeTicketClosed
					err = fmt.Errorf("错误代码: %s, 错误信息: 当前工单(uuid=%s)状态(status=%s)不允许执行关闭操作。", rc, ticket.UUID, currentStatus)
				}
			case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumLgtm]: // 标记工单人工审核通过
				if currentStatus == gqlapi.TicketStatusEnumVldWarning ||
					currentStatus == gqlapi.TicketStatusEnumWaitingForMrv {
					ticket.Status = status
				} else {
					// TODO: 处理rc
					// rc = g.ReturnCodeTicketClosed
					err = fmt.Errorf("错误代码: %s, 错误信息: 当前工单(uuid=%s)状态(status=%s)不允许执行人工审核操作。", rc, ticket.UUID, currentStatus)
				}
			case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumMrvFailure]: // 标记工单人工审核失败
				if currentStatus == gqlapi.TicketStatusEnumVldWarning ||
					currentStatus == gqlapi.TicketStatusEnumWaitingForMrv {
					ticket.Status = status
				} else {
					// TODO: 处理rc
					// rc = g.ReturnCodeTicketClosed
					err = fmt.Errorf("错误代码: %s, 错误信息: 当前工单(uuid=%s)状态(status=%s)不允许执行人工审核操作。", rc, ticket.UUID, currentStatus)
				}
			default:
				// TODO: 处理rc
				// rc = g.ReturnCodeTicketClosed
				err = fmt.Errorf("错误代码: %s, 错误信息: 无效的工单状态修改请求。", rc)
			}
			if err != nil {
				break
			}
		}
		if _, err = g.Engine.ID(ticket.TicketID).Update(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.ClusterID == ticket.ClusterID {
				return true
			}
			return false
		})
		events.Fire(events.EventTicketStatusPatched, &events.TicketStatusPatchedArgs{
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
		found := false
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User

		if strings.TrimSpace(user.Name) == "" {
			rc = gqlapi.ReturnCodeRegistrationIncomplete
			err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)信息不完整。", rc, user.UUID)
			break
		}

		ticket := &models.Ticket{
			UUID: input.TicketUUID,
		}
		found, err = g.Engine.Get(ticket)
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		if !found {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			break
		}

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
		cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
			if elem.ClusterID == ticket.ClusterID {
				return true
			}
			return false
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
		s := crons.NewScheduler()
		local, _ := time.LoadLocation("Local")
		when, _ := time.ParseInLocation("2006-01-02 15:04:05", input.Schedule, local)
		cronUUID, _ := s.RunAt(when, ticket.Subject, "execute", "-T", ticket.UUID)

		cron = &models.Cron{}
		if _, err = g.Engine.Where("`uuid` = ?", cronUUID).Get(cron); err != nil {
			cron = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		ticket.CronID = sql.NullInt64{
			Int64: int64(cron.CronID),
			Valid: true,
		}
		if _, err = g.Engine.ID(ticket.TicketID).Update(ticket); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		events.Fire(events.EventTicketScheduled, &events.TicketScheduledArgs{
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
		found := false
		ticket = &models.Ticket{
			UUID: id,
		}
		found, err = g.Engine.Get(ticket)
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		if !found {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, id)
			break
		}

		break
	}

	if err != nil {
		ticket = nil
	}

	return
}

// Tickets 分页查看全部工单
func (r *queryRootResolver) Tickets(ctx context.Context, after *string, before *string, first *int, last *int) (*gqlapi.TicketConnection, error) {
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
	edges := []gqlapi.TicketEdge{}
	tickets := []*models.Ticket{}

	var err error

	// 列表不允许查询Content内容
	if err = g.Engine.Omit("content").Desc("ticket_id").Where("ticket_id < ?", from).Limit(*first).Find(&tickets); err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		edges = append(edges, gqlapi.TicketEdge{
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
	pageInfo := gqlapi.PageInfo{
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
	cluster = caches.ClustersMap.Any(func(elem *models.Cluster) bool {
		if elem.ClusterID == obj.ClusterID {
			return true
		}
		return false
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
	user = caches.UsersMap.Any(func(elem *models.User) bool {
		if elem.UserID == obj.UserID {
			return true
		}
		return false
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
	user = caches.UsersMap.Any(func(elem *models.User) bool {
		if elem.UserID == obj.ReviewerID {
			return true
		}
		return false
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
	cron = &models.Cron{}
	if _, err = g.Engine.ID(obj.CronID.Int64).Get(cron); err != nil {
		cron = nil
		rc = gqlapi.ReturnCodeUnknowError
		err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
	}

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

	stmts := []*models.Statement{}
	if err := g.Engine.Where("ticket_id = ?", obj.TicketID).Find(&stmts); err != nil {
		return nil, err
	}
	edges := []gqlapi.StatementEdge{}
	for _, stmt := range stmts {
		edges = append(edges, gqlapi.StatementEdge{
			Node:   stmt,
			Cursor: EncodeCursor(fmt.Sprintf("%d", stmt.Sequence)),
		})
	}
	if len(edges) == 0 {
		return nil, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", stmts[0].Sequence))
	endCursor := EncodeCursor(fmt.Sprintf("%d", stmts[len(stmts)-1].Sequence))
	pageInfo := gqlapi.PageInfo{
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

	comments := []*models.Comment{}
	g.Engine.Where("ticket_id = ?", obj.TicketID).Find(&comments)
	edges := []gqlapi.CommentEdge{}
	for _, comment := range comments {
		edges = append(edges, gqlapi.CommentEdge{
			Node:   comment,
			Cursor: EncodeCursor(fmt.Sprintf("%d", comment.CommentID)),
		})
	}
	if len(edges) == 0 {
		return &gqlapi.CommentConnection{}, nil
	}
	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", comments[0].CommentID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", comments[len(comments)-1].CommentID))
	pageInfo := gqlapi.PageInfo{
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
func validation(stmts []*models.Statement, cluster *models.Cluster, ticket *models.Ticket) {
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

	// 更新到数据库
	session := g.Engine.NewSession()
	defer session.Close()
	session.Begin()
	for _, stmt := range stmts {
		if _, err := session.ID(core.PK{stmt.TicketID, stmt.Sequence}).Update(stmt); err != nil {
			session.Rollback()
			log.Errorf("[E] An unexpected error occured during data updating, err: %s", err.Error())
			return
		}
	}

	// 更新工单
	if _, err := session.ID(ticket.TicketID).Update(ticket); err != nil {
		session.Rollback()
		log.Errorf("[E] An unexpected error occured during data updating, err: %s", err.Error())
		return
	}

	if err := session.Commit(); err != nil {
		log.Errorf("[E] An unexpected error occured during data updating, err: %s", err.Error())
	}
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
