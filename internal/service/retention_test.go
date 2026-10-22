package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

type retentionOptionRepo struct {
	domain.OptionRepository
	values map[string]string
}

func (r retentionOptionRepo) GetByName(_ context.Context, name string) (*domain.Option, error) {
	value, ok := r.values[name]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &domain.Option{Name: name, Parameter: value}, nil
}

type retentionMetricRepo struct {
	domain.MetricRepository
	bucket string
}

func (r *retentionMetricRepo) DeleteBucketsBefore(_ context.Context, bucket string) (int64, error) {
	r.bucket = bucket
	return 1, nil
}

type retentionHitRepo struct {
	domain.RuleHitRepository
	before uint
}

func (r *retentionHitRepo) DeleteBefore(_ context.Context, before uint) (int64, error) {
	r.before = before
	return 2, nil
}

type retentionLogRepo struct {
	domain.LogRepository
	before uint
}

func (r *retentionLogRepo) DeleteBefore(_ context.Context, before uint) (int64, error) {
	r.before = before
	return 3, nil
}

type retentionMessageRepo struct {
	domain.MessageRepository
	before uint
}

func (r *retentionMessageRepo) DeleteBefore(_ context.Context, before uint) (int64, error) {
	r.before = before
	return 4, nil
}

type retentionTaskRepo struct {
	domain.TaskRepository
	before uint
}

func (r *retentionTaskRepo) DeleteBefore(_ context.Context, before uint) (int64, error) {
	r.before = before
	return 5, nil
}

func newRetentionService(options domain.OptionRepository) (*RetentionService, *retentionMetricRepo, *retentionHitRepo, *retentionLogRepo, *retentionMessageRepo, *retentionTaskRepo) {
	metric := &retentionMetricRepo{}
	hit := &retentionHitRepo{}
	log := &retentionLogRepo{}
	message := &retentionMessageRepo{}
	task := &retentionTaskRepo{}
	return NewRetentionService(nil, options, metric, hit, log, message, task), metric, hit, log, message, task
}

func TestRetentionCleanupUsesDefaultsWhenOptionsMissing(t *testing.T) {
	svc, metric, hit, log, message, task := newRetentionService(retentionOptionRepo{})

	result, err := svc.Cleanup(context.Background())
	require.NoError(t, err)
	require.Equal(t, RetentionResult{Metrics: 1, Hits: 2, Logs: 3, Messages: 4, Tasks: 5}, result)

	// 缺省窗口：指标/命中/站内信 90 天，日志 180 天，任务 30 天。
	require.Equal(t, time.Now().AddDate(0, 0, -90).Format("2006-01-02"), metric.bucket)
	requireWithinDays(t, hit.before, 90)
	requireWithinDays(t, log.before, 180)
	requireWithinDays(t, message.before, 90)
	requireWithinDays(t, task.before, 30)
}

func TestRetentionCleanupHonoursConfiguredWindows(t *testing.T) {
	svc, metric, _, log, _, _ := newRetentionService(retentionOptionRepo{values: map[string]string{
		retentionMetricsDays:  "7",
		retentionHitsDays:     "7",
		retentionLogsDays:     "30",
		retentionMessagesDays: "14",
		retentionTasksDays:    "3",
	}})

	_, err := svc.Cleanup(context.Background())
	require.NoError(t, err)

	require.Equal(t, time.Now().AddDate(0, 0, -7).Format("2006-01-02"), metric.bucket)
	requireWithinDays(t, log.before, 30)
}

func TestRetentionCleanupRejectsInvalidWindow(t *testing.T) {
	svc, _, _, _, _, _ := newRetentionService(retentionOptionRepo{values: map[string]string{
		retentionMetricsDays: "abc",
	}})

	_, err := svc.Cleanup(context.Background())
	require.ErrorIs(t, err, domain.ErrInvalidInput)

	svc, _, _, _, _, _ = newRetentionService(retentionOptionRepo{values: map[string]string{
		retentionMetricsDays: "0",
	}})
	_, err = svc.Cleanup(context.Background())
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

// requireWithinDays allows a minute of drift for the wall clock crossing
// midnight between computing and asserting the cutoff.
func requireWithinDays(t *testing.T, unix uint, days int) {
	t.Helper()
	want := int64(time.Now().AddDate(0, 0, -days).Unix())
	got := int64(unix)
	require.InDelta(t, want, got, 60)
}
