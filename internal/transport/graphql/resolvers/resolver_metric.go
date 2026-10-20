package resolvers

import (
	"context"

	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/models"
)

// Metrics 按作用域与名称筛选指标
func (r *queryRootResolver) Metrics(ctx context.Context, scope *uint8, names []string) (L []*models.Metric, err error) {
	L = cache.MetricsMap.Filter(func(elem *models.Metric) bool {
		if scope != nil && elem.Scope != *scope {
			return false
		}
		if len(names) > 0 && !containsName(names, elem.Name) {
			return false
		}
		return true
	})

	return
}

// containsName reports whether the metric name is in the requested list.
func containsName(names []string, name string) bool {
	for _, item := range names {
		if item == name {
			return true
		}
	}
	return false
}
