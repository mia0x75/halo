package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

// RetentionResult reports how many rows each retention run removed.
type RetentionResult struct {
	Metrics  int64
	Hits     int64
	Logs     int64
	Messages int64
	Tasks    int64
}

// RetentionService removes the data that grows without bound: the daily
// metric buckets, the rule hit details, the audit log, the in-app messages and
// the task execution records.
//
// The windows live in mm_options (retention.*-days) so an administrator can
// change them from the settings page; a missing option falls back to the
// built-in default.
type RetentionService struct {
	tx       TxManager
	options  domain.OptionRepository
	metrics  domain.MetricRepository
	hits     domain.RuleHitRepository
	logs     domain.LogRepository
	messages domain.MessageRepository
	tasks    domain.TaskRepository
}

// NewRetentionService wires the retention cleanup.
func NewRetentionService(
	tx TxManager,
	options domain.OptionRepository,
	metrics domain.MetricRepository,
	hits domain.RuleHitRepository,
	logs domain.LogRepository,
	messages domain.MessageRepository,
	tasks domain.TaskRepository,
) *RetentionService {
	return &RetentionService{
		tx:       tx,
		options:  options,
		metrics:  metrics,
		hits:     hits,
		logs:     logs,
		messages: messages,
		tasks:    tasks,
	}
}

// Retention window names and their defaults, in days.
const (
	retentionMetricsDays  = "retention.metrics-days"
	retentionHitsDays     = "retention.hits-days"
	retentionLogsDays     = "retention.logs-days"
	retentionMessagesDays = "retention.messages-days"
	retentionTasksDays    = "retention.tasks-days"
	retentionDefaultDays  = 90
	retentionLogsDefault  = 180
	retentionTasksDefault = 30
	retentionMinDays      = 1
	retentionMaxDays      = 3650
)

// Cleanup removes everything older than the configured windows.
func (s *RetentionService) Cleanup(ctx context.Context) (RetentionResult, error) {
	result := RetentionResult{}

	metricsDays, err := s.window(ctx, retentionMetricsDays, retentionDefaultDays)
	if err != nil {
		return result, err
	}
	hitsDays, err := s.window(ctx, retentionHitsDays, retentionDefaultDays)
	if err != nil {
		return result, err
	}
	logsDays, err := s.window(ctx, retentionLogsDays, retentionLogsDefault)
	if err != nil {
		return result, err
	}
	messagesDays, err := s.window(ctx, retentionMessagesDays, retentionDefaultDays)
	if err != nil {
		return result, err
	}
	tasksDays, err := s.window(ctx, retentionTasksDays, retentionTasksDefault)
	if err != nil {
		return result, err
	}

	now := time.Now()

	bucket := now.AddDate(0, 0, -metricsDays).Format("2006-01-02")
	metrics, err := s.metrics.DeleteBucketsBefore(ctx, bucket)
	if err != nil {
		return result, fmt.Errorf("clean metric buckets: %w", err)
	}
	result.Metrics = metrics

	hits, err := s.hits.DeleteBefore(ctx, uint(now.AddDate(0, 0, -hitsDays).Unix()))
	if err != nil {
		return result, fmt.Errorf("clean rule hits: %w", err)
	}
	result.Hits = hits

	logs, err := s.logs.DeleteBefore(ctx, uint(now.AddDate(0, 0, -logsDays).Unix()))
	if err != nil {
		return result, fmt.Errorf("clean logs: %w", err)
	}
	result.Logs = logs

	messages, err := s.messages.DeleteBefore(ctx, uint(now.AddDate(0, 0, -messagesDays).Unix()))
	if err != nil {
		return result, fmt.Errorf("clean messages: %w", err)
	}
	result.Messages = messages

	tasks, err := s.tasks.DeleteBefore(ctx, uint(now.AddDate(0, 0, -tasksDays).Unix()))
	if err != nil {
		return result, fmt.Errorf("clean tasks: %w", err)
	}
	result.Tasks = tasks

	return result, nil
}

// window reads one retention window from the options table.
func (s *RetentionService) window(ctx context.Context, name string, fallback int) (int, error) {
	option, err := s.options.GetByName(ctx, name)
	if err != nil {
		if isNotFound(err) {
			return fallback, nil
		}
		return 0, fmt.Errorf("read %s: %w", name, err)
	}
	if option.Parameter == "" {
		return fallback, nil
	}
	days, err := strconv.Atoi(option.Parameter)
	if err != nil {
		return 0, fmt.Errorf("%w: 选项 %s 的值必须是整数", domain.ErrInvalidInput, name)
	}
	if days < retentionMinDays || days > retentionMaxDays {
		return 0, fmt.Errorf("%w: 选项 %s 的天数必须在 %d 到 %d 之间", domain.ErrInvalidInput, name, retentionMinDays, retentionMaxDays)
	}
	return days, nil
}
