package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type OptionService struct {
	tx      TxManager
	options domain.OptionRepository
}

func NewOptionService(tx TxManager, options domain.OptionRepository) *OptionService {
	return &OptionService{tx: tx, options: options}
}

func (s *OptionService) GetByUUID(ctx context.Context, uuid string) (*domain.Option, error) {
	return s.options.GetByUUID(ctx, uuid)
}

func (s *OptionService) GetByName(ctx context.Context, name string) (*domain.Option, error) {
	return s.options.GetByName(ctx, name)
}

func (s *OptionService) List(ctx context.Context, cursor string, limit int) ([]*domain.Option, error) {
	return s.options.List(ctx, cursor, limit)
}

func (s *OptionService) Create(ctx context.Context, o *domain.Option, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Option.Create(ctx, o); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "option.create", Target: "option:" + o.UUID})
	})
}

func (s *OptionService) Update(ctx context.Context, o *domain.Option, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Option.Update(ctx, o); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "option.update", Target: "option:" + o.UUID})
	})
}

func (s *OptionService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Option.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "option.delete", Target: "option:" + uuid})
	})
}
