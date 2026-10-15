package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// EdgeFromDomain converts a domain entity to its legacy model.
func EdgeFromDomain(e *domain.Edge) *models.Edge {
	if e == nil {
		return nil
	}
	return &models.Edge{
		EdgeID:       e.EdgeID,
		UUID:         e.UUID,
		Type:         uint(e.Type),
		AncestorID:   e.AncestorID,
		DescendantID: e.DescendantID,
		Version:      e.Version,
		UpdateAt:     deref(e.UpdateAt),
		CreateAt:     e.CreateAt,
	}
}

// EdgeToDomain converts a legacy model to its domain entity.
func EdgeToDomain(e *models.Edge) *domain.Edge {
	if e == nil {
		return nil
	}
	return &domain.Edge{
		EdgeID:       e.EdgeID,
		UUID:         e.UUID,
		Type:         domain.EdgeType(e.Type),
		AncestorID:   e.AncestorID,
		DescendantID: e.DescendantID,
		Version:      e.Version,
		UpdateAt:     uintPtr(e.UpdateAt),
		CreateAt:     e.CreateAt,
	}
}
