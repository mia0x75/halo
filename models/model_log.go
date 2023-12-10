package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Log 日志表，记录谁在什么时间做了什么事情
type Log struct {
	LogID     int    `xorm:"'log_id' notnull int pk autoincr"         valid:"-"                                json:"log_id"    gqlgen:"-"`         //
	UUID      string `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                json:"uuid"      gqlgen:"UUID"`      //
	UserID    uint   `xorm:"'user_id' notnull int index(index_1)"     valid:"required,int,range(0|4294967295)" json:"user_id"   gqlgen:"-"`         //
	Operation string `xorm:"'operation' notnull text"                 valid:"required,runelength(1|255)"       json:"operation" gqlgen:"Operation"` //
	Version   int    `xorm:"'version'"                                valid:"-"                                json:"version"   gqlgen:"-"`         //
	CreateAt  uint   `xorm:"'create_at' notnull int index(index_1)"   valid:"-"                                json:"create_at" gqlgen:"CreateAt"`  //
}

// TableName 结构体到数据库表名称的映射
func (m *Log) TableName() string {
	return "mm_logs"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Log) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Log) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Log) String() string {
	return fmt.Sprintf("uuid: %s, user_id: %d, operation: %s",
		m.UUID,
		m.UserID,
		m.Operation,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Log) IsNode() {}
