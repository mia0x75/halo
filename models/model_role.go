package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Role 系统角色
type Role struct {
	RoleID      uint   `db:"role_id"           valid:"-"                         json:"role_id"     gqlgen:"-"`                            //
	UUID        string `db:"uuid"       valid:"-"                         json:"uuid"        gqlgen:"UUID"`                                //
	Name        string `db:"name" valid:"required,runelength(1|50)" json:"name"        gqlgen:"Name"`                                      //
	Description string `db:"description"           valid:"required,runelength(1|75)" json:"description" gqlgen:"Description"`              //
	Version     int    `db:"version"                                      valid:"-"                         json:"version"     gqlgen:"-"` //
	UpdateAt    uint   `db:"update_at"                        valid:"-"                         json:"update_at"   gqlgen:"UpdateAt"`      //
	CreateAt    uint   `db:"create_at"                        valid:"-"                         json:"create_at"   gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Role) TableName() string {
	return "mm_roles"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Role) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Role) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Role) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, description: %s",
		m.UUID,
		m.Name,
		m.Description,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Role) IsNode() {}

// 创建时间
func (m *Role) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Role) GetUpdateAt() *uint {
	return &m.UpdateAt
}
