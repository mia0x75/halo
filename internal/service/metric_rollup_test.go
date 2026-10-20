package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

type fakeMetricRepo struct {
	domain.MetricRepository
	values map[string]float64
}

func (f *fakeMetricRepo) SetGauge(_ context.Context, key domain.MetricKey, value float64) error {
	if f.values == nil {
		f.values = map[string]float64{}
	}
	f.values[key.Name+"/"+key.Bucket] = value
	return nil
}

func (f *fakeMetricRepo) GetByKey(_ context.Context, key domain.MetricKey) (*domain.Metric, error) {
	if f.values == nil {
		return nil, domain.ErrNotFound
	}
	value, ok := f.values[key.Name+"/"+key.Bucket]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &domain.Metric{Name: key.Name, Scope: key.Scope, Bucket: key.Bucket, Value: value}, nil
}

type fakeTicketCountRepo struct {
	domain.TicketRepository
	counts []domain.TicketStatusCount
}

func (f *fakeTicketCountRepo) CountByStatus(context.Context, uint, uint) ([]domain.TicketStatusCount, error) {
	return f.counts, nil
}

type fakeHitCountRepo struct {
	domain.RuleHitRepository
	counts []domain.RuleHitCount
}

func (f *fakeHitCountRepo) CountByLevel(context.Context, uint, uint) ([]domain.RuleHitCount, error) {
	return f.counts, nil
}

func TestMetricRollupDay(t *testing.T) {
	ctx := context.Background()
	metrics := &fakeMetricRepo{}
	tickets := &fakeTicketCountRepo{counts: []domain.TicketStatusCount{
		{Status: domain.TicketStatusWaitingForVld, Count: 3},
		{Status: domain.TicketStatusLgtm, Count: 2},
		{Status: domain.TicketStatusDone, Count: 5},
		{Status: domain.TicketStatusVldFailure, Count: 1},
		{Status: domain.TicketStatusExecFailure, Count: 2},
	}}
	hits := &fakeHitCountRepo{counts: []domain.RuleHitCount{
		{Level: uint8(domain.RuleLevelError), Count: 4},
		{Level: uint8(domain.RuleLevelWarning), Count: 6},
	}}

	svc := NewMetricService(&fakeTx{}, metrics, tickets, hits)
	result, err := svc.RollupDay(ctx, "2026-10-19")
	require.NoError(t, err)
	require.Equal(t, "2026-10-19", result.Date)
	require.Equal(t, 8, result.Metrics)

	require.Equal(t, 3.0, metrics.values["tickets-created/2026-10-19"])
	require.Equal(t, 2.0, metrics.values["tickets-approved/2026-10-19"])
	require.Equal(t, 5.0, metrics.values["tickets-executed/2026-10-19"])
	require.Equal(t, 3.0, metrics.values["tickets-failed/2026-10-19"])
	require.Equal(t, 4.0, metrics.values["violations-error/2026-10-19"])
	require.Equal(t, 6.0, metrics.values["violations-warning/2026-10-19"])
	require.Equal(t, 1.0, metrics.values["rollup-done/2026-10-19"])
}

func TestMetricRollupDayIsIdempotent(t *testing.T) {
	ctx := context.Background()
	metrics := &fakeMetricRepo{}
	tickets := &fakeTicketCountRepo{counts: []domain.TicketStatusCount{
		{Status: domain.TicketStatusDone, Count: 2},
	}}
	hits := &fakeHitCountRepo{}

	svc := NewMetricService(&fakeTx{}, metrics, tickets, hits)
	_, err := svc.RollupDay(ctx, "2026-10-19")
	require.NoError(t, err)
	_, err = svc.RollupDay(ctx, "2026-10-19")
	require.NoError(t, err)
	require.Equal(t, 2.0, metrics.values["tickets-executed/2026-10-19"])
}

func TestMetricRollupStatus(t *testing.T) {
	ctx := context.Background()
	metrics := &fakeMetricRepo{}
	svc := NewMetricService(&fakeTx{}, metrics, &fakeTicketCountRepo{}, &fakeHitCountRepo{})

	done, err := svc.RollupStatus(ctx, "2026-10-19")
	require.NoError(t, err)
	require.False(t, done)

	_, err = svc.RollupDay(ctx, "2026-10-19")
	require.NoError(t, err)

	done, err = svc.RollupStatus(ctx, "2026-10-19")
	require.NoError(t, err)
	require.True(t, done)
}

func TestMetricRollupRejectsBadDate(t *testing.T) {
	ctx := context.Background()
	svc := NewMetricService(&fakeTx{}, &fakeMetricRepo{}, &fakeTicketCountRepo{}, &fakeHitCountRepo{})

	_, err := svc.RollupDay(ctx, "2026/10/19")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}
