package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.UserRepository = (*UserRepo)(nil)

type UserRepo struct {
	q *Queries
}

func NewUserRepo(q *Queries) *UserRepo {
	return &UserRepo{q: q}
}

func userFromRow(r MmUser) *domain.User {
	return &domain.User{
		UserID:   uint(r.UserID),
		UUID:     r.Uuid,
		Email:    r.Email,
		Password: r.Password,
		Status:   domain.UserStatus(r.Status),
		Name:     r.Name,
		Phone:    r.Phone,
		AvatarID: uint(r.AvatarID),
		Version:  int(r.Version),
		UpdateAt: timeFromNull(r.UpdateAt),
		CreateAt: timeToUnix(r.CreateAt),
	}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) (int64, error) {
	id, err := r.q.CreateUser(ctx, CreateUserParams{
		Uuid:     u.UUID,
		Email:    u.Email,
		Password: u.Password,
		Status:   uint8(u.Status),
		Name:     u.Name,
		Phone:    u.Phone,
		AvatarID: uint32(u.AvatarID),
		Version:  uint32(u.Version),
	})
	if err != nil {
		return 0, wrap("user.Create", err)
	}
	return id, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	row, err := r.q.GetUserByID(ctx, clampUint32(id))
	if err != nil {
		return nil, wrap("user.GetByID", err)
	}
	return userFromRow(row), nil
}

func (r *UserRepo) GetByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	row, err := r.q.GetUserByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("user.GetByUUID", err)
	}
	return userFromRow(row), nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	err := r.q.UpdateUser(ctx, UpdateUserParams{
		Uuid:     u.UUID,
		Email:    u.Email,
		Password: u.Password,
		Status:   uint8(u.Status),
		Name:     u.Name,
		Phone:    u.Phone,
		AvatarID: uint32(u.AvatarID),
		Version:  uint32(u.Version),
		UserID:   uint32(u.UserID),
	})
	return wrap("user.Update", err)
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	return wrap("user.Delete", r.q.DeleteUserByID(ctx, clampUint32(id)))
}

func (r *UserRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("user.DeleteByUUID", r.q.DeleteUserByUUID(ctx, uuid))
}

func (r *UserRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.User, error) {
	rows, err := r.q.ListUsers(ctx, ListUsersParams{
		UserID: clampUint32(cursor),
		Limit:  clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("user.List", err)
	}
	out := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		out = append(out, userFromRow(row))
	}
	return out, nil
}

// Count returns the total number of users.
func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	n, err := r.q.CountUsers(ctx)
	if err != nil {
		return 0, wrap("user.Count", err)
	}
	return n, nil
}
