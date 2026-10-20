package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// MetricsMap 线程安全的数据缓存对象，MetricsMap 指标缓存对象
var MetricsMap = &safecache.SafeMap[models.Metric]{}

// Init 缓存初始化
func MetricsInit() (err error) {
	_, err = MetricsMap.WithInit(func() ([]*models.Metric, error) {
		return loadUUIDKeyed(repos.Metric.List, bridge.MetricFromDomain)
	})
	return
}
