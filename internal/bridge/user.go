package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
)

// UserFromDomain converts a domain entity to its legacy model.
func UserFromDomain(u *domain.User) *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		UserID:   u.UserID,
		UUID:     u.UUID,
		Email:    u.Email,
		Password: u.Password,
		Status:   uint8(u.Status),
		Name:     u.Name,
		Phone:    u.Phone,
		AvatarID: u.AvatarID,
		Version:  u.Version,
		UpdateAt: deref(u.UpdateAt),
		CreateAt: u.CreateAt,
	}
}

// UserToDomain converts a legacy model to its domain entity.
func UserToDomain(u *models.User) *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		UserID:   u.UserID,
		UUID:     u.UUID,
		Email:    u.Email,
		Password: u.Password,
		Status:   domain.UserStatus(u.Status),
		Name:     u.Name,
		Phone:    u.Phone,
		AvatarID: u.AvatarID,
		Version:  u.Version,
		UpdateAt: uintPtr(u.UpdateAt),
		CreateAt: u.CreateAt,
	}
}
