package resolvers

import (
	"context"

	"github.com/mia0x75/halo/internal/platform/config"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Messages 当前用户的站内消息
func (r *queryRootResolver) Messages(ctx context.Context, unreadOnly *bool, first *int) (L []*gqlapi.Message, err error) {
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	limit := 20
	if first != nil && *first > 0 {
		limit = *first
	}

	messages, err := r.app.Services.Message.ListByUser(ctx, credential.User.UserID, unreadOnly != nil && *unreadOnly, limit)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		L = append(L, &gqlapi.Message{
			UUID:     message.UUID,
			Category: message.Category,
			Title:    message.Title,
			Content:  message.Content,
			Link:     message.Link,
			ReadAt:   message.ReadAt,
			CreateAt: message.CreateAt,
		})
	}
	return L, nil
}

// UnreadMessageCount 当前用户的未读消息数
func (r *queryRootResolver) UnreadMessageCount(ctx context.Context) (uint, error) {
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	return r.app.Services.Message.CountUnread(ctx, credential.User.UserID)
}

// MarkMessageRead 把一条消息标记为已读
func (r *mutationRootResolver) MarkMessageRead(ctx context.Context, uuid string) (*bool, error) {
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	if err := r.app.Services.Message.MarkRead(ctx, uuid, credential.User.UserID); err != nil {
		return nil, err
	}
	ok := true
	return &ok, nil
}

// MarkAllMessagesRead 把当前用户的全部未读消息标记为已读，返回更新条数
func (r *mutationRootResolver) MarkAllMessagesRead(ctx context.Context) (int, error) {
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	rows, err := r.app.Services.Message.MarkAllRead(ctx, credential.User.UserID)
	if err != nil {
		return 0, err
	}
	return int(rows), nil
}
