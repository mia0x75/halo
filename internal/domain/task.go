package domain

import "context"

// TaskStatus is the state of a task execution.
const (
	// TaskStatusPending waits for a worker.
	TaskStatusPending uint8 = iota
	// TaskStatusRunning is being executed.
	TaskStatusRunning
	// TaskStatusSucceeded finished without error.
	TaskStatusSucceeded
	// TaskStatusFailed failed and may be retried until max attempt.
	TaskStatusFailed
	// TaskStatusCancelled was cancelled before or during execution.
	TaskStatusCancelled
)

// Task is a row of mm_tasks: one execution of an asynchronous or scheduled
// job. The schedule itself lives in mm_crons.
type Task struct {
	TaskID      uint
	UUID        string
	Name        string
	Type        string
	Payload     string
	Status      uint8
	Priority    uint8
	Attempt     uint8
	MaxAttempt  uint8
	CronID      uint
	ScheduledAt uint
	StartedAt   *uint
	FinishedAt  *uint
	Error       string
	Version     int
	UpdateAt    *uint
	CreateAt    uint
}

// TaskRepository stores the task executions.
type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	GetByUUID(ctx context.Context, uuid string) (*Task, error)
	Update(ctx context.Context, t *Task) error
	DeleteByUUID(ctx context.Context, uuid string) error
	// List returns the task executions, newest first.
	List(ctx context.Context, cursor string, limit int) ([]*Task, error)
	// ListPending returns the tasks waiting for a worker.
	ListPending(ctx context.Context, limit int) ([]*Task, error)
	// CountByStatus counts the task executions per status.
	CountByStatus(ctx context.Context) ([]TaskStatusCount, error)
}

// TaskStatusCount counts the task executions of one status.
type TaskStatusCount struct {
	Status uint8
	Count  uint
}
