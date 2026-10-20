package mysqlrepo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.MetricRepository = (*MetricRepo)(nil)

type MetricRepo struct {
	q *Queries
}

func NewMetricRepo(q *Queries) *MetricRepo {
	return &MetricRepo{q: q}
}

// metricFromRow parses the DECIMAL(18,4) column: sqlc reads numeric types as
// string, but the domain exposes float64.
func metricFromRow(r MmMetric) (*domain.Metric, error) {
	value, err := strconv.ParseFloat(r.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("metric.value %q: %w", r.Value, err)
	}
	return &domain.Metric{
		MetricID: uint(r.MetricID),
		UUID:     r.Uuid,
		Name:     r.Name,
		Scope:    r.Scope,
		ScopeID:  r.ScopeID,
		Type:     r.Type,
		Bucket:   r.Bucket,
		Value:    value,
		Version:  int(r.Version),
		UpdateAt: timeFromNull(r.UpdateAt),
		CreateAt: timeToUnix(r.CreateAt),
	}, nil
}

func (r *MetricRepo) Create(ctx context.Context, m *domain.Metric) error {
	return wrap("metric.Create", r.q.CreateMetric(ctx, CreateMetricParams{
		Uuid:    m.UUID,
		Name:    m.Name,
		Scope:   m.Scope,
		ScopeID: m.ScopeID,
		Type:    m.Type,
		Bucket:  m.Bucket,
		Value:   strconv.FormatFloat(m.Value, 'f', -1, 64),
		Version: uint32(m.Version),
	}))
}

func (r *MetricRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Metric, error) {
	row, err := r.q.GetMetricByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("metric.GetByUUID", err)
	}
	return metricFromRow(row)
}

// GetByKey resolves the row by its natural composite key.
func (r *MetricRepo) GetByKey(ctx context.Context, key domain.MetricKey) (*domain.Metric, error) {
	row, err := r.q.GetMetricByKey(ctx, GetMetricByKeyParams{
		Name:    key.Name,
		Scope:   key.Scope,
		ScopeID: key.ScopeID,
		Bucket:  key.Bucket,
	})
	if err != nil {
		return nil, wrap("metric.GetByKey", err)
	}
	return metricFromRow(row)
}

func (r *MetricRepo) Update(ctx context.Context, m *domain.Metric) error {
	return wrap("metric.Update", r.q.UpdateMetric(ctx, UpdateMetricParams{
		Value:   strconv.FormatFloat(m.Value, 'f', -1, 64),
		Version: uint32(m.Version),
		Uuid:    m.UUID,
	}))
}

func (r *MetricRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("metric.DeleteByUUID", r.q.DeleteMetricByUUID(ctx, uuid))
}

func (r *MetricRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Metric, error) {
	rows, err := r.q.ListMetrics(ctx, ListMetricsParams{
		Uuid:  cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("metric.List", err)
	}
	out := make([]*domain.Metric, 0, len(rows))
	for _, row := range rows {
		m, err := metricFromRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

// Increment atomically adds delta to a counter metric, inserting the row when
// it does not exist yet. Single statement, so concurrent increments never lose
// an update.
func (r *MetricRepo) Increment(ctx context.Context, key domain.MetricKey, delta float64) error {
	return wrap("metric.Increment", r.q.IncrementMetric(ctx, IncrementMetricParams{
		Uuid:    uuid.NewString(),
		Name:    key.Name,
		Scope:   key.Scope,
		ScopeID: key.ScopeID,
		Bucket:  key.Bucket,
		Value:   strconv.FormatFloat(delta, 'f', -1, 64),
	}))
}

// SetGauge overwrites the value of a gauge metric.
func (r *MetricRepo) SetGauge(ctx context.Context, key domain.MetricKey, value float64) error {
	return wrap("metric.SetGauge", r.q.SetGaugeMetric(ctx, SetGaugeMetricParams{
		Uuid:    uuid.NewString(),
		Name:    key.Name,
		Scope:   key.Scope,
		ScopeID: key.ScopeID,
		Bucket:  key.Bucket,
		Value:   strconv.FormatFloat(value, 'f', -1, 64),
	}))
}

// DeleteBucketsBefore removes the daily buckets older than the given date.
func (r *MetricRepo) DeleteBucketsBefore(ctx context.Context, bucket string) (int64, error) {
	return r.q.DeleteMetricBucketsBefore(ctx, bucket)
}
