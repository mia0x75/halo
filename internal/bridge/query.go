package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// QueryFromDomain converts a domain entity to its legacy model.
func QueryFromDomain(q *domain.Query) *models.Query {
	if q == nil {
		return nil
	}
	return &models.Query{
		QueryID:   q.QueryID,
		UUID:      q.UUID,
		Type:      uint8(q.Type),
		ClusterID: q.ClusterID,
		Database:  q.Database,
		Content:   q.Content,
		Plan:      q.Plan,
		UserID:    q.UserID,
		Version:   q.Version,
		UpdateAt:  deref(q.UpdateAt),
		CreateAt:  q.CreateAt,
	}
}

// QueryToDomain converts a legacy model to its domain entity.
func QueryToDomain(q *models.Query) *domain.Query {
	if q == nil {
		return nil
	}
	return &domain.Query{
		QueryID:   q.QueryID,
		UUID:      q.UUID,
		Type:      domain.QueryType(q.Type),
		ClusterID: q.ClusterID,
		Database:  q.Database,
		Content:   q.Content,
		Plan:      q.Plan,
		UserID:    q.UserID,
		Version:   q.Version,
		UpdateAt:  uintPtr(q.UpdateAt),
		CreateAt:  q.CreateAt,
	}
}
