package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeOptionsMap 线程安全的数据缓存对象
type SafeOptionsMap struct {
	sync.RWMutex
	M []*models.Option
}

// OptionsMap 系统配置缓存对象
var OptionsMap = &SafeOptionsMap{}

// Count 返回缓存条数
func (c *SafeOptionsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeOptionsMap) Include(f func(*models.Option) bool) bool {
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
func (c *SafeOptionsMap) Any(f func(*models.Option) bool) *models.Option {
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
func (c *SafeOptionsMap) All() []*models.Option {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeOptionsMap) Filter(f func(*models.Option) bool) (L []*models.Option) {
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
func (c *SafeOptionsMap) Map(f func(*models.Option) *models.Option) []*models.Option {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Option, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Init 缓存初始化
func (c *SafeOptionsMap) Init() {
	var m []*models.Option

	if err := g.Engine.Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "options", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()
	c.M = m
}
