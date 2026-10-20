package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Cron 计划任务的模型
type Cron struct {
	CronID    uint   `db:"cron_id"        valid:"-"                                   json:"cron_id"   gqlgen:"-"`                         //
	UUID      string `db:"uuid" valid:"-"                                   json:"uuid"      gqlgen:"UUID"`                                //
	Status    string `db:"status"                 valid:"required,matches(^(C|D|E|F|P|R|S)$)" json:"status"    gqlgen:"Status"`            //
	Name      string `db:"name"               valid:"-"                                   json:"name"      gqlgen:"Name"`                  //
	Cmd       string `db:"cmd"               valid:"-"                                   json:"cmd"       gqlgen:"Cmd"`                    //
	Params    string `db:"params"             valid:"-"                                   json:"params"    gqlgen:"Params"`                //
	Interval  string `db:"interval"           valid:"-"                                   json:"interval"  gqlgen:"Interval"`              //
	Duration  string `db:"duration"           valid:"-"                                   json:"duration"  gqlgen:"Duration"`              //
	LastRun   string `db:"last_run"           valid:"-"                                   json:"last_run"  gqlgen:"LastRun"`               //
	NextRun   string `db:"next_run"           valid:"-"                                   json:"next_run"  gqlgen:"NextRun"`               //
	Recurrent uint8  `db:"recurrent"              valid:"-"                                   json:"recurrent" gqlgen:"-"`                 //
	UserID    uint   `db:"user_id"                               valid:"-"                                   json:"user_id"  gqlgen:"-"`   //
	ClusterID uint   `db:"cluster_id"                            valid:"-"                                   json:"cluster_id" gqlgen:"-"` //
	Version   int    `db:"version"                                valid:"-"                                   json:"version"   gqlgen:"-"` //
	UpdateAt  uint   `db:"update_at"                  valid:"-"                                   json:"update_at" gqlgen:"UpdateAt"`      //
	CreateAt  uint   `db:"create_at"                  valid:"-"                                   json:"create_at" gqlgen:"CreateAt"`      //
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

// 创建时间
func (m *Cron) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Cron) GetUpdateAt() *uint {
	return &m.UpdateAt
}
