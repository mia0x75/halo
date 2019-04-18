package caches

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Init 初始化缓存
func Init() {
	log.Info("[I] cache begin")

	log.Info("[I] #0 Avatars...")
	AvatarsMap.Init()

	log.Info("[I] #1 Clusters...")
	ClustersMap.Init()

	log.Info("[I] #2 Options...")
	OptionsMap.Init()

	log.Info("[I] #3 Roles...")
	RolesMap.Init()

	log.Info("[I] #4 Edges...")
	EdgesMap.Init()

	log.Info("[I] #5 Rules...")
	RulesMap.Init()

	log.Info("[I] #6 Glossaries...")
	GlossariesMap.Init()

	log.Info("[I] #7 Users...")
	UsersMap.Init()

	log.Info("[I] #8 Statistics...")
	StatisticsMap.Init()

	log.Info("[I] cache done")

	LoopInit()
}

// LoopInit 定期刷新缓存
func LoopInit() {
	go func() {
		d := time.Duration(10) * time.Second
		for range time.Tick(d) {
			EdgesMap.Init()
			UsersMap.Init()
			StatisticsMap.Init()
		}
	}()

	go func() {
		d := time.Duration(60) * time.Second
		for range time.Tick(d) {
			ClustersMap.Init()
			OptionsMap.Init()
			RulesMap.Init()
		}
	}()
}
