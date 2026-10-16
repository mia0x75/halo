package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.OptionRepository = (*OptionRepo)(nil)

type OptionRepo struct {
	q *Queries
}

func NewOptionRepo(q *Queries) *OptionRepo {
	return &OptionRepo{q: q}
}

func optionFromRow(r MmOption) *domain.Option {
	return &domain.Option{
		Name:        r.Name,
		UUID:        r.Uuid,
		Parameter:   r.Parameter,
		Description: r.Description,
		Element:     r.Element,
		Writable:    uint8(r.Writable),
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *OptionRepo) Create(ctx context.Context, o *domain.Option) error {
	return wrap("option.Create", r.q.CreateOption(ctx, CreateOptionParams{
		Name:        o.Name,
		Uuid:        o.UUID,
		Parameter:   o.Parameter,
		Description: o.Description,
		Element:     o.Element,
		Writable:    int16(o.Writable),
		Version:     int64(o.Version),
	}))
}

func (r *OptionRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Option, error) {
	row, err := r.q.GetOptionByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("option.GetByUUID", err)
	}
	return optionFromRow(row), nil
}

// GetByName resolves the row by its natural key.
func (r *OptionRepo) GetByName(ctx context.Context, name string) (*domain.Option, error) {
	row, err := r.q.GetOptionByName(ctx, name)
	if err != nil {
		return nil, wrap("option.GetByName", err)
	}
	return optionFromRow(row), nil
}

func (r *OptionRepo) Update(ctx context.Context, o *domain.Option) error {
	return wrap("option.Update", r.q.UpdateOption(ctx, UpdateOptionParams{
		Parameter:   o.Parameter,
		Description: o.Description,
		Element:     o.Element,
		Writable:    int16(o.Writable),
		Version:     int64(o.Version),
		Uuid:        o.UUID,
	}))
}

func (r *OptionRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("option.DeleteByUUID", r.q.DeleteOptionByUUID(ctx, uuid))
}

func (r *OptionRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Option, error) {
	rows, err := r.q.ListOptions(ctx, ListOptionsParams{
		Uuid:  cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("option.List", err)
	}
	out := make([]*domain.Option, 0, len(rows))
	for _, row := range rows {
		out = append(out, optionFromRow(row))
	}
	return out, nil
}
