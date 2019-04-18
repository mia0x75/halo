package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeClustersMap 线程安全的数据缓存对象
type SafeClustersMap struct {
	sync.RWMutex
	M []*models.Cluster
}

// ClustersMap 群集缓存对象
var ClustersMap = &SafeClustersMap{}

// Count 返回缓存条数
func (c *SafeClustersMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeClustersMap) Append(item *models.Cluster) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素，每次仅删除一个
func (c *SafeClustersMap) Remove(f func(*models.Cluster) bool) {
	c.RLock()
	defer c.RUnlock()
	for i, cluster := range c.M {
		if f(cluster) {
			c.M = append(c.M[:i], c.M[i+1:]...)
			break
		}
	}
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeClustersMap) Include(f func(*models.Cluster) bool) bool {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			return true
		}
	}
	return false
}

// Any returns the element if one of the element in the sliece satisfies the predicate f.
func (c *SafeClustersMap) Any(f func(*models.Cluster) bool) *models.Cluster {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			return v
		}
	}
	return nil
}

// All returns all of the slice.
func (c *SafeClustersMap) All() []*models.Cluster {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeClustersMap) Filter(f func(*models.Cluster) bool) (L []*models.Cluster) {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			L = append(L, v)
		}
	}
	return
}

// Map returns a new slice containing the results of applying the function f to each string in the original slice.
func (c *SafeClustersMap) Map(f func(*models.Cluster) *models.Cluster) []*models.Cluster {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Cluster, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// GetPage 返回一页缓存数据
func (c *SafeClustersMap) GetPage(offset, limit int) (clusters []*models.Cluster) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		clusters = c.M[offset:]
	default:
		clusters = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeClustersMap) Init() {
	var m []*models.Cluster

	if err := g.Engine.Desc("cluster_id").Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "clusters", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()

	c.M = m
}
