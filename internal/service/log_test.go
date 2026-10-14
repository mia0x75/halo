package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogServiceFormat(t *testing.T) {
	repo := &fakeLogRepo{}
	svc := NewLogService(repo)

	err := svc.Log(context.Background(), LogEntry{
		UserID:    7,
		Operation: "ticket.create",
		Target:    "ticket:abc",
		Detail:    "why",
	})
	require.NoError(t, err)
	require.Len(t, repo.entries, 1)

	e := repo.entries[0]
	require.Equal(t, uint(7), e.UserID)
	require.Equal(t, "ticket.create|ticket:abc|why", e.Operation)
	require.NotEmpty(t, e.UUID)
	require.Equal(t, 1, e.Version)
	require.NotZero(t, e.CreateAt)
}

func TestLogServicePropagatesError(t *testing.T) {
	svc := NewLogService(&fakeLogRepo{err: context.DeadlineExceeded})
	require.ErrorIs(t, svc.Log(context.Background(), LogEntry{Operation: "x"}), context.DeadlineExceeded)
}
