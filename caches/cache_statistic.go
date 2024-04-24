package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// StatisticsMap 线程安全的数据缓存对象，StatisticsMap 规则缓存对象
var StatisticsMap = &SafeMap[models.Statistic]{}

// Init 缓存初始化
func StatisticsInit() {
	StatisticsMap.WithInit(
		func() (m []*models.Statistic, err error) {
			if err = g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "statistics", err.Error())
			}
			return
		},
	)
}
