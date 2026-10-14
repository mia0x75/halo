package service

import (
	"context"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/domain"
)

type fakeStatementRepo struct {
	byKey map[[2]uint]*domain.Statement
}

func newFakeStatementRepo() *fakeStatementRepo {
	return &fakeStatementRepo{byKey: make(map[[2]uint]*domain.Statement)}
}

func stmtKey(ticketID uint, sequence uint16) [2]uint {
	return [2]uint{ticketID, uint(sequence)}
}

func (f *fakeStatementRepo) Create(_ context.Context, s *domain.Statement) error {
	k := stmtKey(s.TicketID, s.Sequence)
	if _, ok := f.byKey[k]; ok {
		return domain.ErrConflict
	}
	cp := *s
	f.byKey[k] = &cp
	return nil
}

func (f *fakeStatementRepo) GetByUUID(_ context.Context, uuid string) (*domain.Statement, error) {
	for _, s := range f.byKey {
		if s.UUID == uuid {
			cp := *s
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeStatementRepo) GetByTicketSequence(_ context.Context, ticketID uint, sequence uint16) (*domain.Statement, error) {
	s, ok := f.byKey[stmtKey(ticketID, sequence)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *s
	return &cp, nil
}

func (f *fakeStatementRepo) Update(_ context.Context, s *domain.Statement) error {
	k := stmtKey(s.TicketID, s.Sequence)
	if _, ok := f.byKey[k]; !ok {
		return domain.ErrNotFound
	}
	cp := *s
	f.byKey[k] = &cp
	return nil
}

func (f *fakeStatementRepo) DeleteByUUID(_ context.Context, uuid string) error {
	for k, s := range f.byKey {
		if s.UUID == uuid {
			delete(f.byKey, k)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeStatementRepo) DeleteByTicketSequence(_ context.Context, ticketID uint, sequence uint16) error {
	delete(f.byKey, stmtKey(ticketID, sequence))
	return nil
}

func (f *fakeStatementRepo) DeleteByTicket(_ context.Context, ticketID uint) error {
	for k := range f.byKey {
		if k[0] == ticketID {
			delete(f.byKey, k)
		}
	}
	return nil
}

func (f *fakeStatementRepo) List(_ context.Context, _ string, _ int) ([]*domain.Statement, error) {
	return nil, nil
}

type fakeClusterRepo struct {
	byID   map[uint]*domain.Cluster
	nextID uint
}

func newFakeClusterRepo() *fakeClusterRepo {
	return &fakeClusterRepo{byID: make(map[uint]*domain.Cluster)}
}

func (f *fakeClusterRepo) Create(_ context.Context, c *domain.Cluster) (int64, error) {
	f.nextID++
	if c.ClusterID == 0 {
		c.ClusterID = f.nextID
	}
	cp := *c
	f.byID[c.ClusterID] = &cp
	return int64(c.ClusterID), nil
}

func (f *fakeClusterRepo) GetByID(_ context.Context, id int64) (*domain.Cluster, error) {
	c, ok := f.byID[uint(id)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (f *fakeClusterRepo) GetByUUID(_ context.Context, uuid string) (*domain.Cluster, error) {
	for _, c := range f.byID {
		if c.UUID == uuid {
			cp := *c
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeClusterRepo) Update(_ context.Context, c *domain.Cluster) error {
	if _, ok := f.byID[c.ClusterID]; !ok {
		return domain.ErrNotFound
	}
	cp := *c
	f.byID[c.ClusterID] = &cp
	return nil
}

func (f *fakeClusterRepo) Delete(_ context.Context, id int64) error {
	delete(f.byID, uint(id))
	return nil
}

func (f *fakeClusterRepo) DeleteByUUID(_ context.Context, uuid string) error {
	for id, c := range f.byID {
		if c.UUID == uuid {
			delete(f.byID, id)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeClusterRepo) List(_ context.Context, _ int64, _ int) ([]*domain.Cluster, error) {
	return nil, nil
}

type fakeExecutor struct {
	results map[string]*agent.ExecResult
	err     error
	calls   int
}

func (f *fakeExecutor) Execute(_ context.Context, _ *domain.Cluster, _, query string) (*agent.ExecResult, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if r, ok := f.results[query]; ok {
		return r, nil
	}
	return &agent.ExecResult{RowsAffected: 1}, nil
}

func (f *fakeExecutor) LoadSchema(_ context.Context, _ *domain.Cluster, _ string) (*agent.Schema, error) {
	return nil, nil
}

type fakeValidator struct {
	violations []domain.Violation
	err        error
}

func (f *fakeValidator) Validate(_ context.Context, _ domain.ValidateRequest) ([]domain.Violation, error) {
	return f.violations, f.err
}
