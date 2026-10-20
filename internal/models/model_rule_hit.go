package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RuleHit 规则命中明细模型
type RuleHit struct {
	HitID    uint   `db:"violation_id" valid:"-"                              json:"hit_id"    gqlgen:"-"`        //
	UUID     string `db:"uuid"         valid:"-"                              json:"uuid"      gqlgen:"UUID"`     //
	TicketID uint   `db:"ticket_id"    valid:"required"                       json:"ticket_id" gqlgen:"TicketID"` //
	Sequence uint16 `db:"sequence"     valid:"-"                              json:"sequence"  gqlgen:"Sequence"` //
	RuleName string `db:"rule_name"    valid:"required,length(1|75)"         json:"rule"      gqlgen:"RuleName"`  //
	Level    uint8  `db:"level"        valid:"required"                       json:"level"     gqlgen:"Level"`    //
	Message  string `db:"message"      valid:"required,length(1|150)"        json:"message"   gqlgen:"Message"`   //
	CreateAt uint   `db:"create_at"    valid:"-"                              json:"create_at" gqlgen:"CreateAt"` //
	UpdateAt uint   `db:"-"            valid:"-"                              json:"update_at" gqlgen:"UpdateAt"` //
}

// TableName 结构体到数据库表名称的映射
func (m *RuleHit) TableName() string {
	return "mm_violations"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *RuleHit) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *RuleHit) String() string {
	return fmt.Sprintf("uuid: %s, rule: %s, level: %d", m.UUID, m.RuleName, m.Level)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (RuleHit) IsNode() {}

// 创建时间
func (m *RuleHit) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *RuleHit) GetUpdateAt() *uint {
	return &m.UpdateAt
}
