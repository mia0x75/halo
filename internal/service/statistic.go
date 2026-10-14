package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type StatisticService struct {
	tx         TxManager
	statistics domain.StatisticRepository
}

func NewStatisticService(tx TxManager, statistics domain.StatisticRepository) *StatisticService {
	return &StatisticService{tx: tx, statistics: statistics}
}

func (s *StatisticService) GetByUUID(ctx context.Context, uuid string) (*domain.Statistic, error) {
	return s.statistics.GetByUUID(ctx, uuid)
}

func (s *StatisticService) GetByKey(ctx context.Context, group, key string) (*domain.Statistic, error) {
	return s.statistics.GetByKey(ctx, group, key)
}

func (s *StatisticService) List(ctx context.Context, cursor string, limit int) ([]*domain.Statistic, error) {
	return s.statistics.List(ctx, cursor, limit)
}

func (s *StatisticService) Create(ctx context.Context, st *domain.Statistic, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Statistic.Create(ctx, st); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "statistic.create", Target: "statistic:" + st.UUID})
	})
}

func (s *StatisticService) Update(ctx context.Context, st *domain.Statistic, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Statistic.Update(ctx, st); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "statistic.update", Target: "statistic:" + st.UUID})
	})
}

func (s *StatisticService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Statistic.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "statistic.delete", Target: "statistic:" + uuid})
	})
}
