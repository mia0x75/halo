package postgresrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestCommentRepo(t *testing.T) {
	repo := NewCommentRepo(setupTestDB(t))
	ctx := context.Background()

	c := &domain.Comment{
		UUID:     newUUID(),
		Content:  "please review",
		TicketID: 1,
		UserID:   1,
		Version:  1,
		CreateAt: nowUnix(),
	}
	id, err := repo.Create(ctx, c)
	require.NoError(t, err)
	require.NotZero(t, id)
	c.CommentID = uint(id)

	got, err := repo.GetByUUID(ctx, c.UUID)
	require.NoError(t, err)
	require.Equal(t, c.Content, got.Content)
	require.Equal(t, c.TicketID, got.TicketID)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, c.UUID, byID.UUID)

	updated := nowUnix()
	c.Content = "lgtm"
	c.Version = 2
	c.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, c))

	got, err = repo.GetByUUID(ctx, c.UUID)
	require.NoError(t, err)
	require.Equal(t, "lgtm", got.Content)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasComment(list, c.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, c.UUID))
	_, err = repo.GetByUUID(ctx, c.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasComment(list []*domain.Comment, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
