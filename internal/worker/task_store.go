package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mia0x75/halo/internal/domain"
)

// Task status codes used by the worker, mapped onto the mm_tasks status.
const (
	statusPending   = "P"
	statusRunning   = "R"
	statusSucceeded = "S"
	statusFailed    = "F"
	statusCancelled = "C"
	statusHalted    = "H"
)

// taskStatusFromCode maps a worker status code to the stored task status.
func taskStatusFromCode(code string) uint8 {
	switch code {
	case statusRunning:
		return domain.TaskStatusRunning
	case statusSucceeded:
		return domain.TaskStatusSucceeded
	case statusFailed:
		return domain.TaskStatusFailed
	case statusCancelled:
		return domain.TaskStatusCancelled
	default:
		return domain.TaskStatusPending
	}
}

// newTaskUUID returns the uuid of a task execution row.
func newTaskUUID() string {
	return uuid.NewString()
}

// syncTask records the execution of a scheduled task in mm_tasks. The row is
// created with the schedule and updated when the task runs, so the execution
// history survives a restart.
func (s *Scheduler) syncTask(task *Task, cronID uint) error {
	ctx := context.Background()

	payload, err := json.Marshal(task.Params)
	if err != nil {
		return fmt.Errorf("marshal task params: %w", err)
	}

	existing, getErr := repos.Task.GetLatestByCron(ctx, cronID)
	if getErr != nil && !errors.Is(getErr, domain.ErrNotFound) {
		return getErr
	}

	now := uint(time.Now().Unix())
	scheduledAt := now
	if !task.NextRun.IsZero() {
		scheduledAt = uint(task.NextRun.Unix())
	}
	if existing == nil {
		return repos.Task.Create(ctx, &domain.Task{
			UUID:        newTaskUUID(),
			Name:        task.Name,
			Type:        task.Cmd,
			Payload:     string(payload),
			Status:      domain.TaskStatusPending,
			MaxAttempt:  3,
			CronID:      cronID,
			ScheduledAt: scheduledAt,
			Version:     1,
		})
	}

	existing.Status = taskStatusFromCode(task.Status)
	existing.Error = ""
	if task.Err != nil {
		existing.Error = task.Err.Error()
	}
	if existing.Status == domain.TaskStatusRunning {
		existing.StartedAt = &now
	}
	if existing.Status == domain.TaskStatusSucceeded || existing.Status == domain.TaskStatusFailed {
		existing.FinishedAt = &now
	}
	existing.Version++

	return repos.Task.Update(ctx, existing)
}
