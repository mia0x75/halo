package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.CronRepository = (*CronRepo)(nil)

type CronRepo struct {
	q *Queries
}

func NewCronRepo(q *Queries) *CronRepo {
	return &CronRepo{q: q}
}

func cronFromRow(r MmCron) *domain.Cron {
	return &domain.Cron{
		CronID:    uint(r.CronID),
		UUID:      r.Uuid,
		Status:    domain.CronStatus(r.Status),
		Name:      r.Name,
		Cmd:       r.Cmd,
		Params:    r.Params,
		LastRun:   stringFromNull(r.LastRun),
		NextRun:   stringFromNull(r.NextRun),
		Interval:  stringFromNull(r.Interval),
		Duration:  stringFromNull(r.Duration),
		Recurrent: uint8FromNull(r.Recurrent),
		Hash:      stringFromNull(r.Hash),
		Version:   int(r.Version),
		UpdateAt:  uintFromNull(r.UpdateAt),
		CreateAt:  uint(r.CreateAt),
	}
}

func (r *CronRepo) Create(ctx context.Context, c *domain.Cron) (int64, error) {
	id, err := r.q.CreateCron(ctx, CreateCronParams{
		Uuid:      c.UUID,
		Status:    string(c.Status),
		Name:      c.Name,
		Cmd:       c.Cmd,
		Params:    c.Params,
		LastRun:   nullFromString(c.LastRun),
		NextRun:   nullFromString(c.NextRun),
		Interval:  nullFromString(c.Interval),
		Duration:  nullFromString(c.Duration),
		Recurrent: nullFromUint8(c.Recurrent),
		Hash:      nullFromString(c.Hash),
		Version:   uint32(c.Version),
		UpdateAt:  nullFromUint(c.UpdateAt),
		CreateAt:  uint32(c.CreateAt),
	})
	if err != nil {
		return 0, wrap("cron.Create", err)
	}
	return id, nil
}

func (r *CronRepo) GetByID(ctx context.Context, id int64) (*domain.Cron, error) {
	row, err := r.q.GetCronByID(ctx, uint32(id))
	if err != nil {
		return nil, wrap("cron.GetByID", err)
	}
	return cronFromRow(row), nil
}

func (r *CronRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Cron, error) {
	row, err := r.q.GetCronByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("cron.GetByUUID", err)
	}
	return cronFromRow(row), nil
}

func (r *CronRepo) Update(ctx context.Context, c *domain.Cron) error {
	err := r.q.UpdateCron(ctx, UpdateCronParams{
		Uuid:      c.UUID,
		Status:    string(c.Status),
		Name:      c.Name,
		Cmd:       c.Cmd,
		Params:    c.Params,
		LastRun:   nullFromString(c.LastRun),
		NextRun:   nullFromString(c.NextRun),
		Interval:  nullFromString(c.Interval),
		Duration:  nullFromString(c.Duration),
		Recurrent: nullFromUint8(c.Recurrent),
		Hash:      nullFromString(c.Hash),
		Version:   uint32(c.Version),
		UpdateAt:  nullFromUint(c.UpdateAt),
		CronID:    uint32(c.CronID),
	})
	return wrap("cron.Update", err)
}

func (r *CronRepo) Delete(ctx context.Context, id int64) error {
	return wrap("cron.Delete", r.q.DeleteCronByID(ctx, uint32(id)))
}

func (r *CronRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("cron.DeleteByUUID", r.q.DeleteCronByUUID(ctx, uuid))
}

func (r *CronRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Cron, error) {
	rows, err := r.q.ListCrons(ctx, ListCronsParams{
		CronID: uint32(cursor),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, wrap("cron.List", err)
	}
	out := make([]*domain.Cron, 0, len(rows))
	for _, row := range rows {
		out = append(out, cronFromRow(row))
	}
	return out, nil
}
