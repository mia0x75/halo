package mysqlrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestStatisticRepo(t *testing.T) {
	repo := NewStatisticRepo(setupTestDB(t))
	ctx := context.Background()

	s := &domain.Statistic{
		Group:    "group-" + newSuffix(),
		Key:      "key",
		UUID:     newUUID(),
		Value:    12.3456,
		Version:  1,
		CreateAt: nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, s))

	got, err := repo.GetByUUID(ctx, s.UUID)
	require.NoError(t, err)
	require.Equal(t, 12.3456, got.Value)

	byKey, err := repo.GetByKey(ctx, s.Group, s.Key)
	require.NoError(t, err)
	require.Equal(t, s.UUID, byKey.UUID)

	updated := nowUnix()
	s.Value = 99.5
	s.Version = 2
	s.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, s))

	got, err = repo.GetByUUID(ctx, s.UUID)
	require.NoError(t, err)
	require.Equal(t, 99.5, got.Value)

	list, err := repo.List(ctx, s.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasStatistic(list, s.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, s.UUID))
	_, err = repo.GetByUUID(ctx, s.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasStatistic(list []*domain.Statistic, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
