package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// GroupsMap 线程安全的数据缓存对象，GroupsMap 分组缓存对象
var GroupsMap = &safecache.SafeMap[models.Group]{}

// Init 缓存初始化
func GroupsInit() (err error) {
	_, err = GroupsMap.WithInit(func() ([]*models.Group, error) {
		return loadUUIDKeyed(repos.Group.List, bridge.GroupFromDomain)
	})
	return
}
