package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// RolesMap 线程安全的数据缓存对象，RolesMap 角色缓存对象
var RolesMap = &safecache.SafeMap[models.Role]{}

// Init 缓存初始化
func RolesInit() (err error) {
	_, err = RolesMap.WithInit(func() ([]*models.Role, error) {
		return loadIntKeyed(repos.Role.List, bridge.RoleFromDomain)
	})
	return
}
