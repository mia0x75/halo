package caches

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Init 初始化缓存
func Init() {
	log.Info("[I] cache begin")

	log.Info("[I] #0 Avatars...")
	AvatarsInit()

	log.Info("[I] #1 Clusters...")
	ClustersInit()

	log.Info("[I] #2 Options...")
	OptionsInit()

	log.Info("[I] #3 Roles...")
	RolesInit()

	log.Info("[I] #4 Edges...")
	EdgesInit()

	log.Info("[I] #5 Rules...")
	RulesInit()

	log.Info("[I] #6 Glossaries...")
	GlossariesInit()

	log.Info("[I] #7 Users...")
	UsersInit()

	log.Info("[I] #8 Statistics...")
	StatisticsInit()

	log.Info("[I] #9 Templates...")
	TemplatesInit()

	log.Info("[I] cache done")

	LoopInit()
}

// LoopInit 定期刷新缓存
func LoopInit() {
	go func() {
		d := time.Duration(10) * time.Second
		for range time.Tick(d) {
			EdgesInit()
			UsersInit()
			StatisticsInit()
		}
	}()

	go func() {
		d := time.Duration(60) * time.Second
		for range time.Tick(d) {
			ClustersInit()
			OptionsInit()
			RulesInit()
			TemplatesInit()
		}
	}()
}
