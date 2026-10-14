package mysqlrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestRoleRepo(t *testing.T) {
	repo := NewRoleRepo(setupTestDB(t))
	ctx := context.Background()

	r := &domain.Role{
		UUID:        newUUID(),
		Name:        "role-" + newSuffix(),
		Description: "a role",
		Version:     1,
		CreateAt:    nowUnix(),
	}
	id, err := repo.Create(ctx, r)
	require.NoError(t, err)
	require.NotZero(t, id)
	r.RoleID = uint(id)

	got, err := repo.GetByUUID(ctx, r.UUID)
	require.NoError(t, err)
	require.Equal(t, r.Name, got.Name)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, r.UUID, byID.UUID)

	updated := nowUnix()
	r.Description = "updated"
	r.Version = 2
	r.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, r))

	got, err = repo.GetByUUID(ctx, r.UUID)
	require.NoError(t, err)
	require.Equal(t, "updated", got.Description)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasRole(list, r.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, r.UUID))
	_, err = repo.GetByUUID(ctx, r.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasRole(list []*domain.Role, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
