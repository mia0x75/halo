package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// UsersMap 用户缓存对象
var UsersMap = safecache.SafeMap[models.User]{}

// Init 缓存初始化
func UsersInit() (err error) {
	_, err = UsersMap.WithInit(func() ([]*models.User, error) {
		return loadIntKeyed(repos.User.List, bridge.UserFromDomain)
	})
	return
}
