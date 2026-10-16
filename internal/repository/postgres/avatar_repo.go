package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.AvatarRepository = (*AvatarRepo)(nil)

type AvatarRepo struct {
	q *Queries
}

func NewAvatarRepo(q *Queries) *AvatarRepo {
	return &AvatarRepo{q: q}
}

func avatarFromRow(r MmAvatar) *domain.Avatar {
	return &domain.Avatar{
		AvatarID: uint(r.AvatarID),
		UUID:     r.Uuid,
		URL:      r.Url,
		Version:  int(r.Version),
		UpdateAt: timeFromNull(r.UpdateAt),
		CreateAt: timeToUnix(r.CreateAt),
	}
}

func (r *AvatarRepo) Create(ctx context.Context, a *domain.Avatar) (int64, error) {
	id, err := r.q.CreateAvatar(ctx, CreateAvatarParams{
		Uuid:    a.UUID,
		Url:     a.URL,
		Version: int64(a.Version),
	})
	if err != nil {
		return 0, wrap("avatar.Create", err)
	}
	return id, nil
}

func (r *AvatarRepo) GetByID(ctx context.Context, id int64) (*domain.Avatar, error) {
	row, err := r.q.GetAvatarByID(ctx, id)
	if err != nil {
		return nil, wrap("avatar.GetByID", err)
	}
	return avatarFromRow(row), nil
}

func (r *AvatarRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Avatar, error) {
	row, err := r.q.GetAvatarByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("avatar.GetByUUID", err)
	}
	return avatarFromRow(row), nil
}

func (r *AvatarRepo) Update(ctx context.Context, a *domain.Avatar) error {
	err := r.q.UpdateAvatar(ctx, UpdateAvatarParams{
		Uuid:     a.UUID,
		Url:      a.URL,
		Version:  int64(a.Version),
		AvatarID: int64(a.AvatarID),
	})
	return wrap("avatar.Update", err)
}

func (r *AvatarRepo) Delete(ctx context.Context, id int64) error {
	return wrap("avatar.Delete", r.q.DeleteAvatarByID(ctx, id))
}

func (r *AvatarRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("avatar.DeleteByUUID", r.q.DeleteAvatarByUUID(ctx, uuid))
}

func (r *AvatarRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Avatar, error) {
	rows, err := r.q.ListAvatars(ctx, ListAvatarsParams{
		AvatarID: cursor,
		Limit:    clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("avatar.List", err)
	}
	out := make([]*domain.Avatar, 0, len(rows))
	for _, row := range rows {
		out = append(out, avatarFromRow(row))
	}
	return out, nil
}
