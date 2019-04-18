package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// CreateComment 添加审核意见
func (r *mutationRootResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (comment *models.Comment, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		found := false
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User
		ticket := &models.Ticket{
			UUID: input.TicketUUID,
		}
		found, err = g.Engine.Get(ticket)
		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}
		// TODO: 增加ticket.Staus的状态判断，只有在某几个状态才允许添加审核意见
		if ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldWarning] &&
			ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForMrv] &&
			ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumMrvFailure] {
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单的当前状态不允许添加审核意见。", rc)
			break
		}
		if !found {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单(uuid=%s)不存在。", rc, input.TicketUUID)
			break
		}

		comment = &models.Comment{
			Content:  input.Content,
			UserID:   user.UserID,
			TicketID: ticket.TicketID,
		}

		if _, err = g.Engine.Insert(comment); err != nil {
			comment = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break
		}

		events.Fire(events.EventCommentCreated, &events.CommentCreatedArgs{
			User:    *credential.User,
			Ticket:  *ticket,
			Comment: *comment,
		})

		// 退出for循环
		break
	}
	return
}

type commentResolver struct{ *Resolver }

// User 审核意见的发起人
func (r *commentResolver) User(ctx context.Context, obj *models.Comment) (user *models.User, err error) {
	rc := gqlapi.ReturnCodeOK
	user = caches.UsersMap.Any(func(elem *models.User) bool {
		if elem.UserID == obj.UserID {
			return true
		}
		return false
	})

	if user == nil {
		user = nil
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)不存在。", rc, obj.UUID)
	}
	return
}

// Ticket 审核意见关联的工单
func (r *commentResolver) Ticket(ctx context.Context, obj *models.Comment) (ticket *models.Ticket, err error) {
	rc := gqlapi.ReturnCodeOK
	ticket = &models.Ticket{
		TicketID: obj.TicketID,
	}
	if _, err := g.Engine.Get(ticket); err != nil {
		ticket = nil
		rc = gqlapi.ReturnCodeUnknowError
		err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
	}
	return
}
