package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

// Metric names written by the daily rollup. The daily buckets hold a snapshot
// of one day, so they are gauges: recomputing a day overwrites the value
// instead of accumulating it.
const (
	MetricTicketsCreated    = "tickets-created"
	MetricTicketsApproved   = "tickets-approved"
	MetricTicketsExecuted   = "tickets-executed"
	MetricTicketsFailed     = "tickets-failed"
	MetricViolationsError   = "violations-error"
	MetricViolationsWarning = "violations-warning"
	MetricViolationsNotice  = "violations-notice"
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

	tickets, err := s.tickets.CountByStatus(ctx, from, to)
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

	counts := map[domain.TicketStatus]float64{}
	for _, item := range tickets {
		counts[item.Status] = float64(item.Count)
	}
	if err := write(MetricTicketsCreated, counts[domain.TicketStatusWaitingForVld]); err != nil {
		return RollupResult{}, err
	}
	if err := write(MetricTicketsApproved, counts[domain.TicketStatusLgtm]); err != nil {
		return RollupResult{}, err
	}
	if err := write(MetricTicketsExecuted, counts[domain.TicketStatusDone]); err != nil {
		return RollupResult{}, err
	}
	if err := write(MetricTicketsFailed, counts[domain.TicketStatusVldFailure]+counts[domain.TicketStatusExecFailure]); err != nil {
		return RollupResult{}, err
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
