package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.RoleRepository = (*RoleRepo)(nil)

type RoleRepo struct {
	q *Queries
}

func NewRoleRepo(q *Queries) *RoleRepo {
	return &RoleRepo{q: q}
}

func roleFromRow(r MmRole) *domain.Role {
	return &domain.Role{
		RoleID:      uint(r.RoleID),
		UUID:        r.Uuid,
		Name:        r.Name,
		Description: r.Description,
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *RoleRepo) Create(ctx context.Context, role *domain.Role) (int64, error) {
	id, err := r.q.CreateRole(ctx, CreateRoleParams{
		Uuid:        role.UUID,
		Name:        role.Name,
		Description: role.Description,
		Version:     int64(role.Version),
	})
	if err != nil {
		return 0, wrap("role.Create", err)
	}
	return id, nil
}

func (r *RoleRepo) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	row, err := r.q.GetRoleByID(ctx, id)
	if err != nil {
		return nil, wrap("role.GetByID", err)
	}
	return roleFromRow(row), nil
}

func (r *RoleRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Role, error) {
	row, err := r.q.GetRoleByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("role.GetByUUID", err)
	}
	return roleFromRow(row), nil
}

func (r *RoleRepo) Update(ctx context.Context, role *domain.Role) error {
	err := r.q.UpdateRole(ctx, UpdateRoleParams{
		Uuid:        role.UUID,
		Name:        role.Name,
		Description: role.Description,
		Version:     int64(role.Version),
		RoleID:      int64(role.RoleID),
	})
	return wrap("role.Update", err)
}

func (r *RoleRepo) Delete(ctx context.Context, id int64) error {
	return wrap("role.Delete", r.q.DeleteRoleByID(ctx, id))
}

func (r *RoleRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("role.DeleteByUUID", r.q.DeleteRoleByUUID(ctx, uuid))
}

func (r *RoleRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Role, error) {
	rows, err := r.q.ListRoles(ctx, ListRolesParams{
		RoleID: cursor,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, wrap("role.List", err)
	}
	out := make([]*domain.Role, 0, len(rows))
	for _, row := range rows {
		out = append(out, roleFromRow(row))
	}
	return out, nil
}
