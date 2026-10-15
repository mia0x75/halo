package caches_test

import (
	"testing"

	"github.com/mia0x75/halo/internal/platform/cache" // 请替换为您的实际包路径
)

type TestType int

func (t TestType) String() string {
	return "Test"
}

// 辅助函数，初始化缓存
func newCacheWithItems(items ...TestType) *caches.SafeMap[TestType] {
	cache := &caches.SafeMap[TestType]{M: []*TestType{}}
	for _, item := range items {
		cache.Append(&item) // 假设Append是添加元素的方法
	}
	return cache
}

func TestCount(t *testing.T) {
	cache := newCacheWithItems(TestType(1))
	count := cache.Count()

	if count != 1 {
		t.Errorf("Expected count to be 1, but got %d", count)
	}
}

func TestAppend(t *testing.T) {
	cache := newCacheWithItems()
	item := TestType(1)
	cache.Append(&item) // 使用正确的TestType实例

	count := cache.Count()

	if count != 1 {
		t.Errorf("Expected count to be 1 after appending, but got %d", count)
	}
}

// 提示：考虑增加针对Remove方法错误处理的测试用例
func TestRemove(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	cache := newCacheWithItems(item1, item2)
	cache.Remove(func(item *TestType) bool {
		return *item == 1
	})

	count := cache.Count()

	if count != 1 {
		t.Errorf("Expected count to be 1 after removing, but got %d", count)
	}
}

// 提示：考虑对Include方法添加测试，以覆盖更多边界条件
func TestInclude(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	cache := newCacheWithItems(item1, item2)
	includes := cache.Include(func(item *TestType) bool {
		return *item == 1
	})

	if !includes {
		t.Error("Expected includes to be true, but got false")
	}
}

// 提示：Any方法的测试用例可以优化为更明确地检查期望结果
func TestAny(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	cache := newCacheWithItems(item1, item2)
	any := cache.Any(func(item *TestType) bool {
		return *item == 1
	})

	if any == nil {
		t.Error("Expected any to return a non-nil pointer, but got nil")
	}
}

// 提示：All方法的测试用例可以添加更多边界条件，如空缓存的情况
func TestAll(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	cache := newCacheWithItems(item1, item2)
	all := cache.All()

	if len(all) != 2 {
		t.Errorf("Expected all to return 2 elements, but got %d", len(all))
	}
}

// 提示：Filter方法的测试用例可以考虑添加对空结果和错误处理的测试
func TestFilter(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	cache := newCacheWithItems(item1, item2)
	filtered := cache.Filter(func(item *TestType) bool {
		return *item == 1
	})

	if len(filtered) != 1 {
		t.Errorf("Expected filter to return 1 element, but got %d", len(filtered))
	}
}

// 提示：Map方法的测试用例可以考虑添加错误处理和边界条件的测试
func TestMap(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	cache := newCacheWithItems(item1, item2)
	mapped := cache.Map(func(item *TestType) *TestType {
		return item
	})

	expect1 := TestType(1)
	expect2 := TestType(2)
	expected := []*TestType{&expect1, &expect2}
	if len(mapped) != len(expected) {
		t.Errorf("Expected map to return %d elements, but got %d", len(expected), len(mapped))
	}

	for i, item := range mapped {
		if *item != *expected[i] {
			t.Errorf("Expected mapped element at index %d to be %d, but got %d", i, *expected[i], *item)
		}
	}
}

// 提示：GetPage方法的测试用例可以考虑添加边界条件，如页码超出范围的情况
func TestGetPage(t *testing.T) {
	item1 := TestType(1)
	item2 := TestType(2)
	item3 := TestType(3)
	item4 := TestType(4)
	item5 := TestType(5)
	cache := newCacheWithItems(item1, item2, item3, item4, item5)
	page := cache.GetPage(1, 3)

	expect1 := TestType(2)
	expect2 := TestType(3)
	expect3 := TestType(4)
	expected := []*TestType{&expect1, &expect2, &expect3}
	if len(page) != len(expected) {
		t.Errorf("Expected page to return %d elements, but got %d", len(expected), len(page))
	}

	for i, item := range page {
		if *item != *expected[i] {
			t.Errorf("Expected page element at index %d to be %d, but got %d", i, *expected[i], *item)
		}
	}
}

// 提示：Reload方法的测试用例可以考虑添加对错误处理的测试
func TestReload(t *testing.T) {
	cache := &caches.SafeMap[TestType]{}
	cache.Reload()

	if cache.Count() != 0 {
		t.Error("Expected count to be 0 after reload, but got non-zero")
	}
}

// 提示：WithInit方法的测试用例可以考虑添加失败场景的测试，如初始化函数返回错误
func TestWithInit(t *testing.T) {
	cache := &caches.SafeMap[TestType]{}

	item1 := TestType(1)
	item2 := TestType(2)
	item3 := TestType(3)
	items := []*TestType{&item1, &item2, &item3}
	initFunc := func() ([]*TestType, error) {
		return items, nil
	}

	m, err := cache.WithInit(initFunc)

	if err != nil {
		t.Errorf("Expected error to be nil, but got %s", err.Error())
	}

	if len(m) != len(items) {
		t.Errorf("Expected WithInit to return %d elements, but got %d", len(items), len(m))
	}

	for i, item := range m {
		if *item != *items[i] {
			t.Errorf("Expected WithInit element at index %d to be %d, but got %d", i, *items[i], *item)
		}
	}

	if cache.Count() != len(items) {
		t.Errorf("Expected count to be %d after WithInit, but got %d", len(items), cache.Count())
	}
}
