package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Role 系统角色
type Role struct {
	RoleID      uint   `xorm:"'role_id' notnull int pk autoincr <-"           valid:"-"                         json:"role_id"     gqlgen:"-"`           //
	UUID        string `xorm:"'uuid' notnull char(36) unique(unique_1)"       valid:"-"                         json:"uuid"        gqlgen:"UUID"`        //
	Name        string `xorm:"'name' notnull varchar(50) unique(unique_2) <-" valid:"required,runelength(1|50)" json:"name"        gqlgen:"Name"`        //
	Description string `xorm:"'description' notnull varchar(75) <-"           valid:"required,runelength(1|75)" json:"description" gqlgen:"Description"` //
	Version     int    `xorm:"'version'"                                      valid:"-"                         json:"version"     gqlgen:"-"`           //
	UpdateAt    uint   `xorm:"'update_at' notnull int"                        valid:"-"                         json:"update_at"   gqlgen:"UpdateAt"`    //
	CreateAt    uint   `xorm:"'create_at' notnull int"                        valid:"-"                         json:"create_at"   gqlgen:"CreateAt"`    //
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

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Role) AfterSet(colName string, _ xorm.Cell) {
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
