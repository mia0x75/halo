package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
)

// AvatarFromDomain converts a domain entity to its legacy model.
func AvatarFromDomain(a *domain.Avatar) *models.Avatar {
	if a == nil {
		return nil
	}
	return &models.Avatar{
		AvatarID: a.AvatarID,
		UUID:     a.UUID,
		URL:      a.URL,
		Version:  a.Version,
		UpdateAt: deref(a.UpdateAt),
		CreateAt: a.CreateAt,
	}
}

// AvatarToDomain converts a legacy model to its domain entity.
func AvatarToDomain(a *models.Avatar) *domain.Avatar {
	if a == nil {
		return nil
	}
	return &domain.Avatar{
		AvatarID: a.AvatarID,
		UUID:     a.UUID,
		URL:      a.URL,
		Version:  a.Version,
		UpdateAt: uintPtr(a.UpdateAt),
		CreateAt: a.CreateAt,
	}
}
