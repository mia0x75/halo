package resolvers

import (
	"context"

	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// TaskExecutions 任务执行记录
func (r *queryRootResolver) TaskExecutions(ctx context.Context, status *uint8, first *int) (L []*gqlapi.TaskExecution, err error) {
	limit := 50
	if first != nil && *first > 0 {
		limit = *first
	}

	tasks, err := r.app.Services.Task.List(ctx, "", limit)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if status != nil && task.Status != *status {
			continue
		}
		L = append(L, &gqlapi.TaskExecution{
			UUID:        task.UUID,
			Name:        task.Name,
			Type:        task.Type,
			Payload:     task.Payload,
			Status:      task.Status,
			Attempt:     task.Attempt,
			CronID:      task.CronID,
			ScheduledAt: task.ScheduledAt,
			StartedAt:   task.StartedAt,
			FinishedAt:  task.FinishedAt,
			Error:       task.Error,
			CreateAt:    task.CreateAt,
		})
	}
	return L, nil
}

// TaskExecutionSummary 按状态统计任务执行记录
func (r *queryRootResolver) TaskExecutionSummary(ctx context.Context) (L []*gqlapi.TaskExecutionCount, err error) {
	counts, err := r.app.Services.Task.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}
	for _, count := range counts {
		L = append(L, &gqlapi.TaskExecutionCount{Status: count.Status, Count: uint(count.Count)})
	}
	return L, nil
}
