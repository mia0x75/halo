package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

// CleanupResult reports how many rows a retention run removed.
type CleanupResult struct {
	// Metrics is the number of removed daily metric buckets.
	Metrics int64
	// Violations is the number of removed rule hits.
	Violations int64
}

// Cleanup removes the data older than the given number of days: the daily
// metric buckets and the rule hit details. Both grow without bound, so a
// retention window is needed; the platform keeps the detail of recent hits
// because the dashboard aggregates over it.
func (s *MetricService) Cleanup(ctx context.Context, days int) (CleanupResult, error) {
	if days <= 0 {
		return CleanupResult{}, fmt.Errorf("%w: 保留天数必须大于 0", domain.ErrInvalidInput)
	}

	before := uint(time.Now().AddDate(0, 0, -days).Unix())
	bucket := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	result := CleanupResult{}
	metrics, err := s.metric.DeleteBucketsBefore(ctx, bucket)
	if err != nil {
		return result, err
	}
	result.Metrics = metrics

	violations, err := s.hits.DeleteBefore(ctx, before)
	if err != nil {
		return result, err
	}
	result.Violations = violations

	return result, nil
}
