package caches

import (
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// ClustersMap 线程安全的数据缓存对象，ClustersMap 群集缓存对象
var ClustersMap = &SafeMap[models.Cluster]{}

// GetPage 返回一页缓存数据
// func (c *SafeClustersMap) GetPage(offset, limit int) (clusters []*models.Cluster) {
// 	c.RLock()
// 	defer c.RUnlock()
// 	switch {
// 	case offset >= len(c.M) || offset < 0:
// 	case offset+int(limit) >= len(c.M):
// 		clusters = c.M[offset:]
// 	default:
// 		clusters = c.M[offset : offset+limit]
// 	}
// 	return
// }

// Init 缓存初始化
func ClustersInit() (err error) {
	ClustersMap.WithInit(
		func() (m []*models.Cluster, err error) {
			if err = g.Engine.Desc("cluster_id").Find(&m); err != nil {
				log.Printf("查询数据表`%s`时发生一个错误:%s", "clusters", err.Error())
			}
			return
		},
	)
	return
}
