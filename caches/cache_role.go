package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// RolesMap 线程安全的数据缓存对象，RolesMap 角色缓存对象
var RolesMap = &SafeMap[models.Role]{}

// Init 缓存初始化
func RolesInit() {
	RolesMap.WithInit(
		func() (m []*models.Role, err error) {
			if err = g.Engine.Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "roles", err.Error())
			}
			return
		},
	)
}
