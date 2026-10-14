package caches

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

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
		InitFuncs:       []func() error{AvatarsInit, EdgesInit, UsersInit, StatisticsInit},
		RefreshFuncs:    []func() error{EdgesInit, UsersInit, StatisticsInit},
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

// Init 初始化缓存：先同步加载全部缓存，再启动周期刷新。
func Init() {
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
