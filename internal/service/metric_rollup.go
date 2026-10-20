package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

// Metric names of the daily buckets.
//
// 工单类的每日计数（tickets-*）由事件流实时累加：工单在当天创建、在当天执行，
// 这些事实只在发生的瞬间可知，事后按「创建日 + 当前状态」统计会让历史数值
// 随状态推进而回溯变化，趋势不可信。汇总作业不再写这些桶。
const (
	// MetricTicketsCreated counts the tickets created on a day.
	MetricTicketsCreated = "tickets-created"
	// MetricTicketsExecuted counts the tickets executed on a day.
	MetricTicketsExecuted = "tickets-executed"
	// MetricTicketsFailed counts the tickets that failed to execute on a day.
	MetricTicketsFailed = "tickets-failed"
	// MetricViolationsError counts the error level rule hits of a day.
	MetricViolationsError = "violations-error"
	// MetricViolationsWarning counts the warning level rule hits of a day.
	MetricViolationsWarning = "violations-warning"
	// MetricViolationsNotice counts the notice level rule hits of a day.
	MetricViolationsNotice = "violations-notice"
	// MetricRollupDone marks a day whose rollup finished, so a retried job
	// can tell whether the day is already aggregated.
	MetricRollupDone = "rollup-done"
)

// RollupResult summarises what a daily rollup wrote.
type RollupResult struct {
	// Date is the day that was aggregated, in UTC.
	Date string
	// Metrics is the number of metric buckets written.
	Metrics int
}

// RollupDay aggregates one day of tickets and rule hits into the daily metric
// buckets. It is idempotent: recomputing a day overwrites the same buckets, so
// a retried job converges to the same state.
func (s *MetricService) RollupDay(ctx context.Context, date string) (RollupResult, error) {
	from, to, err := dayRange(date)
	if err != nil {
		return RollupResult{}, err
	}

	violations, err := s.hits.CountByLevel(ctx, from, to)
	if err != nil {
		return RollupResult{}, err
	}

	written := 0
	write := func(name string, value float64) error {
		key := domain.MetricKey{Name: name, Scope: domain.MetricScopeDaily, Bucket: date}
		if err := s.metric.SetGauge(ctx, key, value); err != nil {
			return err
		}
		written++
		return nil
	}

	levels := map[uint8]float64{}
	for _, item := range violations {
		levels[item.Level] = float64(item.Count)
	}
	if err := write(MetricViolationsError, levels[uint8(domain.RuleLevelError)]); err != nil {
		return RollupResult{}, err
	}
	if err := write(MetricViolationsWarning, levels[uint8(domain.RuleLevelWarning)]); err != nil {
		return RollupResult{}, err
	}
	if err := write(MetricViolationsNotice, levels[uint8(domain.RuleLevelMessage)]); err != nil {
		return RollupResult{}, err
	}

	// 工单类的每日计数由事件流写入，这里只补齐缺失的日子，
	// 让趋势图不会出现断点；已有计数的日子保持不动。
	for _, name := range []string{MetricTicketsCreated, MetricTicketsExecuted, MetricTicketsFailed} {
		key := domain.MetricKey{Name: name, Scope: domain.MetricScopeDaily, Bucket: date}
		if _, err := s.metric.GetByKey(ctx, key); err == nil {
			continue
		} else if !isNotFound(err) {
			return RollupResult{}, err
		}
		if err := s.metric.Create(ctx, &domain.Metric{
			UUID:     uuid.NewString(),
			Name:     name,
			Scope:    domain.MetricScopeDaily,
			Bucket:   date,
			Type:     domain.MetricTypeCounter,
			Value:    0,
			Version:  1,
			CreateAt: uint(time.Now().Unix()),
		}); err != nil {
			return RollupResult{}, err
		}
		written++
	}

	// 完成标记放在最后：中途失败时它不会出现，重跑时会重新覆盖当天。
	if err := write(MetricRollupDone, 1); err != nil {
		return RollupResult{}, err
	}

	return RollupResult{Date: date, Metrics: written}, nil
}

// RollupStatus reports whether a day was already aggregated.
func (s *MetricService) RollupStatus(ctx context.Context, date string) (bool, error) {
	key := domain.MetricKey{Name: MetricRollupDone, Scope: domain.MetricScopeDaily, Bucket: date}
	metric, err := s.metric.GetByKey(ctx, key)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return metric.Value > 0, nil
}

// dayRange returns the unix second range of one UTC day.
func dayRange(date string) (uint, uint, error) {
	day, err := time.ParseInLocation("2006-01-02", date, time.UTC)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: 日期格式应为 2006-01-02", domain.ErrInvalidInput)
	}
	start := day.Unix()
	return uint(start), uint(day.AddDate(0, 0, 1).Unix()), nil
}
