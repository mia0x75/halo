package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestMessageRepo(t *testing.T) {
	repo := NewMessageRepo(setupTestDB(t))
	ctx := context.Background()

	userID := uint(newID())
	msg := &domain.Message{
		UUID:     newUUID(),
		UserID:   userID,
		Category: "ticket",
		Title:    "工单状态更新",
		Content:  "工单「订单归档字段」的状态已更新。",
		Link:     "/tickets/x",
	}
	require.NoError(t, repo.Create(ctx, msg))

	got, err := repo.GetByUUID(ctx, msg.UUID)
	require.NoError(t, err)
	require.Equal(t, "工单状态更新", got.Title)
	require.Nil(t, got.ReadAt)

	count, err := repo.CountUnread(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, uint(1), count)

	unread, err := repo.ListByUser(ctx, userID, true, 10)
	require.NoError(t, err)
	require.Len(t, unread, 1)

	require.NoError(t, repo.MarkRead(ctx, msg.UUID, userID))
	count, err = repo.CountUnread(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, uint(0), count)

	all, err := repo.ListByUser(ctx, userID, false, 10)
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.NotNil(t, all[0].ReadAt)

	second := &domain.Message{
		UUID:     newUUID(),
		UserID:   userID,
		Category: "ticket",
		Title:    "第二条消息",
		Content:  "内容",
	}
	require.NoError(t, repo.Create(ctx, second))
	rows, err := repo.MarkAllRead(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int64(1), rows)

	count, err = repo.CountUnread(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, uint(0), count)

	again, err := repo.MarkAllRead(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, int64(0), again)
}
