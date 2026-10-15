package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestGlossaryRepo(t *testing.T) {
	repo := NewGlossaryRepo(setupTestDB(t))
	ctx := context.Background()

	g := &domain.Glossary{
		Group:       "group-" + newSuffix(),
		Key:         1,
		Value:       "value-" + newSuffix(),
		UUID:        newUUID(),
		Description: "a glossary entry",
		Version:     1,
		CreateAt:    nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, g))

	got, err := repo.GetByUUID(ctx, g.UUID)
	require.NoError(t, err)
	require.Equal(t, g.Value, got.Value)
	require.Equal(t, uint8(1), got.Key)

	byKey, err := repo.GetByKey(ctx, g.Group, g.Key)
	require.NoError(t, err)
	require.Equal(t, g.UUID, byKey.UUID)

	updated := nowUnix()
	g.Value = "value-updated-" + newSuffix()
	g.Version = 2
	g.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, g))

	got, err = repo.GetByUUID(ctx, g.UUID)
	require.NoError(t, err)
	require.Equal(t, g.Value, got.Value)

	list, err := repo.List(ctx, g.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasGlossary(list, g.UUID))

	require.NoError(t, repo.DeleteByKey(ctx, g.Group, g.Key))
	_, err = repo.GetByUUID(ctx, g.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasGlossary(list []*domain.Glossary, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
