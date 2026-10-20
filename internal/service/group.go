package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// GroupService manages the cluster groups: the servers of one group belong
// together, and roles are granted access to groups.
type GroupService struct {
	groups domain.GroupRepository
}

func NewGroupService(groups domain.GroupRepository) *GroupService {
	return &GroupService{groups: groups}
}

func (s *GroupService) GetByUUID(ctx context.Context, uuid string) (*domain.Group, error) {
	return s.groups.GetByUUID(ctx, uuid)
}

func (s *GroupService) GetByName(ctx context.Context, name string) (*domain.Group, error) {
	return s.groups.GetByName(ctx, name)
}

func (s *GroupService) List(ctx context.Context, cursor string, limit int) ([]*domain.Group, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.groups.List(ctx, cursor, limit)
}

func (s *GroupService) Create(ctx context.Context, g *domain.Group) error {
	return s.groups.Create(ctx, g)
}

func (s *GroupService) Update(ctx context.Context, g *domain.Group) error {
	return s.groups.Update(ctx, g)
}

func (s *GroupService) DeleteByUUID(ctx context.Context, uuid string) error {
	return s.groups.DeleteByUUID(ctx, uuid)
}
