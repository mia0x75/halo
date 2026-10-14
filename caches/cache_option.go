package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// OptionsMap 线程安全的数据缓存对象，OptionsMap 系统配置缓存对象
var OptionsMap = &SafeMap[models.Option]{}

// Init 缓存初始化
func OptionsInit() (err error) {
	OptionsMap.WithInit(
		func() (m []*models.Option, err error) {
			if err = g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "options", err.Error())
			}
			return
		},
	)
	return
}
