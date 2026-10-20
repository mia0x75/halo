package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// TaskService keeps the task executions: every run of an asynchronous or
// scheduled job is recorded, so a restart cannot lose the history and a
// failure can be traced.
type TaskService struct {
	tasks domain.TaskRepository
}

func NewTaskService(tasks domain.TaskRepository) *TaskService {
	return &TaskService{tasks: tasks}
}

// Enqueue records a new task execution waiting for a worker.
func (s *TaskService) Enqueue(ctx context.Context, t *domain.Task) error {
	return s.tasks.Create(ctx, t)
}

func (s *TaskService) GetByUUID(ctx context.Context, uuid string) (*domain.Task, error) {
	return s.tasks.GetByUUID(ctx, uuid)
}

// Update writes the new state of a task execution.
func (s *TaskService) Update(ctx context.Context, t *domain.Task) error {
	return s.tasks.Update(ctx, t)
}

func (s *TaskService) DeleteByUUID(ctx context.Context, uuid string) error {
	return s.tasks.DeleteByUUID(ctx, uuid)
}

// List returns the task executions, newest first.
func (s *TaskService) List(ctx context.Context, cursor string, limit int) ([]*domain.Task, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.tasks.List(ctx, cursor, limit)
}

// ListPending returns the tasks waiting for a worker.
func (s *TaskService) ListPending(ctx context.Context, limit int) ([]*domain.Task, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.tasks.ListPending(ctx, limit)
}

// CountByStatus counts the task executions per status.
func (s *TaskService) CountByStatus(ctx context.Context) ([]domain.TaskStatusCount, error) {
	return s.tasks.CountByStatus(ctx)
}
