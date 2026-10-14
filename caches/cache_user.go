package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// UsersMap 用户缓存对象
var UsersMap = SafeMap[models.User]{}

// Init 缓存初始化
func UsersInit() (err error) {
	UsersMap.WithInit(
		func() (m []*models.User, err error) {
			if err = g.Engine.Desc("user_id").Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "users", err.Error())
			}
			return
		},
	)
	return
}
