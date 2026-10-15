package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// RoleFromDomain converts a domain entity to its legacy model.
func RoleFromDomain(r *domain.Role) *models.Role {
	if r == nil {
		return nil
	}
	return &models.Role{
		RoleID:      r.RoleID,
		UUID:        r.UUID,
		Name:        r.Name,
		Description: r.Description,
		Version:     r.Version,
		UpdateAt:    deref(r.UpdateAt),
		CreateAt:    r.CreateAt,
	}
}

// RoleToDomain converts a legacy model to its domain entity.
