package postgresrepo

import (
	"context"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.LogRepository = (*LogRepo)(nil)

type LogRepo struct {
	q *Queries
}

func NewLogRepo(q *Queries) *LogRepo {
	return &LogRepo{q: q}
}

func logFromRow(r MmLog) *domain.Log {
	return &domain.Log{
		LogID:     uint(r.LogID),
		UUID:      r.Uuid,
		UserID:    uint(r.UserID),
		Operation: r.Operation,
		Version:   int(r.Version),
		CreateAt:  timeToUnix(r.CreateAt),
	}
}

func (r *LogRepo) Create(ctx context.Context, l *domain.Log) (int64, error) {
	id, err := r.q.CreateLog(ctx, CreateLogParams{
		Uuid:      l.UUID,
		UserID:    int64(l.UserID),
		Operation: l.Operation,
		Version:   int64(l.Version),
	})
	if err != nil {
		return 0, wrap("log.Create", err)
	}
	return id, nil
}

func (r *LogRepo) GetByID(ctx context.Context, id int64) (*domain.Log, error) {
	row, err := r.q.GetLogByID(ctx, id)
	if err != nil {
		return nil, wrap("log.GetByID", err)
	}
	return logFromRow(row), nil
}

func (r *LogRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Log, error) {
	row, err := r.q.GetLogByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("log.GetByUUID", err)
	}
	return logFromRow(row), nil
}

func (r *LogRepo) Update(ctx context.Context, l *domain.Log) error {
	err := r.q.UpdateLog(ctx, UpdateLogParams{
		Uuid:      l.UUID,
		UserID:    int64(l.UserID),
		Operation: l.Operation,
		Version:   int64(l.Version),
		LogID:     int64(l.LogID),
	})
	return wrap("log.Update", err)
}

func (r *LogRepo) Delete(ctx context.Context, id int64) error {
	return wrap("log.Delete", r.q.DeleteLogByID(ctx, id))
}

func (r *LogRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("log.DeleteByUUID", r.q.DeleteLogByUUID(ctx, uuid))
}

// DeleteBefore removes the log rows created before the given time.
func (r *LogRepo) DeleteBefore(ctx context.Context, before uint) (int64, error) {
	rows, err := r.q.DeleteLogsBefore(ctx, time.Unix(int64(before), 0).UTC())
	if err != nil {
		return 0, wrap("log.DeleteBefore", err)
	}
	return rows, nil
}

func (r *LogRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Log, error) {
	rows, err := r.q.ListLogs(ctx, ListLogsParams{
		LogID: cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("log.List", err)
	}
	out := make([]*domain.Log, 0, len(rows))
	for _, row := range rows {
		out = append(out, logFromRow(row))
	}
	return out, nil
}
