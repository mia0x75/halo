package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeGlossariesMap 线程安全的数据缓存对象
type SafeGlossariesMap struct {
	sync.RWMutex
	M []*models.Glossary
}

// GlossariesMap 字典缓存对象
var GlossariesMap = &SafeGlossariesMap{}

// Count 返回缓存条数
func (c *SafeGlossariesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeGlossariesMap) Include(f func(*models.Glossary) bool) bool {
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
func (c *SafeGlossariesMap) Any(f func(*models.Glossary) bool) *models.Glossary {
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
func (c *SafeGlossariesMap) All() []*models.Glossary {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeGlossariesMap) Filter(f func(*models.Glossary) bool) (L []*models.Glossary) {
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
func (c *SafeGlossariesMap) Map(f func(*models.Glossary) *models.Glossary) []*models.Glossary {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Glossary, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Init 缓存初始化
func (c *SafeGlossariesMap) Init() {
	var m []*models.Glossary

	if err := g.Engine.Asc("group").Asc("key").Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "glossaries", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()
	c.M = m
}
