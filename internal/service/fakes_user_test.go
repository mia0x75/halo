package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

type fakeUserRepo struct {
	byID   map[uint]*domain.User
	nextID uint
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byID: make(map[uint]*domain.User)}
}

func (f *fakeUserRepo) Create(_ context.Context, u *domain.User) (int64, error) {
	f.nextID++
	if u.UserID == 0 {
		u.UserID = f.nextID
	}
	cp := *u
	f.byID[u.UserID] = &cp
	return int64(u.UserID), nil
}

func (f *fakeUserRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	u, ok := f.byID[uint(id)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (f *fakeUserRepo) GetByUUID(_ context.Context, uuid string) (*domain.User, error) {
	for _, u := range f.byID {
		if u.UUID == uuid {
			cp := *u
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepo) Update(_ context.Context, u *domain.User) error {
	if _, ok := f.byID[u.UserID]; !ok {
		return domain.ErrNotFound
	}
	cp := *u
	f.byID[u.UserID] = &cp
	return nil
}

func (f *fakeUserRepo) Delete(_ context.Context, id int64) error {
	delete(f.byID, uint(id))
	return nil
}

func (f *fakeUserRepo) DeleteByUUID(_ context.Context, uuid string) error {
	for id, u := range f.byID {
		if u.UUID == uuid {
			delete(f.byID, id)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeUserRepo) List(_ context.Context, _ int64, _ int) ([]*domain.User, error) {
	return nil, nil
}

func (f *fakeUserRepo) Count(_ context.Context) (int64, error) {
	return int64(len(f.byID)), nil
}
