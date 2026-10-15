package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/models"
)

// EdgesMap 线程安全的数据缓存对象，EdgesMap 关联缓存对象
var EdgesMap = &SafeMap[models.Edge]{}

// Init 缓存初始化
func EdgesInit() (err error) {
	_, err = EdgesMap.WithInit(func() ([]*models.Edge, error) {
		return loadIntKeyed(repos.Edge.List, bridge.EdgeFromDomain)
	})
	return
}
