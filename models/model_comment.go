package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Comment 审核意见建议的模型
type Comment struct {
	CommentID uint   `xorm:"'comment_id' notnull int pk autoincr"     valid:"-"                                json:"comment_id" gqlgen:"-"`        //
	UUID      string `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                json:"uuid"       gqlgen:"UUID"`     //
	Content   string `xorm:"'content' notnull tinytext"               valid:"required,length(1|255)"           json:"content"    gqlgen:"Content"`  //
	TicketID  uint   `xorm:"'ticket_id' notnull int index(index_1)"   valid:"required,int,range(0|4294967295)" json:"ticket_id"  gqlgen:"-"`        //
	UserID    uint   `xorm:"'user_id' notnull int index(index_2)"     valid:"required,int,range(0|4294967295)" json:"user_id"    gqlgen:"-"`        //
	Version   int    `xorm:"'version'"                                valid:"-"                                json:"version"    gqlgen:"-"`        //
	UpdateAt  uint   `xorm:"'update_at' notnull int"                  valid:"-"                                json:"update_at"  gqlgen:"UpdateAt"` //
	CreateAt  uint   `xorm:"'create_at' notnull int"                  valid:"-"                                json:"create_at"  gqlgen:"CreateAt"` //
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

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Comment) AfterSet(colName string, _ xorm.Cell) {
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
