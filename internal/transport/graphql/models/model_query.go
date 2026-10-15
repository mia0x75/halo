package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Query 记录用户发起一般数据查询及SHOW语句
type Query struct {
	QueryID   uint   `db:"query_id"       valid:"-"                                json:"query_id"    gqlgen:"-"`                         //
	UUID      string `db:"uuid" valid:"-"                                json:"uuid"        gqlgen:"UUID"`                                //
	Type      uint8  `db:"type"                   valid:"required,int,range(0|255)"        json:"type"        gqlgen:"Type"`              //
	ClusterID uint   `db:"cluster_id"  valid:"required,int,range(0|4294967295)" json:"cluster_id"  gqlgen:"-"`                            //
	Database  string `db:"database"           valid:"required,length(1|50)"            json:"database"    gqlgen:"Database"`              //
	Content   string `db:"content"                   valid:"required,length(1|65535)"         json:"content"     gqlgen:"Content"`        //
	Plan      string `db:"plan"                      valid:"required,length(1|65535),ascii"   json:"plan"        gqlgen:"Plan"`           //
	UserID    uint   `db:"user_id"     valid:"required,int,range(0|4294967295)" json:"user_id"     gqlgen:"-"`                            //
	Version   int    `db:"version"                                valid:"-"                                json:"version"     gqlgen:"-"` //
	UpdateAt  uint   `db:"update_at"                  valid:"-"                                json:"update_at"   gqlgen:"UpdateAt"`      //
	CreateAt  uint   `db:"create_at"                  valid:"-"                                json:"create_at"   gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Query) TableName() string {
	return "mm_queries"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Query) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.Content = strings.TrimSpace(m.Content)
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Query) BeforeUpdate() {
	m.Content = strings.TrimSpace(m.Content)
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Query) String() string {
	return fmt.Sprintf("uuid: %s, type: %d, content: %s, plan: %s, database: %s",
		m.UUID,
		m.Type,
		m.Content,
		m.Plan,
		m.Database,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Query) IsNode() {}

// 创建时间
func (m *Query) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Query) GetUpdateAt() *uint {
	return &m.UpdateAt
}
