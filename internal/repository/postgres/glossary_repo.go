package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.GlossaryRepository = (*GlossaryRepo)(nil)

type GlossaryRepo struct {
	q *Queries
}

func NewGlossaryRepo(q *Queries) *GlossaryRepo {
	return &GlossaryRepo{q: q}
}

func glossaryFromRow(r MmGlossary) *domain.Glossary {
	return &domain.Glossary{
		Group:       r.Group,
		Key:         uint8(r.Key),
		Value:       r.Value,
		UUID:        r.Uuid,
		Description: r.Description,
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *GlossaryRepo) Create(ctx context.Context, g *domain.Glossary) error {
	return wrap("glossary.Create", r.q.CreateGlossary(ctx, CreateGlossaryParams{
		Group:       g.Group,
		Key:         int16(g.Key),
		Value:       g.Value,
		Uuid:        g.UUID,
		Description: g.Description,
		Version:     int64(g.Version),
	}))
}

func (r *GlossaryRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Glossary, error) {
	row, err := r.q.GetGlossaryByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("glossary.GetByUUID", err)
	}
	return glossaryFromRow(row), nil
}

// GetByKey resolves the row by its natural composite key.
func (r *GlossaryRepo) GetByKey(ctx context.Context, group string, key uint8) (*domain.Glossary, error) {
	row, err := r.q.GetGlossaryByKey(ctx, GetGlossaryByKeyParams{Group: group, Key: int16(key)})
	if err != nil {
		return nil, wrap("glossary.GetByKey", err)
	}
	return glossaryFromRow(row), nil
}

func (r *GlossaryRepo) Update(ctx context.Context, g *domain.Glossary) error {
	return wrap("glossary.Update", r.q.UpdateGlossary(ctx, UpdateGlossaryParams{
		Value:       g.Value,
		Description: g.Description,
		Version:     int64(g.Version),
		Uuid:        g.UUID,
	}))
}

func (r *GlossaryRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("glossary.DeleteByUUID", r.q.DeleteGlossaryByUUID(ctx, uuid))
}

func (r *GlossaryRepo) DeleteByKey(ctx context.Context, group string, key uint8) error {
	return wrap("glossary.DeleteByKey", r.q.DeleteGlossaryByKey(ctx, DeleteGlossaryByKeyParams{Group: group, Key: int16(key)}))
}

func (r *GlossaryRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Glossary, error) {
	rows, err := r.q.ListGlossaries(ctx, ListGlossariesParams{
		Uuid:  cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("glossary.List", err)
	}
	out := make([]*domain.Glossary, 0, len(rows))
	for _, row := range rows {
		out = append(out, glossaryFromRow(row))
	}
	return out, nil
}
