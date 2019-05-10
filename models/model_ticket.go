package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Ticket 工单模型
type Ticket struct {
	TicketID   uint          `xorm:"'ticket_id' notnull int pk autoincr"      valid:"-"                                 json:"ticket_id"   gqlgen:"-"`        //
	UUID       string        `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                 json:"uuid"        gqlgen:"UUID"`     //
	ClusterID  uint          `xorm:"'cluster_id' notnull int index(index_1)"  valid:"required,int,range(0|4294967295)"  json:"cluster_id"  gqlgen:"-"`        //
	Database   string        `xorm:"'database' notnull varchar(50)"           valid:"required,length(1|50)"             json:"database"    gqlgen:"Database"` //
	Subject    string        `xorm:"'subject' notnull varchar(50)"            valid:"required,runelength(1|50)"         json:"subject"     gqlgen:"Subject"`  //
	Content    string        `xorm:"'content' notnull text"                   valid:"required,runelength(1|65535)"      json:"content"     gqlgen:"Content"`  //
	Status     uint8         `xorm:"'status' notnull int"                     valid:"required,matches(^([1-9]?[0-9])$)" json:"status"      gqlgen:"Status"`   // 状态 0-99
	UserID     uint          `xorm:"'user_id' notnull int index(index_2)"     valid:"required,int,range(0|4294967295)"  json:"user_id"     gqlgen:"-"`        //
	ReviewerID uint          `xorm:"'reviewer_id' notnull int index(index_3)" valid:"required,int,range(0|4294967295)"  json:"reviewer_id" gqlgen:"-"`        //
	CronID     sql.NullInt64 `xorm:"'cron_id' notnull int index(index_4)"     valid:"required,int,range(0|4294967295)"  json:"cron_id"     gqlgen:"-"`        //
	Version    int           `xorm:"'version'"                                valid:"-"                                 json:"version"     gqlgen:"-"`        //
	UpdateAt   uint          `xorm:"'update_at' notnull int"                  valid:"-"                                 json:"update_at"   gqlgen:"UpdateAt"` //
	CreateAt   uint          `xorm:"'create_at' notnull int"                  valid:"-"                                 json:"create_at"   gqlgen:"CreateAt"` //
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

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Ticket) AfterSet(colName string, _ xorm.Cell) {
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

// IsSearchable GraphQL的基类需要实现的接口，暂时不动
func (Ticket) IsSearchable() {}
