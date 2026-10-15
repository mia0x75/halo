package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// ClustersMap 线程安全的数据缓存对象，ClustersMap 群集缓存对象
var ClustersMap = &safecache.SafeMap[models.Cluster]{}

// Init 缓存初始化
func ClustersInit() (err error) {
	_, err = ClustersMap.WithInit(func() ([]*models.Cluster, error) {
		return loadIntKeyed(repos.Cluster.List, bridge.ClusterFromDomain)
	})
	return
}
