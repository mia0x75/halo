package postgresrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestUserRepo(t *testing.T) {
	repo := NewUserRepo(setupTestDB(t))
	ctx := context.Background()

	u := &domain.User{
		UUID:     newUUID(),
		Email:    newEmail(),
		Password: "hashed",
		Status:   domain.UserStatusNormal,
		Name:     "tester",
		Phone:    13800000000,
		AvatarID: 1,
		Version:  1,
		CreateAt: nowUnix(),
	}
	id, err := repo.Create(ctx, u)
	require.NoError(t, err)
	require.NotZero(t, id)
	u.UserID = uint(id)

	got, err := repo.GetByUUID(ctx, u.UUID)
	require.NoError(t, err)
	require.Equal(t, u.Email, got.Email)
	require.Equal(t, uint64(13800000000), got.Phone)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, u.UUID, byID.UUID)

	updated := nowUnix()
	u.Name = "tester-2"
	u.Status = domain.UserStatusBlocked
	u.Version = 2
	u.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, u))

	got, err = repo.GetByUUID(ctx, u.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.UserStatusBlocked, got.Status)
	require.Equal(t, "tester-2", got.Name)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasUser(list, u.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, u.UUID))
	_, err = repo.GetByUUID(ctx, u.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasUser(list []*domain.User, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
