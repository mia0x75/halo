package postgresrepo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestTaskRepo(t *testing.T) {
	repo := NewTaskRepo(setupTestDB(t))
	ctx := context.Background()

	now := uint(time.Now().Unix())
	task := &domain.Task{
		UUID:        newUUID(),
		Name:        "execute-ticket",
		Type:        "ticket",
		Payload:     `{"ticket":"x"}`,
		Status:      domain.TaskStatusPending,
		Priority:    1,
		MaxAttempt:  3,
		ScheduledAt: now,
		Version:     1,
	}
	require.NoError(t, repo.Create(ctx, task))

	got, err := repo.GetByUUID(ctx, task.UUID)
	require.NoError(t, err)
	require.Equal(t, "execute-ticket", got.Name)
	require.Equal(t, domain.TaskStatusPending, got.Status)

	pending, err := repo.ListPending(ctx, 100)
	require.NoError(t, err)
	require.NotEmpty(t, pending)

	started := uint(time.Now().Unix())
	got.Status = domain.TaskStatusSucceeded
	got.StartedAt = &started
	got.FinishedAt = &started
	got.Version = 2
	require.NoError(t, repo.Update(ctx, got))

	got, err = repo.GetByUUID(ctx, task.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.TaskStatusSucceeded, got.Status)
	require.NotNil(t, got.StartedAt)

	counts, err := repo.CountByStatus(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, counts)

	require.NoError(t, repo.DeleteByUUID(ctx, task.UUID))
	_, err = repo.GetByUUID(ctx, task.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
