package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Metric 指标模型
type Metric struct {
	MetricID uint    `db:"metric_id" valid:"-"                              json:"metric_id"   gqlgen:"-"`            //
	UUID     string  `db:"uuid"      valid:"-"                              json:"uuid"        gqlgen:"UUID"`         //
	Name     string  `db:"name"      valid:"required,length(1|75)"         json:"name"        gqlgen:"Name"`         //
	Scope    uint8   `db:"scope"     valid:"required"                       json:"scope"       gqlgen:"Scope"`        //
	ScopeID  string  `db:"scope_id"  valid:"-"                              json:"scope_id"    gqlgen:"ScopeID"`      //
	Type     uint8   `db:"type"      valid:"required"                       json:"type"        gqlgen:"Type"`         //
	Bucket   string  `db:"bucket"    valid:"-"                              json:"bucket"      gqlgen:"Bucket"`       //
	Value    float64 `db:"value"     valid:"required"                       json:"value"       gqlgen:"Value"`        //
	Version  int     `db:"version"   valid:"-"                              json:"version"     gqlgen:"-"`            //
	UpdateAt uint    `db:"update_at" valid:"-"                              json:"update_at"   gqlgen:"UpdateAt"`     //
	CreateAt uint    `db:"create_at" valid:"-"                              json:"create_at"   gqlgen:"CreateAt"`     //
}

// TableName 结构体到数据库表名称的映射
func (m *Metric) TableName() string {
	return "mm_metrics"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Metric) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Metric) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Metric) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, scope: %d, value: %f",
		m.UUID,
		m.Name,
		m.Scope,
		m.Value,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Metric) IsNode() {}

// 创建时间
func (m *Metric) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Metric) GetUpdateAt() *uint {
	return &m.UpdateAt
}
