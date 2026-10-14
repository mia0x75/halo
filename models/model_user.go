package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User 用户模型
type User struct {
	UserID   uint   `xorm:"'user_id' notnull int pk autoincr" db:"user_id"            valid:"-"                                json:"user_id"   gqlgen:"-"`          //
	UUID     string `xorm:"'uuid' notnull char(36) unique(unique_1)" db:"uuid"     valid:"-"                                json:"uuid"      gqlgen:"UUID"`          //
	Email    string `xorm:"'email' notnull varchar(75) unique(unique_2)" db:"email" valid:"required,email,length(3|75)"      json:"email"     gqlgen:"Email"`        //
	Password string `xorm:"'password' notnull char(60)" db:"password"                  valid:"required"                         json:"-"         gqlgen:"-"`         //
	Status   uint8  `xorm:"'status' notnull tinyint" db:"status"                     valid:"required,matches(^(0|1)$)"        json:"status"    gqlgen:"Status"`      // 0 - 禁用 | 1 - 有效
	Name     string `xorm:"'name' notnull varchar(25)" db:"name"                   valid:"required,runelength(1|25)"        json:"name"      gqlgen:"Name"`          //
	Phone    uint64 `xorm:"'phone' bigint" db:"phone"                               valid:""                                 json:"phone"     gqlgen:"Phone"`        //
	AvatarID uint   `xorm:"'avatar_id' notnull int" db:"avatar_id"                      valid:"required,int,range(0|4294967295)" json:"avatar_id" gqlgen:"-"`        //
	Version  int    `xorm:"'version'" db:"version"                                    valid:"-"                                json:"version"   gqlgen:"-"`          //
	UpdateAt uint   `xorm:"'update_at' notnull int" db:"update_at"                      valid:"-"                                json:"update_at" gqlgen:"UpdateAt"` //
	CreateAt uint   `xorm:"'create_at' notnull int" db:"create_at"                      valid:"-"                                json:"create_at" gqlgen:"CreateAt"` //
}

// TableName 结构体到数据库表名称的映射
func (m *User) TableName() string {
	return "mm_users"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *User) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.Email = strings.TrimSpace(m.Email)
	m.Name = strings.TrimSpace(m.Name)
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *User) BeforeUpdate() {
	m.Email = strings.TrimSpace(m.Email)
	m.Name = strings.TrimSpace(m.Name)
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *User) String() string {
	return fmt.Sprintf("uuid: %s, email: %s, phone: %d, name: %s, status: %d",
		m.UUID,
		m.Email,
		m.Phone,
		m.Name,
		m.Status,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (User) IsNode() {}

// 创建时间
func (m *User) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *User) GetUpdateAt() *uint {
	return &m.UpdateAt
}

// IsSearchable GraphQL的基类需要实现的接口，暂时不动
func (User) IsSearchable() {}
