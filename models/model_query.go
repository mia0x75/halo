package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Query 记录用户发起一般数据查询及SHOW语句
type Query struct {
	QueryID   uint   `xorm:"'query_id' notnull int pk autoincr"       valid:"-"                                json:"query_id"    gqlgen:"-"`        //
	UUID      string `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                json:"uuid"        gqlgen:"UUID"`     //
	Type      uint8  `xorm:"'type' notnull tinyint"                   valid:"required,int,range(0|255)"        json:"type"        gqlgen:"Type"`     //
	ClusterID uint   `xorm:"'cluster_id' notnull int index(index_1)"  valid:"required,int,range(0|4294967295)" json:"cluster_id"  gqlgen:"-"`        //
	Database  string `xorm:"'database' notnull varchar(50)"           valid:"required,length(1|50)"            json:"database"    gqlgen:"Database"` //
	Content   string `xorm:"'content' notnull text"                   valid:"required,length(1|65535)"         json:"content"     gqlgen:"Content"`  //
	Plan      string `xorm:"'plan' notnull text"                      valid:"required,length(1|65535),ascii"   json:"plan"        gqlgen:"Plan"`     //
	UserID    uint   `xorm:"'user_id' notnull int index(index_2)"     valid:"required,int,range(0|4294967295)" json:"user_id"     gqlgen:"-"`        //
	Version   int    `xorm:"'version'"                                valid:"-"                                json:"version"     gqlgen:"-"`        //
	UpdateAt  uint   `xorm:"'update_at' notnull int"                  valid:"-"                                json:"update_at"   gqlgen:"UpdateAt"` //
	CreateAt  uint   `xorm:"'create_at' notnull int"                  valid:"-"                                json:"create_at"   gqlgen:"CreateAt"` //
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

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Query) AfterSet(colName string, _ xorm.Cell) {
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
