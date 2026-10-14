package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

type UserService struct {
	tx    TxManager
	users domain.UserRepository
}

func NewUserService(tx TxManager, users domain.UserRepository) *UserService {
	return &UserService{tx: tx, users: users}
}

func (s *UserService) Get(ctx context.Context, id uint) (*domain.User, error) {
	return s.users.GetByID(ctx, int64(id))
}

func (s *UserService) GetByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	return s.users.GetByUUID(ctx, uuid)
}

func (s *UserService) List(ctx context.Context, cursor int64, limit int) ([]*domain.User, error) {
	return s.users.List(ctx, cursor, limit)
}

func (s *UserService) Create(ctx context.Context, u *domain.User, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.User.Create(ctx, u)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "user.create", Target: "user:" + u.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserService) Update(ctx context.Context, u *domain.User, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.User.Update(ctx, u); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "user.update", Target: "user:" + u.UUID})
	})
}

func (s *UserService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.User.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "user.delete", Target: fmt.Sprintf("user:%d", id)})
	})
}
