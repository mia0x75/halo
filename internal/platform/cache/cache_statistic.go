package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/models"
)

// StatisticsMap 线程安全的数据缓存对象，StatisticsMap 规则缓存对象
var StatisticsMap = &SafeMap[models.Statistic]{}

// Init 缓存初始化
func StatisticsInit() (err error) {
	_, err = StatisticsMap.WithInit(func() ([]*models.Statistic, error) {
		return loadUUIDKeyed(repos.Statistic.List, bridge.StatisticFromDomain)
	})
	return
}
