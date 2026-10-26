package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/service"
)

// fakeCronStore keeps crons in memory and mirrors the schedulable filter of
// ListSchedulable so a cancelled schedule cannot sneak back through reload.
type fakeCronStore struct {
	byUUID map[string]*domain.Cron
	nextID uint
}

func (f *fakeCronStore) Create(_ context.Context, c *domain.Cron) (int64, error) {
	f.nextID++
	c.CronID = f.nextID
	cp := *c
	f.byUUID[c.UUID] = &cp
	return int64(c.CronID), nil
}

func (f *fakeCronStore) GetByID(_ context.Context, id int64) (*domain.Cron, error) {
	for _, c := range f.byUUID {
		if c.CronID == uint(id) {
			cp := *c
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeCronStore) GetByUUID(_ context.Context, uuid string) (*domain.Cron, error) {
	c, ok := f.byUUID[uuid]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (f *fakeCronStore) Update(_ context.Context, c *domain.Cron) error {
	if _, ok := f.byUUID[c.UUID]; !ok {
		return domain.ErrNotFound
	}
	cp := *c
	f.byUUID[c.UUID] = &cp
	return nil
}

func (f *fakeCronStore) Delete(_ context.Context, _ int64) error { return nil }

func (f *fakeCronStore) DeleteByUUID(_ context.Context, uuid string) error {
	delete(f.byUUID, uuid)
	return nil
}

func (f *fakeCronStore) List(_ context.Context, _ int64, _ int) ([]*domain.Cron, error) {
	out := make([]*domain.Cron, 0, len(f.byUUID))
	for _, c := range f.byUUID {
		cp := *c
		out = append(out, &cp)
	}
	return out, nil
}

func (f *fakeCronStore) ListSchedulable(_ context.Context) ([]*domain.Cron, error) {
	out := []*domain.Cron{}
	for _, c := range f.byUUID {
		recurrent := c.Recurrent != nil && *c.Recurrent == 1
		pending := c.Status == domain.CronStatusPending
		if recurrent || pending {
			cp := *c
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (f *fakeCronStore) Count(_ context.Context) (int64, error) {
	return int64(len(f.byUUID)), nil
}

// fakeTaskStore only records executions; reload never reads from it.
type fakeTaskStore struct {
	created int
}

func (f *fakeTaskStore) Create(_ context.Context, _ *domain.Task) error {
	f.created++
	return nil
}

func (f *fakeTaskStore) GetByUUID(_ context.Context, _ string) (*domain.Task, error) {
	return nil, domain.ErrNotFound
}

func (f *fakeTaskStore) Update(_ context.Context, _ *domain.Task) error { return nil }

func (f *fakeTaskStore) DeleteByUUID(_ context.Context, _ string) error { return nil }

func (f *fakeTaskStore) List(_ context.Context, _ string, _ int) ([]*domain.Task, error) {
	return nil, nil
}

func (f *fakeTaskStore) ListPending(_ context.Context, _ int) ([]*domain.Task, error) {
	return nil, nil
}

func (f *fakeTaskStore) GetLatestByCron(_ context.Context, _ uint) (*domain.Task, error) {
	return nil, domain.ErrNotFound
}

func (f *fakeTaskStore) CountByStatus(_ context.Context) ([]domain.TaskStatusCount, error) {
	return nil, nil
}

func (f *fakeTaskStore) DeleteBefore(_ context.Context, _ uint) (int64, error) {
	return 0, nil
}

func TestSchedulerCancelPersistsCancelledStatus(t *testing.T) {
	crons := &fakeCronStore{byUUID: make(map[string]*domain.Cron)}
	tasks := &fakeTaskStore{}
	Init(service.Repos{Cron: crons, Task: tasks})

	interval := "0s"
	recurrent := uint8(0)
	_, err := crons.Create(context.Background(), &domain.Cron{
		CronID:    1,
		UUID:      "cron-1",
		Status:    domain.CronStatusPending,
		Name:      "execute ticket",
		Cmd:       "halocli",
		Params:    `["execute","-T","t1"]`,
		Interval:  &interval,
		Recurrent: &recurrent,
	})
	require.NoError(t, err)

	// A scheduler without Start: register/Cancel only touch the task map.
	s := &Scheduler{tasks: map[string]*Task{}}
	task := NewTask("execute ticket", "halocli", []string{"execute", "-T", "t1"})
	task.UUID = "cron-1"
	task.NextRun = time.Now().Add(time.Hour)
	s.tasks[task.UUID] = task

	require.NoError(t, s.Cancel(task.UUID))

	stored, err := crons.GetByUUID(context.Background(), "cron-1")
	require.NoError(t, err)
	// The row must leave the schedulable set, otherwise a later reload
	// re-registers the cancelled schedule.
	require.Equal(t, domain.CronStatusCancelled, stored.Status)
	require.NotContains(t, s.tasks, "cron-1")

	require.NoError(t, s.reload())
	require.NotContains(t, s.tasks, "cron-1")

	require.Error(t, s.Cancel("missing"))
}
