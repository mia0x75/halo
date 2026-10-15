package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// GlossariesMap 线程安全的数据缓存对象，GlossariesMap 字典缓存对象
var GlossariesMap = &safecache.SafeMap[models.Glossary]{}

// Init 缓存初始化
func GlossariesInit() (err error) {
	_, err = GlossariesMap.WithInit(func() ([]*models.Glossary, error) {
		return loadUUIDKeyed(repos.Glossary.List, bridge.GlossaryFromDomain)
	})
	return
}
