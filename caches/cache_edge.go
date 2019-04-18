package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeEdgesMap 线程安全的数据缓存对象
type SafeEdgesMap struct {
	sync.RWMutex
	M []*models.Edge
}

// EdgesMap 关联缓存对象
var EdgesMap = &SafeEdgesMap{}

// Count 返回缓存条数
func (c *SafeEdgesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Reload 刷新缓存
func (c *SafeEdgesMap) Reload() {
	c.Init()
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeEdgesMap) Include(f func(*models.Edge) bool) bool {
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
func (c *SafeEdgesMap) Any(f func(*models.Edge) bool) *models.Edge {
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
func (c *SafeEdgesMap) All() []*models.Edge {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeEdgesMap) Filter(f func(*models.Edge) bool) (L []*models.Edge) {
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
func (c *SafeEdgesMap) Map(f func(*models.Edge) *models.Edge) []*models.Edge {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Edge, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Init 缓存初始化
func (c *SafeEdgesMap) Init() {
	var m []*models.Edge

	if err := g.Engine.Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "edges", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()
	c.M = m
}
