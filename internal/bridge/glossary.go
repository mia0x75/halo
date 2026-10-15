package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// GlossaryFromDomain converts a domain entity to its legacy model.
func GlossaryFromDomain(g *domain.Glossary) *models.Glossary {
	if g == nil {
		return nil
	}
	return &models.Glossary{
		Group:       g.Group,
		Key:         uint(g.Key),
		Value:       g.Value,
		UUID:        g.UUID,
		Description: g.Description,
		Version:     g.Version,
		UpdateAt:    deref(g.UpdateAt),
		CreateAt:    g.CreateAt,
	}
}

// GlossaryToDomain converts a legacy model to its domain entity.
