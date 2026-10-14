package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// TemplatesMap 线程安全的数据缓存对象，TemplatesMap 规则缓存对象
var TemplatesMap = SafeMap[models.Template]{}

// Init 缓存初始化
func TemplatesInit() (err error) {
	TemplatesMap.WithInit(
		func() (m []*models.Template, err error) {
			if err = g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "templates", err.Error())
			}
			return
		},
	)
	return
}
