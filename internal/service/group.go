package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mia0x75/halo/internal/domain"
)

// GroupService manages the cluster groups: the servers of one group belong
// together, and roles are granted access to groups.
type GroupService struct {
	tx       TxManager
	groups   domain.GroupRepository
	clusters domain.ClusterRepository
	log      *LogService
}

func NewGroupService(tx TxManager, groups domain.GroupRepository, clusters domain.ClusterRepository, log *LogService) *GroupService {
	return &GroupService{tx: tx, groups: groups, clusters: clusters, log: log}
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

// Create stores a new group. The name identifies the group, so it must be
// unique.
func (s *GroupService) Create(ctx context.Context, g *domain.Group, operatorID uint) error {
	name := strings.TrimSpace(g.Name)
	if name == "" {
		return fmt.Errorf("%w: 分组名称不能为空", domain.ErrInvalidInput)
	}
	if _, err := s.groups.GetByName(ctx, name); err == nil {
		return fmt.Errorf("%w: 分组名称(name=%s)已存在", domain.ErrConflict, name)
	} else if !isNotFound(err) {
		return err
	}
	g.Name = name
	if g.Status == 0 {
		g.Status = domain.GroupStatusNormal
	}
	g.Version = 1

	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Group.Create(ctx, g); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    operatorID,
			Operation: "group.create",
			Target:    "group:" + g.UUID,
			Detail:    g.Name,
		})
	})
}

// Update changes the name, description or status of a group.
func (s *GroupService) Update(ctx context.Context, g *domain.Group, operatorID uint) error {
	current, err := s.groups.GetByUUID(ctx, g.UUID)
	if err != nil {
		return err
	}
	name := strings.TrimSpace(g.Name)
	if name == "" {
		return fmt.Errorf("%w: 分组名称不能为空", domain.ErrInvalidInput)
	}
	if name != current.Name {
		existing, err := s.groups.GetByName(ctx, name)
		switch {
		case err == nil && existing.UUID != g.UUID:
			return fmt.Errorf("%w: 分组名称(name=%s)已存在", domain.ErrConflict, name)
		case err != nil && !isNotFound(err):
			return err
		}
	}
	g.Name = name
	g.Version = current.Version + 1

	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Group.Update(ctx, g); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    operatorID,
			Operation: "group.update",
			Target:    "group:" + g.UUID,
			Detail:    g.Name,
		})
	})
}

// DeleteByUUID removes a group. Clusters must be moved to another group first,
// otherwise they would keep a dangling group_id.
func (s *GroupService) DeleteByUUID(ctx context.Context, uuid string, operatorID uint) error {
	group, err := s.groups.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	clusters, err := s.clusters.CountByGroup(ctx, group.GroupID)
	if err != nil {
		return fmt.Errorf("count group clusters: %w", err)
	}
	if clusters > 0 {
		return fmt.Errorf("%w: 分组下还有 %d 个群集，请先迁移", domain.ErrConflict, clusters)
	}

	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Group.DeleteByUUID(ctx, uuid); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    operatorID,
			Operation: "group.delete",
			Target:    "group:" + uuid,
		})
	})
}

// isNotFound reports whether the error means "no such row".
func isNotFound(err error) bool {
	return err != nil && errors.Is(err, domain.ErrNotFound)
}
