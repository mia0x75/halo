package g

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

// Engine 全局XORM的引擎
var Engine *xorm.Engine

// InitDB 初始化数据库连接
func InitDB() (err error) {
	cfg := Config()
	Engine, err = xorm.NewEngine("mysql", cfg.Database.Addr)
	if err != nil {
		log.Fatalf("[F] open db fail: %v", err)
	}

	Engine.SetMaxIdleConns(cfg.Database.MaxIdle)
	Engine.SetMaxOpenConns(cfg.Database.MaxConnections)
	Engine.SetConnMaxLifetime(time.Duration(cfg.Database.WaitTimeout) * time.Second)

	err = Engine.DB().Ping()
	if err != nil {
		log.Fatalf("[F] ping db fail: %v", err)
	}
	return
}
