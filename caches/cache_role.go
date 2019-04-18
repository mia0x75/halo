package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeRolesMap 线程安全的数据缓存对象
type SafeRolesMap struct {
	sync.RWMutex
	M map[string]*models.Role
}

// RolesMap 角色缓存对象
var RolesMap = &SafeRolesMap{M: make(map[string]*models.Role)}

// Count 返回缓存条数
func (c *SafeRolesMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeRolesMap) Include(f func(*models.Role) bool) bool {
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
func (c *SafeRolesMap) Any(f func(*models.Role) bool) *models.Role {
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
func (c *SafeRolesMap) All() []*models.Role {
	c.RLock()
	defer c.RUnlock()
	return c.Map(func(elem *models.Role) *models.Role {
		return elem
	})
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeRolesMap) Filter(f func(*models.Role) bool) (L []*models.Role) {
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
func (c *SafeRolesMap) Map(f func(*models.Role) *models.Role) []*models.Role {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.Role, len(c.M))
	i := 0
	for _, v := range c.M {
		m[i] = f(v)
		i++
	}
	return m
}

// Init 缓存初始化
func (c *SafeRolesMap) Init() {
	var m []*models.Role

	if err := g.Engine.Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "roles", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()
	for _, role := range m {
		c.M[role.Name] = role
	}
}
