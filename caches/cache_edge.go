package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// EdgesMap 线程安全的数据缓存对象，EdgesMap 关联缓存对象
var EdgesMap = &SafeMap[models.Edge]{}

// Init 缓存初始化
func EdgesInit() {
	EdgesMap.WithInit(
		func() (m []*models.Edge, err error) {
			if err = g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "edges", err.Error())
			}
			return
		},
	)
}
