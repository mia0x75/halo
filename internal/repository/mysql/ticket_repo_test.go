package mysqlrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestTicketRepo(t *testing.T) {
	repo := NewTicketRepo(setupTestDB(t))
	ctx := context.Background()

	tk := &domain.Ticket{
		UUID:       newUUID(),
		ClusterID:  1,
		Database:   "halo",
		Subject:    "add index",
		Content:    "alter table ...",
		Status:     domain.TicketStatusWaitingForVld,
		UserID:     1,
		ReviewerID: 2,
		Version:    1,
		CreateAt:   nowUnix(),
	}
	id, err := repo.Create(ctx, tk)
	require.NoError(t, err)
	require.NotZero(t, id)
	tk.TicketID = uint(id)

	got, err := repo.GetByUUID(ctx, tk.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusWaitingForVld, got.Status)
	require.Nil(t, got.CronID)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, tk.UUID, byID.UUID)

	updated := nowUnix()
	tk.Status = domain.TicketStatusLgtm
	tk.Version = 2
	tk.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, tk))

	got, err = repo.GetByUUID(ctx, tk.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusLgtm, got.Status)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasTicket(list, tk.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, tk.UUID))
	_, err = repo.GetByUUID(ctx, tk.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasTicket(list []*domain.Ticket, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
