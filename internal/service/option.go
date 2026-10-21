package service

import (
	"context"
	"fmt"

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

// Update stores a new value for an option and records the old and new value in
// the audit log.
func (s *OptionService) Update(ctx context.Context, o *domain.Option, operatorID uint) error {
	if err := domain.ValidateOptionValue(o.Name, o.Parameter); err != nil {
		return err
	}
	current, err := s.options.GetByUUID(ctx, o.UUID)
	if err != nil {
		return err
	}

	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Option.Update(ctx, o); err != nil {
			return err
		}
		detail := fmt.Sprintf("%s->%s",
			domain.MaskOptionValue(current.Name, current.Parameter),
			domain.MaskOptionValue(o.Name, o.Parameter))
		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    operatorID,
			Operation: "option.update",
			Target:    "option:" + o.UUID,
			Detail:    detail,
		})
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
