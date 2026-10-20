package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// fakeTx is a non-transactional TxManager: it simply hands the same Repos to fn.
// It is enough to exercise service orchestration; atomicity is the repository
// TxManager's concern and is covered there.
type fakeTx struct {
	repos Repos
}

func (f *fakeTx) WithTx(ctx context.Context, fn func(ctx context.Context, repos Repos) error) error {
	return fn(ctx, f.repos)
}

type fakeTicketRepo struct {
	byID   map[uint]*domain.Ticket
	nextID uint
	// reviewer 记录按审核人查询时返回的工单
	reviewer []*domain.Ticket
}

func newFakeTicketRepo() *fakeTicketRepo {
	return &fakeTicketRepo{byID: make(map[uint]*domain.Ticket)}
}

func (f *fakeTicketRepo) Create(_ context.Context, t *domain.Ticket) (int64, error) {
	f.nextID++
	if t.TicketID == 0 {
		t.TicketID = f.nextID
	}
	if _, ok := f.byID[t.TicketID]; ok {
		return 0, domain.ErrConflict
	}
	cp := *t
	f.byID[t.TicketID] = &cp
	return int64(t.TicketID), nil
}

func (f *fakeTicketRepo) GetByID(_ context.Context, id int64) (*domain.Ticket, error) {
	t, ok := f.byID[uint(id)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (f *fakeTicketRepo) GetByUUID(_ context.Context, uuid string) (*domain.Ticket, error) {
	for _, t := range f.byID {
		if t.UUID == uuid {
			cp := *t
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeTicketRepo) Update(_ context.Context, t *domain.Ticket) error {
	if _, ok := f.byID[t.TicketID]; !ok {
		return domain.ErrNotFound
	}
	cp := *t
	f.byID[t.TicketID] = &cp
	return nil
}

func (f *fakeTicketRepo) Delete(_ context.Context, id int64) error {
	delete(f.byID, uint(id))
	return nil
}

func (f *fakeTicketRepo) DeleteByUUID(_ context.Context, uuid string) error {
	for id, t := range f.byID {
		if t.UUID == uuid {
			delete(f.byID, id)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeTicketRepo) List(_ context.Context, _ int64, _ int) ([]*domain.Ticket, error) {
	return nil, nil
}

func (f *fakeTicketRepo) ListFiltered(_ context.Context, _ int64, _ int, _ uint8, _ string) ([]*domain.Ticket, error) {
	return f.reviewer, nil
}

func (f *fakeTicketRepo) ListByUser(_ context.Context, _ uint) ([]*domain.Ticket, error) {
	return nil, nil
}

func (f *fakeTicketRepo) ListByReviewer(_ context.Context, _ uint) ([]*domain.Ticket, error) {
	return f.reviewer, nil
}

func (f *fakeTicketRepo) GetByCronID(_ context.Context, _ uint) (*domain.Ticket, error) {
	return nil, domain.ErrNotFound
}

type fakeLogRepo struct {
	entries []*domain.Log
	err     error
}

func (f *fakeLogRepo) Create(_ context.Context, l *domain.Log) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.entries = append(f.entries, l)
	return int64(len(f.entries)), nil
}

func (f *fakeLogRepo) GetByID(_ context.Context, _ int64) (*domain.Log, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeLogRepo) GetByUUID(_ context.Context, _ string) (*domain.Log, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeLogRepo) Update(_ context.Context, _ *domain.Log) error { return nil }
func (f *fakeLogRepo) Delete(_ context.Context, _ int64) error       { return nil }
func (f *fakeLogRepo) DeleteByUUID(_ context.Context, _ string) error {
	return domain.ErrNotFound
}
func (f *fakeLogRepo) List(_ context.Context, _ int64, _ int) ([]*domain.Log, error) {
	return nil, nil
}

// CountByStatus counts the stored tickets per status inside a time window.
func (f *fakeTicketRepo) CountByStatus(context.Context, uint, uint) ([]domain.TicketStatusCount, error) {
	counts := map[domain.TicketStatus]uint{}
	for _, ticket := range f.byID {
		counts[ticket.Status]++
	}
	out := make([]domain.TicketStatusCount, 0, len(counts))
	for status, count := range counts {
		out = append(out, domain.TicketStatusCount{Status: status, Count: count})
	}
	return out, nil
}
