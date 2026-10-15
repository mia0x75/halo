package postgresrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestQueryRepo(t *testing.T) {
	repo := NewQueryRepo(setupTestDB(t))
	ctx := context.Background()

	q := &domain.Query{
		UUID:      newUUID(),
		Type:      domain.QueryTypeQuery,
		ClusterID: 1,
		Database:  "halo",
		Content:   "select 1",
		Plan:      "plan",
		UserID:    1,
		Version:   1,
		CreateAt:  nowUnix(),
	}
	id, err := repo.Create(ctx, q)
	require.NoError(t, err)
	require.NotZero(t, id)
	q.QueryID = uint(id)

	got, err := repo.GetByUUID(ctx, q.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.QueryTypeQuery, got.Type)
	require.Equal(t, q.Content, got.Content)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, q.UUID, byID.UUID)

	updated := nowUnix()
	q.Type = domain.QueryTypeAnalyze
	q.Version = 2
	q.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, q))

	got, err = repo.GetByUUID(ctx, q.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.QueryTypeAnalyze, got.Type)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasQuery(list, q.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, q.UUID))
	_, err = repo.GetByUUID(ctx, q.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasQuery(list []*domain.Query, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
