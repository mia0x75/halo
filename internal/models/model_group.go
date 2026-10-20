package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Group 群集分组模型
type Group struct {
	GroupID     uint   `db:"group_id" valid:"-"                              json:"group_id"   gqlgen:"-"`             //
	UUID        string `db:"uuid"      valid:"-"                              json:"uuid"        gqlgen:"UUID"`        //
	Name        string `db:"name"      valid:"required,length(1|75)"         json:"name"        gqlgen:"Name"`         //
	Description string `db:"description" valid:"-"                            json:"description" gqlgen:"Description"` //
	Status      uint8  `db:"status"    valid:"required"                       json:"status"      gqlgen:"Status"`      //
	Version     int    `db:"version"   valid:"-"                              json:"version"     gqlgen:"-"`           //
	UpdateAt    uint   `db:"update_at" valid:"-"                              json:"update_at"   gqlgen:"UpdateAt"`    //
	CreateAt    uint   `db:"create_at" valid:"-"                              json:"create_at"   gqlgen:"CreateAt"`    //
}

// TableName 结构体到数据库表名称的映射
func (m *Group) TableName() string {
	return "mm_groups"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Group) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Group) String() string {
	return fmt.Sprintf("uuid: %s, name: %s", m.UUID, m.Name)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Group) IsNode() {}

// 创建时间
func (m *Group) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Group) GetUpdateAt() *uint {
	return &m.UpdateAt
}
