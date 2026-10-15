package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Template 用户模型
type Template struct {
	UUID        string `db:"uuid"        valid:"-"                            json:"uuid"      gqlgen:"UUID"`                  //
	Subject     string `db:"subject"    valid:"required,email,length(1|100)" json:"email"     gqlgen:"Subject"`                //
	Body        string `db:"body"               valid:"required,length(1|65535)"     json:"-"         gqlgen:"Body"`           //
	Description string `db:"description" valid:"required,length(1|50)"        json:"-"         gqlgen:"Description"`           //
	Version     int    `db:"version"                         valid:"-"                            json:"version"   gqlgen:"-"` //
	UpdateAt    uint   `db:"update_at"           valid:"-"                            json:"update_at" gqlgen:"UpdateAt"`      //
	CreateAt    uint   `db:"create_at"           valid:"-"                            json:"create_at" gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Template) TableName() string {
	return "mm_templates"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Template) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.Subject = strings.TrimSpace(m.Subject)
	m.Body = strings.TrimSpace(m.Body)
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Template) BeforeUpdate() {
	m.Subject = strings.TrimSpace(m.Subject)
	m.Body = strings.TrimSpace(m.Body)
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Template) String() string {
	return fmt.Sprintf("uuid: %s, subject: %s",
		m.UUID,
		m.Subject,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Template) IsNode() {}

// 创建时间
func (m *Template) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Template) GetUpdateAt() *uint {
	return &m.UpdateAt
}
