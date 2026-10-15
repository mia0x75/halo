package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Avatar 用户可选的头像
type Avatar struct {
	AvatarID uint   `db:"avatar_id"     valid:"-"                     json:"avatar_id" gqlgen:"-"`                            //
	UUID     string `db:"uuid"   valid:"-"                     json:"uuid"      gqlgen:"UUID"`                                //
	URL      string `db:"url" valid:"required,length(1|75)" json:"url"       gqlgen:"URL"`                                    //
	Version  int    `db:"version"                                  valid:"-"                     json:"version"   gqlgen:"-"` //
	UpdateAt uint   `db:"update_at"                    valid:"-"                     json:"update_at" gqlgen:"UpdateAt"`      //
	CreateAt uint   `db:"create_at"                    valid:"-"                     json:"create_at" gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Avatar) TableName() string {
	return "mm_avatars"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Avatar) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Avatar) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Avatar) String() string {
	return fmt.Sprintf("uuid: %s, url: %s",
		m.UUID,
		m.URL,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Avatar) IsNode() {}

// 创建时间
func (m *Avatar) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Avatar) GetUpdateAt() *uint {
	return &m.UpdateAt
}
