package postgresrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestCronRepo(t *testing.T) {
	repo := NewCronRepo(setupTestDB(t))
	ctx := context.Background()

	recurrent := uint8(1)
	c := &domain.Cron{
		UUID:      newUUID(),
		Status:    domain.CronStatusPending,
		Name:      "cron-" + newSuffix(),
		Cmd:       "cleanup",
		Params:    "{}",
		Recurrent: &recurrent,
		Version:   1,
		CreateAt:  nowUnix(),
	}
	id, err := repo.Create(ctx, c)
	require.NoError(t, err)
	require.NotZero(t, id)
	c.CronID = uint(id)

	got, err := repo.GetByUUID(ctx, c.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.CronStatusPending, got.Status)
	require.NotNil(t, got.Recurrent)
	require.Equal(t, uint8(1), *got.Recurrent)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, c.UUID, byID.UUID)

	lastRun := "2026-01-01 00:00:00"
	updated := nowUnix()
	c.Status = domain.CronStatusSuccess
	c.LastRun = &lastRun
	c.Version = 2
	c.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, c))

	got, err = repo.GetByUUID(ctx, c.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.CronStatusSuccess, got.Status)
	require.NotNil(t, got.LastRun)
	require.Equal(t, lastRun, *got.LastRun)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasCron(list, c.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, c.UUID))
	_, err = repo.GetByUUID(ctx, c.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

// Ticket bookings must stay out of the system job list and its count.
func TestCronRepoSkipsTicketBookings(t *testing.T) {
	repo := NewCronRepo(setupTestDB(t))
	ctx := context.Background()

	base, err := repo.Count(ctx)
	require.NoError(t, err)

	oneShot := uint8(0)
	booking := &domain.Cron{
		UUID:      newUUID(),
		Status:    domain.CronStatusPending,
		Name:      "booking-" + newSuffix(),
		Cmd:       "halocli",
		Params:    `["execute","-T","some-ticket"]`,
		Recurrent: &oneShot,
		Version:   1,
		CreateAt:  nowUnix(),
	}
	_, err = repo.Create(ctx, booking)
	require.NoError(t, err)

	list, err := repo.List(ctx, int64(math.MaxUint32), 1000)
	require.NoError(t, err)
	require.False(t, hasCron(list, booking.UUID))

	after, err := repo.Count(ctx)
	require.NoError(t, err)
	require.Equal(t, base, after)
}

func hasCron(list []*domain.Cron, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
