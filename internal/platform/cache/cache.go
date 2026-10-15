package cache

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// SafeMap 是一个并发安全的映射，存储切片类型的元素。
// 它提供了多种操作切片的方法，包括添加、删除、过滤等。
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

// Append 添加元素到缓存。
func (c *SafeMap[T]) Append(item *T) {
	c.Lock()
	defer c.Unlock()
	c.M = append(c.M, item)
}

// Remove 删除满足条件的元素。
func (c *SafeMap[T]) Remove(f func(*T) bool) {
	c.Lock()
	defer c.Unlock()
	var s []*T
	for _, v := range c.M {
		if !f(v) {
			s = append(s, v)
		}
	}
	c.M = s
}

// Include 判断切片中是否存在满足条件的元素。
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

// Any 返回切片中第一个满足条件的元素。
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

// All 返回切片中的所有元素。
func (c *SafeMap[T]) All() []*T {
	c.RLock()
	defer c.RUnlock()
	return c.M
}

// Filter 返回一个包含所有满足条件的元素的新切片。
func (c *SafeMap[T]) Filter(f func(*T) bool) (L []*T) {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.M {
		if f(v) {
			L = append(L, v)
		}
	}
	return L
}

// Map 根据函数f的结果，返回一个新的切片。
func (c *SafeMap[T]) Map(f func(*T) *T) []*T {
	c.RLock()
	defer c.RUnlock()
	m := make([]*T, len(c.M))
	for i, v := range c.M {
		m[i] = f(v)
	}
	return m
}

// GetPage 返回指定偏移量和限制大小的一页数据。
func (c *SafeMap[T]) GetPage(offset, limit int) (users []*T) {
	c.RLock()
	defer c.RUnlock()
	if offset < 0 || offset >= len(c.M) || limit <= 0 {
		return nil
	}
	end := offset + limit
	if end > len(c.M) {
		end = len(c.M)
	}
	return c.M[offset:end]
}

// Reload 重新加载数据。如果fetchFunc为nil，则不执行任何操作。
func (c *SafeMap[T]) Reload() {
	c.Lock()
	defer c.Unlock()
	if c.fetchFunc == nil {
		log.Println("Reload called but fetchFunc is nil")
		return
	}
	m, err := c.fetchFunc()
	if err != nil {
		log.Printf("Error reloading data: %s", err)
		return
	}
	c.M = m
}

// Init 初始化缓存。如果初始化失败，将返回nil并记录错误。
func (c *SafeMap[T]) WithInit(fn func() ([]*T, error)) (m []*T, err error) {
	if m, err = fn(); err != nil {
		log.Printf("Error initializing cache: %s", err)
		return
	} else {
		c.Lock()
		defer c.Unlock()
		c.M = m
		c.fetchFunc = fn
		return
	}
}
