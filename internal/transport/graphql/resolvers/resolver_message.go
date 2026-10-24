package resolvers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/service"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Messages 当前用户的站内消息，按游标分页
func (r *queryRootResolver) Messages(
	ctx context.Context,
	unreadOnly *bool,
	category *string,
	keyword *string,
	first *int,
	after *string,
) (*gqlapi.MessageConnection, error) {
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)

	filter := service.MessageFilter{
		UnreadOnly: unreadOnly != nil && *unreadOnly,
		Category:   strings.TrimSpace(valueOrEmpty(category)),
		Keyword:    strings.TrimSpace(valueOrEmpty(keyword)),
		Limit:      20,
	}
	if first != nil && *first > 0 {
		filter.Limit = *first
	}
	beforeID, err := decodeMessageCursor(after)
	if err != nil {
		return nil, err
	}
	filter.BeforeID = beforeID

	// Probe one extra row for hasNextPage. CountByUser ignores Limit, so the
	// probe row never leaks into the totals.
	pageLimit := filter.Limit
	filter.Limit = probeLimit(pageLimit)
	messages, err := r.app.Services.Message.ListByUser(ctx, credential.User.UserID, filter)
	if err != nil {
		return nil, err
	}
	pageMessages, hasNextPage := trimKeysetPage(messages, pageLimit)
	total, err := r.app.Services.Message.CountByUser(ctx, credential.User.UserID, filter)
	if err != nil {
		return nil, err
	}

	edges := make([]*gqlapi.MessageEdge, 0, len(pageMessages))
	for _, message := range pageMessages {
		edges = append(edges, &gqlapi.MessageEdge{
			Cursor: EncodeCursor(strconv.FormatUint(uint64(message.MessageID), 10)),
			Node: &gqlapi.Message{
				UUID:     message.UUID,
				Category: message.Category,
				Title:    message.Title,
				Content:  message.Content,
				Link:     message.Link,
				ReadAt:   message.ReadAt,
				CreateAt: message.CreateAt,
			},
		})
	}

	pageInfo := &gqlapi.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: after != nil && *after != "",
		StartCursor:     "",
		EndCursor:       "",
	}
	if len(edges) > 0 {
		pageInfo.StartCursor = edges[0].Cursor
		pageInfo.EndCursor = edges[len(edges)-1].Cursor
	}

	return &gqlapi.MessageConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: int(total),
	}, nil
}

// decodeMessageCursor 把游标还原成 message_id 上界。首屏没有游标，
// 用哨兵值兜住「小于游标」的查询条件。
func decodeMessageCursor(after *string) (uint, error) {
	if after == nil || *after == "" {
		return math.MaxUint32, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(*after)
	if err != nil {
		return 0, fmt.Errorf("错误代码: %s, 错误信息: 游标无效。", gqlapi.ReturnCodeInvalidParams)
	}
	id, err := strconv.ParseUint(strings.TrimPrefix(string(decoded), "cursor"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("错误代码: %s, 错误信息: 游标无效。", gqlapi.ReturnCodeInvalidParams)
	}
	return uint(id), nil
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

// RemoveMessage 删除一条自己的站内消息
func (r *mutationRootResolver) RemoveMessage(ctx context.Context, uuid string) (bool, error) {
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	if err := r.app.Services.Message.Delete(ctx, uuid, credential.User.UserID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return false, fmt.Errorf("错误代码: %s, 错误信息: 消息(uuid=%s)不存在。", gqlapi.ReturnCodeNotFound, uuid)
		}
		return false, err
	}
	return true, nil
}
