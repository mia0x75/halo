package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
)

// SafeUsersMap 线程安全的数据缓存对象
type SafeUsersMap struct {
	sync.RWMutex
	M []*models.User
}

// UsersMap 用户缓存对象
var UsersMap = &SafeUsersMap{}

// Count 返回缓存条数
func (c *SafeUsersMap) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeUsersMap) Append(item *models.User) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeUsersMap) Include(f func(*models.User) bool) bool {
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
func (c *SafeUsersMap) Any(f func(*models.User) bool) *models.User {
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
func (c *SafeUsersMap) All() []*models.User {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeUsersMap) Filter(f func(*models.User) bool) (L []*models.User) {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			L = append(L, v)
		}
	}
	return
}

// Lookup 根据UUID返回一个用户信息
func (c *SafeUsersMap) Lookup(UUID string) *models.User {
	return c.Any(func(elem *models.User) bool {
		if elem.UUID == UUID {
			return true
		}
		return false
	})
}

// Map returns a new slice containing the results of applying the function f to each string in the original slice.
func (c *SafeUsersMap) Map(f func(*models.User) *models.User) []*models.User {
	c.RLock()
	defer c.RUnlock()
	m := make([]*models.User, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// GetPage 返回一页缓存数据
func (c *SafeUsersMap) GetPage(offset, limit int) (users []*models.User) {
	c.RLock()
	defer c.RUnlock()
	switch {
	case offset >= len(c.M) || offset < 0:
	case offset+int(limit) >= len(c.M):
		users = c.M[offset:]
	default:
		users = c.M[offset : offset+limit]
	}
	return
}

// Init 缓存初始化
func (c *SafeUsersMap) Init() {
	var m []*models.User

	if err := g.Engine.Desc("user_id").Find(&m); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "users", err.Error())
		return
	}

	c.Lock()
	defer c.Unlock()
	c.M = m
}
