package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// GroupFromDomain converts a domain entity to its model.
func GroupFromDomain(g *domain.Group) *models.Group {
	if g == nil {
		return nil
	}
	return &models.Group{
		GroupID:     g.GroupID,
		UUID:        g.UUID,
		Name:        g.Name,
		Description: g.Description,
		Status:      g.Status,
		Version:     g.Version,
		UpdateAt:    deref(g.UpdateAt),
		CreateAt:    g.CreateAt,
	}
}

// GroupToDomain converts a model to its domain entity.
func GroupToDomain(g *models.Group) *domain.Group {
	if g == nil {
		return nil
	}
	updateAt := g.UpdateAt
	return &domain.Group{
		GroupID:     g.GroupID,
		UUID:        g.UUID,
		Name:        g.Name,
		Description: g.Description,
		Status:      g.Status,
		Version:     g.Version,
		UpdateAt:    &updateAt,
		CreateAt:    g.CreateAt,
	}
}
