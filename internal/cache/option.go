package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// OptionsMap 线程安全的数据缓存对象，OptionsMap 系统配置缓存对象
var OptionsMap = &safecache.SafeMap[models.Option]{}

// Init 缓存初始化
func OptionsInit() (err error) {
	_, err = OptionsMap.WithInit(func() ([]*models.Option, error) {
		return loadAll(repos.Option.List, bridge.OptionFromDomain)
	})
	return
}
