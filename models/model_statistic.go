package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Statistic 统计模型
type Statistic struct {
	Group    string  `db:"group"              valid:"-"                              json:"group"        gqlgen:"Group"`                 //
	Key      string  `db:"key"             valid:"required,length(1|50),alphanum" json:"name"         gqlgen:"Key"`                      //
	UUID     string  `db:"uuid" valid:"-"                              json:"uuid"         gqlgen:"UUID"`                                //
	Value    float64 `db:"value"            valid:"required"                       json:"value"        gqlgen:"Value"`                   //
	Version  int     `db:"version"                                valid:"-"                              json:"version"      gqlgen:"-"` //
	UpdateAt uint    `db:"update_at"                  valid:"-"                              json:"update_at"    gqlgen:"UpdateAt"`      //
	CreateAt uint    `db:"create_at"                  valid:"-"                              json:"create_at"    gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Statistic) TableName() string {
	return "mm_statistics"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Statistic) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Statistic) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Statistic) String() string {
	return fmt.Sprintf("uuid: %s, group: %s, name: %s, value: %f",
		m.UUID,
		m.Group,
		m.Key,
		m.Value,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Statistic) IsNode() {}

// 创建时间
func (m *Statistic) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Statistic) GetUpdateAt() *uint {
	return &m.UpdateAt
}
