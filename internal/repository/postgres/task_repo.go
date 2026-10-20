package postgresrepo

import (
	"context"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.TaskRepository = (*TaskRepo)(nil)

type TaskRepo struct {
	q *Queries
}

func NewTaskRepo(q *Queries) *TaskRepo {
	return &TaskRepo{q: q}
}

func taskFromRow(r MmTask) *domain.Task {
	return &domain.Task{
		TaskID:      uint(r.TaskID),
		UUID:        r.Uuid,
		Name:        r.Name,
		Type:        r.Type,
		Payload:     r.Payload,
		Status:      uint8(r.Status),
		Priority:    uint8(r.Priority),
		Attempt:     uint8(r.Attempt),
		MaxAttempt:  uint8(r.MaxAttempt),
		CronID:      uint(r.CronID),
		ScheduledAt: uint(r.ScheduledAt.Unix()),
		StartedAt:   timeFromNull(r.StartedAt),
		FinishedAt:  timeFromNull(r.FinishedAt),
		Error:       r.Error,
		Version:     int(r.Version),
		UpdateAt:    timeFromNull(r.UpdateAt),
		CreateAt:    timeToUnix(r.CreateAt),
	}
}

func (r *TaskRepo) Create(ctx context.Context, t *domain.Task) error {
	return wrap("task.Create", r.q.CreateTask(ctx, CreateTaskParams{
		Uuid:        t.UUID,
		Name:        t.Name,
		Type:        t.Type,
		Payload:     t.Payload,
		Status:      int16(t.Status),
		Priority:    int16(t.Priority),
		Attempt:     int16(t.Attempt),
		MaxAttempt:  int16(t.MaxAttempt),
		CronID:      int64(t.CronID),
		ScheduledAt: time.Unix(int64(t.ScheduledAt), 0).UTC(),
		Error:       t.Error,
		Version:     int64(t.Version),
	}))
}

func (r *TaskRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Task, error) {
	row, err := r.q.GetTaskByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("task.GetByUUID", err)
	}
	return taskFromRow(row), nil
}

func (r *TaskRepo) Update(ctx context.Context, t *domain.Task) error {
	return wrap("task.Update", r.q.UpdateTask(ctx, UpdateTaskParams{
		Name:        t.Name,
		Type:        t.Type,
		Payload:     t.Payload,
		Status:      int16(t.Status),
		Priority:    int16(t.Priority),
		Attempt:     int16(t.Attempt),
		MaxAttempt:  int16(t.MaxAttempt),
		CronID:      int64(t.CronID),
		ScheduledAt: time.Unix(int64(t.ScheduledAt), 0).UTC(),
		StartedAt:   nullFromTime(t.StartedAt),
		FinishedAt:  nullFromTime(t.FinishedAt),
		Error:       t.Error,
		Version:     int64(t.Version),
		Uuid:        t.UUID,
	}))
}

func (r *TaskRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("task.DeleteByUUID", r.q.DeleteTaskByUUID(ctx, uuid))
}

func (r *TaskRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Task, error) {
	rows, err := r.q.ListTasks(ctx, ListTasksParams{
		Uuid:  cursor,
		Limit: clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("task.List", err)
	}
	out := make([]*domain.Task, 0, len(rows))
	for _, row := range rows {
		out = append(out, taskFromRow(row))
	}
	return out, nil
}

func (r *TaskRepo) ListPending(ctx context.Context, limit int) ([]*domain.Task, error) {
	rows, err := r.q.ListPendingTasks(ctx, clampInt32(limit))
	if err != nil {
		return nil, wrap("task.ListPending", err)
	}
	out := make([]*domain.Task, 0, len(rows))
	for _, row := range rows {
		out = append(out, taskFromRow(row))
	}
	return out, nil
}

func (r *TaskRepo) CountByStatus(ctx context.Context) ([]domain.TaskStatusCount, error) {
	rows, err := r.q.CountTasksByStatus(ctx)
	if err != nil {
		return nil, wrap("task.CountByStatus", err)
	}
	out := make([]domain.TaskStatusCount, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.TaskStatusCount{Status: uint8(row.Status), Count: uint(row.Count)})
	}
	return out, nil
}
