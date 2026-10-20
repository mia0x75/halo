package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Ticket 工单模型
type Ticket struct {
	TicketID   uint          `db:"ticket_id"      valid:"-"                                 json:"ticket_id"   gqlgen:"TicketID"`                  //
	UUID       string        `db:"uuid" valid:"-"                                 json:"uuid"        gqlgen:"UUID"`                                //
	ClusterID  uint          `db:"cluster_id"  valid:"required,int,range(0|4294967295)"  json:"cluster_id"  gqlgen:"-"`                            //
	Database   string        `db:"database"           valid:"required,length(1|50)"             json:"database"    gqlgen:"Database"`              //
	Subject    string        `db:"subject"            valid:"required,runelength(1|50)"         json:"subject"     gqlgen:"Subject"`               //
	Content    string        `db:"content"                   valid:"required,runelength(1|65535)"      json:"content"     gqlgen:"Content"`        //
	Status     uint8         `db:"status"                     valid:"required,matches(^([1-9]?[0-9])$)" json:"status"      gqlgen:"Status"`        // 状态 0-99
	UserID     uint          `db:"user_id"     valid:"required,int,range(0|4294967295)"  json:"user_id"     gqlgen:"-"`                            //
	ReviewerID uint          `db:"reviewer_id" valid:"required,int,range(0|4294967295)"  json:"reviewer_id" gqlgen:"-"`                            //
	CronID     sql.NullInt64 `db:"cron_id"     valid:"required,int,range(0|4294967295)"  json:"cron_id"     gqlgen:"-"`                            //
	Version    int           `db:"version"                                valid:"-"                                 json:"version"     gqlgen:"-"` //
	UpdateAt   uint          `db:"update_at"                  valid:"-"                                 json:"update_at"   gqlgen:"UpdateAt"`      //
	CreateAt   uint          `db:"create_at"                  valid:"-"                                 json:"create_at"   gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Ticket) TableName() string {
	return "mm_tickets"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Ticket) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CronID = sql.NullInt64{Int64: 0, Valid: false}
	m.Subject = strings.TrimSpace(m.Subject)
	m.Content = strings.TrimSpace(m.Content)
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Ticket) BeforeUpdate() {
	m.Subject = strings.TrimSpace(m.Subject)
	m.Content = strings.TrimSpace(m.Content)
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Ticket) String() string {
	return fmt.Sprintf("uuid: %s, subject: %s, database: %s, status: %d",
		m.UUID,
		m.Subject,
		m.Database,
		m.Status,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Ticket) IsNode() {}

// 创建时间
func (m *Ticket) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Ticket) GetUpdateAt() *uint {
	return &m.UpdateAt
}

// IsSearchable GraphQL的基类需要实现的接口，暂时不动
func (Ticket) IsSearchable() {}
