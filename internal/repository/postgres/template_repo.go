package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.TemplateRepository = (*TemplateRepo)(nil)

type TemplateRepo struct {
	q *Queries
}

func NewTemplateRepo(q *Queries) *TemplateRepo {
	return &TemplateRepo{q: q}
}

func templateFromRow(r MmTemplate) *domain.Template {
	return &domain.Template{
		UUID:        r.Uuid,
		Subject:     r.Subject,
		Body:        r.Body,
		Description: r.Description,
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *TemplateRepo) Create(ctx context.Context, t *domain.Template) error {
	return wrap("template.Create", r.q.CreateTemplate(ctx, CreateTemplateParams{
		Uuid:        t.UUID,
		Subject:     t.Subject,
		Body:        t.Body,
		Description: t.Description,
		Version:     int64(t.Version),
	}))
}

func (r *TemplateRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Template, error) {
	row, err := r.q.GetTemplateByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("template.GetByUUID", err)
	}
	return templateFromRow(row), nil
}

func (r *TemplateRepo) Update(ctx context.Context, t *domain.Template) error {
	return wrap("template.Update", r.q.UpdateTemplate(ctx, UpdateTemplateParams{
		Subject:     t.Subject,
		Body:        t.Body,
		Description: t.Description,
		Version:     int64(t.Version),
		Uuid:        t.UUID,
	}))
}

func (r *TemplateRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("template.DeleteByUUID", r.q.DeleteTemplateByUUID(ctx, uuid))
}

func (r *TemplateRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Template, error) {
	rows, err := r.q.ListTemplates(ctx, ListTemplatesParams{
		Uuid:  cursor,
		Limit: int32(limit),
	})
	if err != nil {
		return nil, wrap("template.List", err)
	}
	out := make([]*domain.Template, 0, len(rows))
	for _, row := range rows {
		out = append(out, templateFromRow(row))
	}
	return out, nil
}
