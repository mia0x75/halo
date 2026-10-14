package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

type RoleService struct {
	tx    TxManager
	roles domain.RoleRepository
}

func NewRoleService(tx TxManager, roles domain.RoleRepository) *RoleService {
	return &RoleService{tx: tx, roles: roles}
}

func (s *RoleService) Get(ctx context.Context, id uint) (*domain.Role, error) {
	return s.roles.GetByID(ctx, int64(id))
}

func (s *RoleService) GetByUUID(ctx context.Context, uuid string) (*domain.Role, error) {
	return s.roles.GetByUUID(ctx, uuid)
}

func (s *RoleService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Role, error) {
	return s.roles.List(ctx, cursor, limit)
}

func (s *RoleService) Create(ctx context.Context, r *domain.Role, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Role.Create(ctx, r)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "role.create", Target: "role:" + r.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *RoleService) Update(ctx context.Context, r *domain.Role, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Role.Update(ctx, r); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "role.update", Target: "role:" + r.UUID})
	})
}

func (s *RoleService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Role.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "role.delete", Target: fmt.Sprintf("role:%d", id)})
	})
}
