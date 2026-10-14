package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// GlossariesMap 线程安全的数据缓存对象，GlossariesMap 字典缓存对象
var GlossariesMap = &SafeMap[models.Glossary]{}

// Init 缓存初始化
func GlossariesInit() (err error) {
	GlossariesMap.WithInit(
		func() (m []*models.Glossary, err error) {
			if err = g.Engine.Asc("group").Asc("key").Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "glossaries", err.Error())
			}
			return
		},
	)
	return
}
