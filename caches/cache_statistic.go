package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeStatisticsMap 线程安全的数据缓存对象
type SafeStatisticsMap struct {
	sync.RWMutex
	M []*models.Statistic
}

// StatisticsMap 规则缓存对象
var StatisticsMap = &SafeStatisticsMap{}

// Count 返回缓存条数
func (c *SafeStatisticsMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeStatisticsMap) Include(f func(*models.Statistic) bool) bool {
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
func (c *SafeStatisticsMap) Any(f func(*models.Statistic) bool) *models.Statistic {
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
func (c *SafeStatisticsMap) All() []*models.Statistic {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeStatisticsMap) Filter(f func(*models.Statistic) bool) (L []*models.Statistic) {
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
func (c *SafeStatisticsMap) Map(f func(*models.Statistic) *models.Statistic) []*models.Statistic {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Statistic, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// Init 缓存初始化
func (c *SafeStatisticsMap) Init() {
	var m []*models.Statistic

	if err := g.Engine.Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "statistics", err.Error())
		return
	}
	c.Lock()
	defer c.Unlock()
	c.M = m
}
