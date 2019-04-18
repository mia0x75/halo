package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Cron 计划任务的模型
type Cron struct {
	CronID    uint   `xorm:"'cron_id' notnull int pk autoincr"        valid:"-"                                   json:"cron_id"   gqlgen:"-"`        //
	UUID      string `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                   json:"uuid"      gqlgen:"UUID"`     //
	Status    string `xorm:"'status' notnull char(1)"                 valid:"required,matches(^(C|D|E|F|P|R|S)$)" json:"status"    gqlgen:"Status"`   //
	Name      string `xorm:"'name' notnull varchar(75)"               valid:"-"                                   json:"name"      gqlgen:"Name"`     //
	Cmd       string `xorm:"'cmd' notnull varchar(100)"               valid:"-"                                   json:"cmd"       gqlgen:"Cmd"`      //
	Params    string `xorm:"'params' notnull varchar(75)"             valid:"-"                                   json:"params"    gqlgen:"Params"`   //
	Interval  string `xorm:"'interval' notnull varchar(20)"           valid:"-"                                   json:"interval"  gqlgen:"Interval"` //
	Duration  string `xorm:"'duration' notnull varchar(20)"           valid:"-"                                   json:"duration"  gqlgen:"Duration"` //
	LastRun   string `xorm:"'last_run' notnull varchar(20)"           valid:"-"                                   json:"last_run"  gqlgen:"LastRun"`  //
	NextRun   string `xorm:"'next_run' notnull varchar(20)"           valid:"-"                                   json:"next_run"  gqlgen:"NextRun"`  //
	Recurrent uint8  `xorm:"'recurrent' notnull tinyint"              valid:"-"                                   json:"recurrent" gqlgen:"-"`        //
	Version   int    `xorm:"'version'"                                valid:"-"                                   json:"version"   gqlgen:"-"`        //
	UpdateAt  uint   `xorm:"'update_at' notnull int"                  valid:"-"                                   json:"update_at" gqlgen:"UpdateAt"` //
	CreateAt  uint   `xorm:"'create_at' notnull int"                  valid:"-"                                   json:"create_at" gqlgen:"CreateAt"` //
}

// TableName 结构体到数据库表名称的映射
func (m *Cron) TableName() string {
	return "mm_crons"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Cron) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Cron) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Cron) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Cron) String() string {
	return fmt.Sprintf("uuid: %s, status: %s, name:%s, cmd: %s, duration: %s, last_run: %s, next_run: %s, recurrent: %d",
		m.UUID,
		m.Status,
		m.Name,
		m.Cmd,
		m.Duration,
		m.LastRun,
		m.NextRun,
		m.Recurrent,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Cron) IsNode() {}
