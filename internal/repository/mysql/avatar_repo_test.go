package mysqlrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestAvatarRepo(t *testing.T) {
	repo := NewAvatarRepo(setupTestDB(t))
	ctx := context.Background()

	a := &domain.Avatar{
		UUID:     newUUID(),
		URL:      "https://example.test/" + newSuffix(),
		Version:  1,
		CreateAt: nowUnix(),
	}
	id, err := repo.Create(ctx, a)
	require.NoError(t, err)
	require.NotZero(t, id)
	a.AvatarID = uint(id)

	got, err := repo.GetByUUID(ctx, a.UUID)
	require.NoError(t, err)
	require.Equal(t, a.URL, got.URL)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, a.UUID, byID.UUID)

	updated := nowUnix()
	a.URL = "https://example.test/updated-" + newSuffix()
	a.Version = 2
	a.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, a))

	got, err = repo.GetByUUID(ctx, a.UUID)
	require.NoError(t, err)
	require.Equal(t, a.URL, got.URL)
	require.Equal(t, 2, got.Version)
	require.NotNil(t, got.UpdateAt)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasAvatar(list, a.UUID))

	require.NoError(t, repo.Delete(ctx, id))
	_, err = repo.GetByUUID(ctx, a.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasAvatar(list []*domain.Avatar, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
