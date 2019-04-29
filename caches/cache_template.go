package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeTemplatesMap 线程安全的数据缓存对象
type SafeTemplatesMap struct {
	sync.RWMutex
	M []*models.Template
}

// TemplatesMap 规则缓存对象
var TemplatesMap = &SafeTemplatesMap{}

// Count 返回缓存条数
func (c *SafeTemplatesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeTemplatesMap) Include(f func(*models.Template) bool) bool {
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
func (c *SafeTemplatesMap) Any(f func(*models.Template) bool) *models.Template {
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
func (c *SafeTemplatesMap) All() []*models.Template {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeTemplatesMap) Filter(f func(*models.Template) bool) (L []*models.Template) {
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
func (c *SafeTemplatesMap) Map(f func(*models.Template) *models.Template) []*models.Template {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Template, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Init 缓存初始化
func (c *SafeTemplatesMap) Init() {
	var m []*models.Template

	if err := g.Engine.Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "templates", err.Error())
		return
	}
	c.Lock()
	defer c.Unlock()
	c.M = m
}
