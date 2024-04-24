package caches

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type SafeMap[T any] struct {
	sync.RWMutex
	M         []*T
	fetchFunc func() ([]*T, error)
}

// Count 返回缓存条数
func (c *SafeMap[T]) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.M)
}

// Append 添加元素
func (c *SafeMap[T]) Append(item *T) {
	c.RLock()
	defer c.RUnlock()
	c.M = append(c.M, item)
}

// Remove 删除元素
func (c *SafeMap[T]) Remove(f func(*T) bool) {
	c.RLock()
	defer c.RUnlock()
	var s []*T
	for _, v := range c.M {
		if !f(v) {
			s = append(s, v)
		}
	}
	c.M = s
}

// Include returns true if one of the element in the sliece satisfies the predicate f.
func (c *SafeMap[T]) Include(f func(*T) bool) bool {
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
func (c *SafeMap[T]) Any(f func(*T) bool) *T {
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
func (c *SafeMap[T]) All() []*T {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter returns a new slice containing all elements in the slice that satisfy the predicate f.
func (c *SafeMap[T]) Filter(f func(*T) bool) (L []*T) {
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
func (c *SafeMap[T]) Map(f func(*T) *T) []*T {
	c.RLock()
	defer c.RUnlock()
	m := make([]*T, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// GetPage 返回一页缓存数据
func (c *SafeMap[T]) GetPage(offset, limit int) (users []*T) {
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

// Reload 重新加载数据
func (c *SafeMap[T]) Reload() {
	c.Lock()
	defer c.Unlock()
	if c.fetchFunc == nil {
		return
	}
	if m, err := c.fetchFunc(); err != nil {
		c.M = m
	}
}

// Init 缓存初始化
func (c *SafeMap[T]) WithInit(fn func() ([]*T, error)) []*T {
	if m, err := fn(); err != nil {
		log.Printf("查询数据表`%s`时发生一个错误:%s", "users", err.Error())
		return nil
	} else {
		c.Lock()
		defer c.Unlock()
		c.M = m
		c.fetchFunc = fn
		return m
	}
}
