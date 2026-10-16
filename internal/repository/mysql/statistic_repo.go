package mysqlrepo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.StatisticRepository = (*StatisticRepo)(nil)

type StatisticRepo struct {
	q *Queries
}

func NewStatisticRepo(q *Queries) *StatisticRepo {
	return &StatisticRepo{q: q}
}

// statisticFromRow parses the DECIMAL(18,4) column: sqlc reads numeric types as
// string, but the domain exposes float64.
func statisticFromRow(r MmStatistic) (*domain.Statistic, error) {
	value, err := strconv.ParseFloat(r.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("statistic.value %q: %w", r.Value, err)
	}
	return &domain.Statistic{
		Group:    r.Group,
		Key:      r.Key,
		UUID:     r.Uuid,
		Value:    value,
		Version:  int(r.Version),
		UpdateAt: timeFromNull(r.UpdateAt),
		CreateAt: timeToUnix(r.CreateAt),
	}, nil
}

func (r *StatisticRepo) Create(ctx context.Context, s *domain.Statistic) error {
	return wrap("statistic.Create", r.q.CreateStatistic(ctx, CreateStatisticParams{
		Group:   s.Group,
		Key:     s.Key,
		Uuid:    s.UUID,
		Value:   strconv.FormatFloat(s.Value, 'f', -1, 64),
		Version: uint32(s.Version),
	}))
}

func (r *StatisticRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Statistic, error) {
	row, err := r.q.GetStatisticByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("statistic.GetByUUID", err)
	}
	return statisticFromRow(row)
}

// GetByKey resolves the row by its natural composite key.
func (r *StatisticRepo) GetByKey(ctx context.Context, group, key string) (*domain.Statistic, error) {
	row, err := r.q.GetStatisticByKey(ctx, GetStatisticByKeyParams{Group: group, Key: key})
	if err != nil {
		return nil, wrap("statistic.GetByKey", err)
	}
	return statisticFromRow(row)
}

func (r *StatisticRepo) Update(ctx context.Context, s *domain.Statistic) error {
	return wrap("statistic.Update", r.q.UpdateStatistic(ctx, UpdateStatisticParams{
		Value:   strconv.FormatFloat(s.Value, 'f', -1, 64),
		Version: uint32(s.Version),
		Uuid:    s.UUID,
	}))
}

func (r *StatisticRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("statistic.DeleteByUUID", r.q.DeleteStatisticByUUID(ctx, uuid))
}

func (r *StatisticRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Statistic, error) {
	rows, err := r.q.ListStatistics(ctx, ListStatisticsParams{
		Uuid:  cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("statistic.List", err)
	}
	out := make([]*domain.Statistic, 0, len(rows))
	for _, row := range rows {
		s, err := statisticFromRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// UpsertIncrement atomically adds delta to a counter, inserting the row when it
// does not exist yet. Single statement, so concurrent increments never lose an
// update (the previous read-modify-write could).
func (r *StatisticRepo) UpsertIncrement(ctx context.Context, group, key string, delta float64) error {
	return wrap("statistic.UpsertIncrement", r.q.UpsertStatistic(ctx, UpsertStatisticParams{
		Group: group,
		Key:   key,
		Uuid:  uuid.NewString(),
		Value: strconv.FormatFloat(delta, 'f', -1, 64),
	}))
}
