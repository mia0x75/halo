package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type MetricService struct {
	tx      TxManager
	metric  domain.MetricRepository
	tickets domain.TicketRepository
	hits    domain.RuleHitRepository
}

func NewMetricService(tx TxManager, metric domain.MetricRepository, tickets domain.TicketRepository, hits domain.RuleHitRepository) *MetricService {
	return &MetricService{tx: tx, metric: metric, tickets: tickets, hits: hits}
}

func (s *MetricService) GetByUUID(ctx context.Context, uuid string) (*domain.Metric, error) {
	return s.metric.GetByUUID(ctx, uuid)
}

func (s *MetricService) GetByKey(ctx context.Context, key domain.MetricKey) (*domain.Metric, error) {
	return s.metric.GetByKey(ctx, key)
}

func (s *MetricService) List(ctx context.Context, cursor string, limit int) ([]*domain.Metric, error) {
	return s.metric.List(ctx, cursor, limit)
}

func (s *MetricService) Create(ctx context.Context, m *domain.Metric, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Metric.Create(ctx, m); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "metric.create", Target: "metric:" + m.UUID})
	})
}

func (s *MetricService) Update(ctx context.Context, m *domain.Metric, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Metric.Update(ctx, m); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "metric.update", Target: "metric:" + m.UUID})
	})
}

func (s *MetricService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Metric.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "metric.delete", Target: "metric:" + uuid})
	})
}
