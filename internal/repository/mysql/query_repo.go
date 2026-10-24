package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.QueryRepository = (*QueryRepo)(nil)

type QueryRepo struct {
	q *Queries
}

func NewQueryRepo(q *Queries) *QueryRepo {
	return &QueryRepo{q: q}
}

func queryFromRow(r MmQuery) *domain.Query {
	return &domain.Query{
		QueryID:   uint(r.QueryID),
		UUID:      r.Uuid,
		Type:      domain.QueryType(r.Type),
		ClusterID: uint(r.ClusterID),
		Database:  r.Database,
		Content:   r.Content,
		Plan:      r.Plan,
		UserID:    uint(r.UserID),
		Version:   int(r.Version),
		UpdateAt:  timeFromNull(r.UpdateAt),
		CreateAt:  timeToUnix(r.CreateAt),
	}
}

func (r *QueryRepo) Create(ctx context.Context, q *domain.Query) (int64, error) {
	id, err := r.q.CreateQuery(ctx, CreateQueryParams{
		Uuid:      q.UUID,
		Type:      uint8(q.Type),
		ClusterID: uint32(q.ClusterID),
		Database:  q.Database,
		Content:   q.Content,
		Plan:      q.Plan,
		UserID:    uint32(q.UserID),
		Version:   uint32(q.Version),
	})
	if err != nil {
		return 0, wrap("query.Create", err)
	}
	return id, nil
}

func (r *QueryRepo) GetByID(ctx context.Context, id int64) (*domain.Query, error) {
	row, err := r.q.GetQueryByID(ctx, clampUint32(id))
	if err != nil {
		return nil, wrap("query.GetByID", err)
	}
	return queryFromRow(row), nil
}

func (r *QueryRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Query, error) {
	row, err := r.q.GetQueryByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("query.GetByUUID", err)
	}
	return queryFromRow(row), nil
}

func (r *QueryRepo) Update(ctx context.Context, q *domain.Query) error {
	err := r.q.UpdateQuery(ctx, UpdateQueryParams{
		Uuid:      q.UUID,
		Type:      uint8(q.Type),
		ClusterID: uint32(q.ClusterID),
		Database:  q.Database,
		Content:   q.Content,
		Plan:      q.Plan,
		UserID:    uint32(q.UserID),
		Version:   uint32(q.Version),
		QueryID:   uint32(q.QueryID),
	})
	return wrap("query.Update", err)
}

func (r *QueryRepo) Delete(ctx context.Context, id int64) error {
	return wrap("query.Delete", r.q.DeleteQueryByID(ctx, clampUint32(id)))
}

func (r *QueryRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("query.DeleteByUUID", r.q.DeleteQueryByUUID(ctx, uuid))
}

func (r *QueryRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Query, error) {
	rows, err := r.q.ListQueries(ctx, ListQueriesParams{
		QueryID: clampUint32(cursor),
		Limit:   clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("query.List", err)
	}
	out := make([]*domain.Query, 0, len(rows))
	for _, row := range rows {
		out = append(out, queryFromRow(row))
	}
	return out, nil
}

// ListByUser returns every query logged for one user, oldest first.
func (r *QueryRepo) ListByUser(ctx context.Context, userID uint) ([]*domain.Query, error) {
	rows, err := r.q.ListQueriesByUser(ctx, clampUint(userID))
	if err != nil {
		return nil, wrap("query.ListByUser", err)
	}
	out := make([]*domain.Query, 0, len(rows))
	for _, row := range rows {
		out = append(out, queryFromRow(row))
	}
	return out, nil
}

// Count returns the total number of queries.
func (r *QueryRepo) Count(ctx context.Context) (int64, error) {
	n, err := r.q.CountQueries(ctx)
	if err != nil {
		return 0, wrap("query.Count", err)
	}
	return n, nil
}
