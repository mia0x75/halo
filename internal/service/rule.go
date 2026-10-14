package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type RuleService struct {
	tx    TxManager
	rules domain.RuleRepository
}

func NewRuleService(tx TxManager, rules domain.RuleRepository) *RuleService {
	return &RuleService{tx: tx, rules: rules}
}

func (s *RuleService) GetByUUID(ctx context.Context, uuid string) (*domain.Rule, error) {
	return s.rules.GetByUUID(ctx, uuid)
}

func (s *RuleService) GetByName(ctx context.Context, name string) (*domain.Rule, error) {
	return s.rules.GetByName(ctx, name)
}

func (s *RuleService) List(ctx context.Context, cursor string, limit int) ([]*domain.Rule, error) {
	return s.rules.List(ctx, cursor, limit)
}

func (s *RuleService) Create(ctx context.Context, r *domain.Rule, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Rule.Create(ctx, r); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "rule.create", Target: "rule:" + r.UUID})
	})
}

func (s *RuleService) Update(ctx context.Context, r *domain.Rule, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Rule.Update(ctx, r); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "rule.update", Target: "rule:" + r.UUID})
	})
}

func (s *RuleService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Rule.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "rule.delete", Target: "rule:" + uuid})
	})
}
