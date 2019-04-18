package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeAvatarsMap 线程安全的数据缓存对象
type SafeAvatarsMap struct {
	sync.RWMutex
	M []*models.Avatar
}

// AvatarsMap 头像缓存对象
var AvatarsMap = &SafeAvatarsMap{}

// Count 返回缓存条数
func (c *SafeAvatarsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeAvatarsMap) Include(f func(*models.Avatar) bool) bool {
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
func (c *SafeAvatarsMap) Any(f func(*models.Avatar) bool) *models.Avatar {
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
func (c *SafeAvatarsMap) All() []*models.Avatar {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeAvatarsMap) Filter(f func(*models.Avatar) bool) (L []*models.Avatar) {
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
func (c *SafeAvatarsMap) Map(f func(*models.Avatar) *models.Avatar) []*models.Avatar {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Avatar, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Init 缓存初始化
func (c *SafeAvatarsMap) Init() {
	var m []*models.Avatar

	if err := g.Engine.Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "avatars", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()
	c.M = m
}
