package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

type fakeCleanupRepo struct {
	domain.MetricRepository
	deletedBucket string
	deletedBefore uint
}

func (f *fakeCleanupRepo) DeleteBucketsBefore(_ context.Context, bucket string) (int64, error) {
	f.deletedBucket = bucket
	return 3, nil
}

type fakeCleanupHitRepo struct {
	domain.RuleHitRepository
	deletedBefore uint
}

func (f *fakeCleanupHitRepo) DeleteBefore(_ context.Context, before uint) (int64, error) {
	f.deletedBefore = before
	return 7, nil
}

func TestMetricCleanup(t *testing.T) {
	ctx := context.Background()
	metrics := &fakeCleanupRepo{}
	hits := &fakeCleanupHitRepo{}
	svc := NewMetricService(&fakeTx{}, metrics, &fakeTicketCountRepo{}, hits)

	result, err := svc.Cleanup(ctx, 30)
	require.NoError(t, err)
	require.Equal(t, int64(3), result.Metrics)
	require.Equal(t, int64(7), result.Violations)
	require.NotEmpty(t, metrics.deletedBucket)
	require.NotZero(t, hits.deletedBefore)
}

func TestMetricCleanupRejectsNonPositiveWindow(t *testing.T) {
	ctx := context.Background()
	svc := NewMetricService(&fakeTx{}, &fakeCleanupRepo{}, &fakeTicketCountRepo{}, &fakeCleanupHitRepo{})

	_, err := svc.Cleanup(ctx, 0)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}
