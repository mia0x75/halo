package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Log 日志表，记录谁在什么时间做了什么事情
type Log struct {
	LogID     int    `db:"log_id"         valid:"-"                                json:"log_id"    gqlgen:"-"`                         //
	UUID      string `db:"uuid" valid:"-"                                json:"uuid"      gqlgen:"UUID"`                                //
	UserID    uint   `db:"user_id"     valid:"required,int,range(0|4294967295)" json:"user_id"   gqlgen:"-"`                            //
	Operation string `db:"operation"                 valid:"required,runelength(1|255)"       json:"operation" gqlgen:"Operation"`      //
	Version   int    `db:"version"                                valid:"-"                                json:"version"   gqlgen:"-"` //
	CreateAt  uint   `db:"create_at"   valid:"-"                                json:"create_at" gqlgen:"CreateAt"`                     //
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
