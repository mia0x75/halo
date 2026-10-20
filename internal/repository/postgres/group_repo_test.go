package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestGroupRepo(t *testing.T) {
	repo := NewGroupRepo(setupTestDB(t))
	ctx := context.Background()

	group := &domain.Group{
		UUID:        newUUID(),
		Name:        "group-" + newSuffix(),
		Description: "机房 A",
		Status:      domain.GroupStatusNormal,
		Version:     1,
	}
	require.NoError(t, repo.Create(ctx, group))

	got, err := repo.GetByUUID(ctx, group.UUID)
	require.NoError(t, err)
	require.Equal(t, "机房 A", got.Description)

	byName, err := repo.GetByName(ctx, group.Name)
	require.NoError(t, err)
	require.Equal(t, group.UUID, byName.UUID)

	list, err := repo.List(ctx, group.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasGroup(list, group.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, group.UUID))
	_, err = repo.GetByUUID(ctx, group.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasGroup(list []*domain.Group, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
