package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

type QueryService struct {
	tx      TxManager
	queries domain.QueryRepository
}

func NewQueryService(tx TxManager, queries domain.QueryRepository) *QueryService {
	return &QueryService{tx: tx, queries: queries}
}

func (s *QueryService) Get(ctx context.Context, id uint) (*domain.Query, error) {
	return s.queries.GetByID(ctx, int64(id))
}

func (s *QueryService) GetByUUID(ctx context.Context, uuid string) (*domain.Query, error) {
	return s.queries.GetByUUID(ctx, uuid)
}

func (s *QueryService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Query, error) {
	return s.queries.List(ctx, cursor, limit)
}

// Count returns the total number of queries, independent of any page cursor.
func (s *QueryService) Count(ctx context.Context) (int64, error) {
	return s.queries.Count(ctx)
}

// ListByUser returns every query logged for one user.
func (s *QueryService) ListByUser(ctx context.Context, userID uint) ([]*domain.Query, error) {
	return s.queries.ListByUser(ctx, userID)
}

func (s *QueryService) Create(ctx context.Context, q *domain.Query, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Query.Create(ctx, q)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "query.create", Target: "query:" + q.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *QueryService) Update(ctx context.Context, q *domain.Query, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Query.Update(ctx, q); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "query.update", Target: "query:" + q.UUID})
	})
}

func (s *QueryService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Query.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "query.delete", Target: fmt.Sprintf("query:%d", id)})
	})
}
