package postgresrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestEdgeRepo(t *testing.T) {
	repo := NewEdgeRepo(setupTestDB(t))
	ctx := context.Background()

	e := &domain.Edge{
		UUID:         newUUID(),
		Type:         domain.EdgeTypeUserToRole,
		AncestorID:   newID(),
		DescendantID: newID(),
		Version:      1,
		CreateAt:     nowUnix(),
	}
	id, err := repo.Create(ctx, e)
	require.NoError(t, err)
	require.NotZero(t, id)
	e.EdgeID = uint(id)

	got, err := repo.GetByUUID(ctx, e.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.EdgeTypeUserToRole, got.Type)
	require.Equal(t, e.AncestorID, got.AncestorID)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, e.UUID, byID.UUID)

	updated := nowUnix()
	e.Type = domain.EdgeTypeUserToCluster
	e.Version = 2
	e.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, e))

	got, err = repo.GetByUUID(ctx, e.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.EdgeTypeUserToCluster, got.Type)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasEdge(list, e.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, e.UUID))
	_, err = repo.GetByUUID(ctx, e.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasEdge(list []*domain.Edge, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
