package g

import (
	"github.com/akhenakh/statgo"
)

const CREDENTIAL_KEY = "Credential"

var (
	// GlobalStat 获取服务器当前系统环境的全局对象
	GlobalStat = statgo.NewStat()
	// Scheduler 定时任务的全局对象
	// Scheduler = crons.NewScheduler()
)
