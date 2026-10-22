package cache

import (
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	safecache "github.com/mia0x75/halo/internal/platform/cache"
)

// RulesMap 线程安全的数据缓存对象，RulesMap 规则缓存对象
var RulesMap = &safecache.SafeMap[models.Rule]{}

// Init 缓存初始化
func RulesInit() (err error) {
	_, err = RulesMap.WithInit(func() ([]*models.Rule, error) {
		// The rule list is ordered by rule ID, so the cache starts at the first
		// rule instead of the UUID sentinel.
		return loadFrom(repos.Rule.List, "", bridge.RuleFromDomain)
	})
	return
}
