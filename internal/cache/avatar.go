package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// AvatarsMap 线程安全的数据缓存对象，AvatarsMap 头像缓存对象
var AvatarsMap = &safecache.SafeMap[models.Avatar]{}

// Init 缓存初始化
func AvatarsInit() (err error) {
	_, err = AvatarsMap.WithInit(func() ([]*models.Avatar, error) {
		return loadIntKeyed(repos.Avatar.List, bridge.AvatarFromDomain)
	})
	return
}
