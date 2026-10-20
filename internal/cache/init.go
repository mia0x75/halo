package cache

import (
	"context"
	"math"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/service"
)

// repos is injected once at startup via Init. The caches read their rows
// through the repository layer.

var repos service.Repos

// The caches hold whole tables in memory. List is keyset-paginated, so a single
// page is requested with a cursor that sorts after every key.
const (
	cacheLimit     = 1 << 20
	cacheIntCursor = int64(math.MaxUint32)
)

// cacheUUIDCursor sorts after every UUID-keyed row. The uuid columns hold hex
// UUIDs except mm_options, whose seed inserts values positionally and so stores
// option names ("ldap.domain") in the uuid column; 'z' repeated past the widest
// column sorts after both under the connection collation.
var cacheUUIDCursor = strings.Repeat("z", 64)

// loadIntKeyed pages in every row of an integer-keyed table and converts it to
// the legacy model the cache stores.
func loadIntKeyed[D any, M any](list func(context.Context, int64, int) ([]D, error), conv func(D) *M) ([]*M, error) {
	items, err := list(context.Background(), cacheIntCursor, cacheLimit)
	if err != nil {
		return nil, err
	}
	out := make([]*M, 0, len(items))
	for _, item := range items {
		out = append(out, conv(item))
	}
	return out, nil
}

// loadUUIDKeyed pages in every row of a UUID-keyed table.
func loadUUIDKeyed[D any, M any](list func(context.Context, string, int) ([]D, error), conv func(D) *M) ([]*M, error) {
	items, err := list(context.Background(), cacheUUIDCursor, cacheLimit)
	if err != nil {
		return nil, err
	}
	out := make([]*M, 0, len(items))
	for _, item := range items {
		out = append(out, conv(item))
	}
	return out, nil
}

// CacheInit 描述一组共享同一刷新间隔的缓存。
// InitFuncs 在启动时同步加载一次；RefreshFuncs 此后按 RefreshInterval 周期刷新。
// RefreshInterval 为 0 表示该组启动后不再刷新。
type CacheInit struct {
	Name            string
	InitFuncs       []func() error
	RefreshFuncs    []func() error
	RefreshInterval time.Duration
}

var cacheInits = []CacheInit{
	{
		Name:            "fast",
		InitFuncs:       []func() error{AvatarsInit, EdgesInit, UsersInit, MetricsInit},
		RefreshFuncs:    []func() error{EdgesInit, UsersInit, MetricsInit},
		RefreshInterval: 10 * time.Second,
	},
	{
		Name:            "slow",
		InitFuncs:       []func() error{ClustersInit, OptionsInit, RolesInit, RulesInit, GlossariesInit, TemplatesInit},
		RefreshFuncs:    []func() error{ClustersInit, OptionsInit, RulesInit, TemplatesInit},
		RefreshInterval: 60 * time.Second,
	},
}

var ctx, cancel = context.WithCancel(context.Background())

// LoopInit 停止所有的刷新循环，释放资源。应在程序退出时调用。
func LoopInit() {
	cancel()
}

// Init 注入仓储并初始化缓存：先同步加载全部缓存，再启动周期刷新。
func Init(r service.Repos) {
	repos = r
	log.Info("[I] cache begin")
	for _, ci := range cacheInits {
		for _, fn := range ci.InitFuncs {
			if err := fn(); err != nil {
				log.Errorf("[E] Error initializing cache: %v", err)
			}
		}
	}
	log.Info("[I] cache done")

	startRefreshLoops()
}

// startRefreshLoop 启动单个刷新循环
func startRefreshLoop(rf func() error, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := rf(); err != nil {
					log.Errorf("[E] Error refreshing cache: %v", err)
				}
			}
		}
	}()
}

// startRefreshLoops 启动所有的刷新循环，每个刷新函数恰好启动一个
func startRefreshLoops() {
	for _, ci := range cacheInits {
		if ci.RefreshInterval <= 0 {
			continue
		}
		for _, rf := range ci.RefreshFuncs {
			startRefreshLoop(rf, ci.RefreshInterval)
		}
	}
}
