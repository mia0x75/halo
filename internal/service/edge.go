package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

type EdgeService struct {
	tx    TxManager
	edges domain.EdgeRepository
}

func NewEdgeService(tx TxManager, edges domain.EdgeRepository) *EdgeService {
	return &EdgeService{tx: tx, edges: edges}
}

func (s *EdgeService) Get(ctx context.Context, id uint) (*domain.Edge, error) {
	return s.edges.GetByID(ctx, int64(id))
}

func (s *EdgeService) GetByUUID(ctx context.Context, uuid string) (*domain.Edge, error) {
	return s.edges.GetByUUID(ctx, uuid)
}

func (s *EdgeService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Edge, error) {
	return s.edges.List(ctx, cursor, limit)
}

func (s *EdgeService) Create(ctx context.Context, e *domain.Edge, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Edge.Create(ctx, e)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "edge.create", Target: "edge:" + e.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *EdgeService) Update(ctx context.Context, e *domain.Edge, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Edge.Update(ctx, e); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "edge.update", Target: "edge:" + e.UUID})
	})
}

func (s *EdgeService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Edge.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "edge.delete", Target: fmt.Sprintf("edge:%d", id)})
	})
}

// Replace swaps every edge of one kind leaving ancestorID for the given set,
// atomically. It backs the grant/revoke endpoints.
func (s *EdgeService) Replace(ctx context.Context, edgeType uint, ancestorID uint, edges []*domain.Edge, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Edge.DeleteByTypeAncestor(ctx, edgeType, ancestorID); err != nil {
			return err
		}
		now := uint(time.Now().Unix())
		for _, e := range edges {
			e.AncestorID = ancestorID
			if e.UUID == "" {
				e.UUID = uuid.NewString()
			}
			if e.CreateAt == 0 {
				e.CreateAt = now
			}
			if _, err := repos.Edge.Create(ctx, e); err != nil {
				return err
			}
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "edge.replace", Target: fmt.Sprintf("edge:%d:%d", edgeType, ancestorID)})
	})
}
