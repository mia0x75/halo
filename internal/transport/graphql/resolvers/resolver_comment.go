package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// CreateComment 添加审核意见
func (r *mutationRootResolver) CreateComment(ctx context.Context, input models.CreateCommentInput) (comment *models.Comment, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
		user := credential.User

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

		// TODO: 增加ticket.Staus的状态判断，只有在某几个状态才允许添加审核意见
		if ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumVldWarning] &&
			ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumWaitingForMrv] &&
			ticket.Status != gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumMrvFailure] {
			err = fmt.Errorf("错误代码: %s, 错误信息: 工单的当前状态不允许添加审核意见。", rc)
			break
		}

		comment = &models.Comment{
			Content:  input.Content,
			UserID:   user.UserID,
			TicketID: ticket.TicketID,
		}
		comment.BeforeInsert()

		newID, createErr := r.app.Services.Comment.Create(ctx, bridge.CommentToDomain(comment), user.UserID)
		if createErr != nil {
			comment = nil
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, createErr.Error())
			break
		}
		comment.CommentID = uint(newID)

		event.Fire(event.EventCommentCreated, &event.CommentCreated{
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
	user = cache.UsersMap.Any(func(elem *models.User) bool {
		return elem.UserID == obj.UserID
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
	t, getErr := r.app.Services.Ticket.Get(ctx, obj.TicketID)
	if getErr != nil {
		ticket = nil
		rc = gqlapi.ReturnCodeUnknowError
		err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, getErr.Error())
		return
	}
	ticket = bridge.TicketFromDomain(t)
	return
}
