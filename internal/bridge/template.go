package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// TemplateFromDomain converts a domain entity to its legacy model.
func TemplateFromDomain(t *domain.Template) *models.Template {
	if t == nil {
		return nil
	}
	return &models.Template{
		UUID:        t.UUID,
		Subject:     t.Subject,
		Body:        t.Body,
		Description: t.Description,
		Version:     t.Version,
		UpdateAt:    deref(t.UpdateAt),
		CreateAt:    t.CreateAt,
	}
}

// TemplateToDomain converts a legacy model to its domain entity.
func TemplateToDomain(t *models.Template) *domain.Template {
	if t == nil {
		return nil
	}
	return &domain.Template{
		UUID:        t.UUID,
		Subject:     t.Subject,
		Body:        t.Body,
		Description: t.Description,
		Version:     t.Version,
		UpdateAt:    uintPtr(t.UpdateAt),
		CreateAt:    t.CreateAt,
	}
}
