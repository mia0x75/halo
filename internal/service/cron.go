package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

type CronService struct {
	tx    TxManager
	crons domain.CronRepository
}

func NewCronService(tx TxManager, crons domain.CronRepository) *CronService {
	return &CronService{tx: tx, crons: crons}
}

func (s *CronService) Get(ctx context.Context, id uint) (*domain.Cron, error) {
	return s.crons.GetByID(ctx, int64(id))
}

func (s *CronService) GetByUUID(ctx context.Context, uuid string) (*domain.Cron, error) {
	return s.crons.GetByUUID(ctx, uuid)
}

func (s *CronService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Cron, error) {
	return s.crons.List(ctx, cursor, limit)
}

func (s *CronService) Create(ctx context.Context, c *domain.Cron, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Cron.Create(ctx, c)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "cron.create", Target: "cron:" + c.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CronService) Update(ctx context.Context, c *domain.Cron, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Cron.Update(ctx, c); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "cron.update", Target: "cron:" + c.UUID})
	})
}

func (s *CronService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Cron.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "cron.delete", Target: fmt.Sprintf("cron:%d", id)})
	})
}
