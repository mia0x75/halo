package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// AvatarsMap 线程安全的数据缓存对象，AvatarsMap 头像缓存对象
var AvatarsMap = &SafeMap[models.Avatar]{}

// Init 缓存初始化
func AvatarsInit() (err error) {
	AvatarsMap.WithInit(
		func() (m []*models.Avatar, err error) {
			if err := g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "avatars", err.Error())
			}
			return
		},
	)
	return
}
