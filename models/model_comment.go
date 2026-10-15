package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Comment 审核意见建议的模型
type Comment struct {
	CommentID uint   `db:"comment_id"     valid:"-"                                json:"comment_id" gqlgen:"-"`                         //
	UUID      string `db:"uuid" valid:"-"                                json:"uuid"       gqlgen:"UUID"`                                //
	Content   string `db:"content"               valid:"required,length(1|255)"           json:"content"    gqlgen:"Content"`            //
	TicketID  uint   `db:"ticket_id"   valid:"required,int,range(0|4294967295)" json:"ticket_id"  gqlgen:"-"`                            //
	UserID    uint   `db:"user_id"     valid:"required,int,range(0|4294967295)" json:"user_id"    gqlgen:"-"`                            //
	Version   int    `db:"version"                                valid:"-"                                json:"version"    gqlgen:"-"` //
	UpdateAt  uint   `db:"update_at"                  valid:"-"                                json:"update_at"  gqlgen:"UpdateAt"`      //
	CreateAt  uint   `db:"create_at"                  valid:"-"                                json:"create_at"  gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Comment) TableName() string {
	return "mm_comments"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Comment) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.Content = strings.TrimSpace(m.Content)
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Comment) BeforeUpdate() {
	m.Content = strings.TrimSpace(m.Content)
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Comment) String() string {
	return fmt.Sprintf("uuid: %s, content: %s",
		m.UUID,
		m.Content,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Comment) IsNode() {}

// 创建时间
func (m *Comment) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Comment) GetUpdateAt() *uint {
	return &m.UpdateAt
}
