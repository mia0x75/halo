package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.GroupRepository = (*GroupRepo)(nil)

type GroupRepo struct {
	q *Queries
}

func NewGroupRepo(q *Queries) *GroupRepo {
	return &GroupRepo{q: q}
}

func groupFromRow(r MmGroup) *domain.Group {
	return &domain.Group{
		GroupID:     uint(r.GroupID),
		UUID:        r.Uuid,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *GroupRepo) Create(ctx context.Context, g *domain.Group) error {
	return wrap("group.Create", r.q.CreateGroup(ctx, CreateGroupParams{
		Uuid:        g.UUID,
		Name:        g.Name,
		Description: g.Description,
		Status:      g.Status,
		Version:     uint32(g.Version),
	}))
}

func (r *GroupRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Group, error) {
	row, err := r.q.GetGroupByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("group.GetByUUID", err)
	}
	return groupFromRow(row), nil
}

func (r *GroupRepo) GetByName(ctx context.Context, name string) (*domain.Group, error) {
	row, err := r.q.GetGroupByName(ctx, name)
	if err != nil {
		return nil, wrap("group.GetByName", err)
	}
	return groupFromRow(row), nil
}

func (r *GroupRepo) Update(ctx context.Context, g *domain.Group) error {
	return wrap("group.Update", r.q.UpdateGroup(ctx, UpdateGroupParams{
		Name:        g.Name,
		Description: g.Description,
		Status:      g.Status,
		Version:     uint32(g.Version),
		Uuid:        g.UUID,
	}))
}

func (r *GroupRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("group.DeleteByUUID", r.q.DeleteGroupByUUID(ctx, uuid))
}

func (r *GroupRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Group, error) {
	rows, err := r.q.ListGroups(ctx, ListGroupsParams{
		Uuid:  cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("group.List", err)
	}
	out := make([]*domain.Group, 0, len(rows))
	for _, row := range rows {
		out = append(out, groupFromRow(row))
	}
	return out, nil
}
