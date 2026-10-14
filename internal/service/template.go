package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type TemplateService struct {
	tx        TxManager
	templates domain.TemplateRepository
}

func NewTemplateService(tx TxManager, templates domain.TemplateRepository) *TemplateService {
	return &TemplateService{tx: tx, templates: templates}
}

func (s *TemplateService) GetByUUID(ctx context.Context, uuid string) (*domain.Template, error) {
	return s.templates.GetByUUID(ctx, uuid)
}

func (s *TemplateService) List(ctx context.Context, cursor string, limit int) ([]*domain.Template, error) {
	return s.templates.List(ctx, cursor, limit)
}

func (s *TemplateService) Create(ctx context.Context, t *domain.Template, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Template.Create(ctx, t); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "template.create", Target: "template:" + t.UUID})
	})
}

func (s *TemplateService) Update(ctx context.Context, t *domain.Template, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Template.Update(ctx, t); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "template.update", Target: "template:" + t.UUID})
	})
}

func (s *TemplateService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Template.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "template.delete", Target: "template:" + uuid})
	})
}
