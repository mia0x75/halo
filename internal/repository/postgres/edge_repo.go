package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.EdgeRepository = (*EdgeRepo)(nil)

type EdgeRepo struct {
	q *Queries
}

func NewEdgeRepo(q *Queries) *EdgeRepo {
	return &EdgeRepo{q: q}
}

func edgeFromRow(r MmEdge) *domain.Edge {
	return &domain.Edge{
		EdgeID:       uint(r.EdgeID),
		UUID:         r.Uuid,
		Type:         domain.EdgeType(r.Type),
		AncestorID:   uint(r.AncestorID),
		DescendantID: uint(r.DescendantID),
		Version:      int(r.Version),
		UpdateAt:     timeFromNull(r.UpdateAt),
		CreateAt:     timeToUnix(r.CreateAt),
	}
}

func (r *EdgeRepo) Create(ctx context.Context, e *domain.Edge) (int64, error) {
	id, err := r.q.CreateEdge(ctx, CreateEdgeParams{
		Uuid:         e.UUID,
		Type:         int64(e.Type),
		AncestorID:   int64(e.AncestorID),
		DescendantID: int64(e.DescendantID),
		Version:      int64(e.Version),
	})
	if err != nil {
		return 0, wrap("edge.Create", err)
	}
	return id, nil
}

func (r *EdgeRepo) GetByID(ctx context.Context, id int64) (*domain.Edge, error) {
	row, err := r.q.GetEdgeByID(ctx, id)
	if err != nil {
		return nil, wrap("edge.GetByID", err)
	}
	return edgeFromRow(row), nil
}

func (r *EdgeRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Edge, error) {
	row, err := r.q.GetEdgeByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("edge.GetByUUID", err)
	}
	return edgeFromRow(row), nil
}

func (r *EdgeRepo) Update(ctx context.Context, e *domain.Edge) error {
	err := r.q.UpdateEdge(ctx, UpdateEdgeParams{
		Uuid:         e.UUID,
		Type:         int64(e.Type),
		AncestorID:   int64(e.AncestorID),
		DescendantID: int64(e.DescendantID),
		Version:      int64(e.Version),
		EdgeID:       int64(e.EdgeID),
	})
	return wrap("edge.Update", err)
}

func (r *EdgeRepo) Delete(ctx context.Context, id int64) error {
	return wrap("edge.Delete", r.q.DeleteEdgeByID(ctx, id))
}

func (r *EdgeRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("edge.DeleteByUUID", r.q.DeleteEdgeByUUID(ctx, uuid))
}

func (r *EdgeRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Edge, error) {
	rows, err := r.q.ListEdges(ctx, ListEdgesParams{
		EdgeID: cursor,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, wrap("edge.List", err)
	}
	out := make([]*domain.Edge, 0, len(rows))
	for _, row := range rows {
		out = append(out, edgeFromRow(row))
	}
	return out, nil
}

// DeleteByTypeAncestor removes every edge of one kind leaving an ancestor, the
// primitive behind "replace this user's grants".
func (r *EdgeRepo) DeleteByTypeAncestor(ctx context.Context, edgeType uint, ancestorID uint) error {
	return wrap("edge.DeleteByTypeAncestor", r.q.DeleteEdgesByTypeAncestor(ctx, DeleteEdgesByTypeAncestorParams{
		Type:       int64(edgeType),
		AncestorID: int64(ancestorID),
	}))
}
