package mysqlrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestLogRepo(t *testing.T) {
	repo := NewLogRepo(setupTestDB(t))
	ctx := context.Background()

	l := &domain.Log{
		UUID:      newUUID(),
		UserID:    1,
		Operation: "ticket.create",
		Version:   1,
		CreateAt:  nowUnix(),
	}
	id, err := repo.Create(ctx, l)
	require.NoError(t, err)
	require.NotZero(t, id)
	l.LogID = uint(id)

	got, err := repo.GetByUUID(ctx, l.UUID)
	require.NoError(t, err)
	require.Equal(t, l.Operation, got.Operation)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, l.UUID, byID.UUID)

	l.Operation = "ticket.delete"
	l.Version = 2
	require.NoError(t, repo.Update(ctx, l))

	got, err = repo.GetByUUID(ctx, l.UUID)
	require.NoError(t, err)
	require.Equal(t, "ticket.delete", got.Operation)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasLog(list, l.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, l.UUID))
	_, err = repo.GetByUUID(ctx, l.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasLog(list []*domain.Log, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
