package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestOptionRepo(t *testing.T) {
	repo := NewOptionRepo(setupTestDB(t))
	ctx := context.Background()

	o := &domain.Option{
		Name:        "option-" + newSuffix(),
		UUID:        newUUID(),
		Value:       "true",
		Description: "an option",
		Element:     "switch",
		Writable:    1,
		Version:     1,
		CreateAt:    nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, o))

	got, err := repo.GetByUUID(ctx, o.UUID)
	require.NoError(t, err)
	require.Equal(t, o.Value, got.Value)

	byName, err := repo.GetByName(ctx, o.Name)
	require.NoError(t, err)
	require.Equal(t, o.UUID, byName.UUID)

	updated := nowUnix()
	o.Value = "false"
	o.Version = 2
	o.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, o))

	got, err = repo.GetByUUID(ctx, o.UUID)
	require.NoError(t, err)
	require.Equal(t, "false", got.Value)

	list, err := repo.List(ctx, o.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasOption(list, o.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, o.UUID))
	_, err = repo.GetByUUID(ctx, o.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasOption(list []*domain.Option, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
