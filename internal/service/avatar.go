package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

type AvatarService struct {
	tx      TxManager
	avatars domain.AvatarRepository
}

func NewAvatarService(tx TxManager, avatars domain.AvatarRepository) *AvatarService {
	return &AvatarService{tx: tx, avatars: avatars}
}

func (s *AvatarService) Get(ctx context.Context, id uint) (*domain.Avatar, error) {
	return s.avatars.GetByID(ctx, int64(id))
}

func (s *AvatarService) GetByUUID(ctx context.Context, uuid string) (*domain.Avatar, error) {
	return s.avatars.GetByUUID(ctx, uuid)
}

func (s *AvatarService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Avatar, error) {
	return s.avatars.List(ctx, cursor, limit)
}

func (s *AvatarService) Create(ctx context.Context, a *domain.Avatar, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Avatar.Create(ctx, a)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "avatar.create", Target: "avatar:" + a.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AvatarService) Update(ctx context.Context, a *domain.Avatar, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Avatar.Update(ctx, a); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "avatar.update", Target: "avatar:" + a.UUID})
	})
}

func (s *AvatarService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Avatar.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "avatar.delete", Target: fmt.Sprintf("avatar:%d", id)})
	})
}
