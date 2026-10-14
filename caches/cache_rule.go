package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// RulesMap 线程安全的数据缓存对象，RulesMap 规则缓存对象
var RulesMap = &SafeMap[models.Rule]{}

// Init 缓存初始化
func RulesInit() (err error) {
	_, err = RulesMap.WithInit(
		func() (m []*models.Rule, err error) {
			if err = g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "rules", err.Error())
			}
			return
		},
	)
	return
}
