package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

// TestThinServicePattern exercises one thin service end to end to lock in the
// shared write pattern: repo write plus audit row, in the same transaction.
func TestThinServicePattern(t *testing.T) {
	ctx := context.Background()
	users := newFakeUserRepo()
	logs := &fakeLogRepo{}
	svc := NewUserService(&fakeTx{repos: Repos{User: users, Log: logs}}, users)

	u := &domain.User{UUID: "u1", Email: "a@b.test"}
	id, err := svc.Create(ctx, u, 7)
	require.NoError(t, err)
	require.NotZero(t, id)

	got, err := svc.GetByUUID(ctx, "u1")
	require.NoError(t, err)
	require.Equal(t, "a@b.test", got.Email)

	require.Len(t, logs.entries, 1)
	require.Equal(t, "user.create|user:u1|", logs.entries[0].Operation)

	require.NoError(t, svc.Delete(ctx, id, 7))
	require.Len(t, logs.entries, 2)
}
