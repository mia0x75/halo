package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type GlossaryService struct {
	tx         TxManager
	glossaries domain.GlossaryRepository
}

func NewGlossaryService(tx TxManager, glossaries domain.GlossaryRepository) *GlossaryService {
	return &GlossaryService{tx: tx, glossaries: glossaries}
}

func (s *GlossaryService) GetByUUID(ctx context.Context, uuid string) (*domain.Glossary, error) {
	return s.glossaries.GetByUUID(ctx, uuid)
}

func (s *GlossaryService) GetByKey(ctx context.Context, group string, key uint8) (*domain.Glossary, error) {
	return s.glossaries.GetByKey(ctx, group, key)
}

func (s *GlossaryService) List(ctx context.Context, cursor string, limit int) ([]*domain.Glossary, error) {
	return s.glossaries.List(ctx, cursor, limit)
}

func (s *GlossaryService) Create(ctx context.Context, g *domain.Glossary, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Glossary.Create(ctx, g); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "glossary.create", Target: "glossary:" + g.UUID})
	})
}

func (s *GlossaryService) Update(ctx context.Context, g *domain.Glossary, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Glossary.Update(ctx, g); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "glossary.update", Target: "glossary:" + g.UUID})
	})
}

func (s *GlossaryService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Glossary.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "glossary.delete", Target: "glossary:" + uuid})
	})
}

func (s *GlossaryService) DeleteByKey(ctx context.Context, group string, key uint8, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Glossary.DeleteByKey(ctx, group, key); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "glossary.delete"})
	})
}
