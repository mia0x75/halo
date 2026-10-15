package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Rule 系统审核规则
type Rule struct {
	Name        string `db:"name"               valid:"required,length(10|10)"            json:"name"        gqlgen:"Name"`                  //
	UUID        string `db:"uuid" valid:"-"                                 json:"uuid"        gqlgen:"UUID"`                                //
	Group       uint8  `db:"group"               valid:"required,range(0|255)"             json:"group"       gqlgen:"Group"`                //
	Level       uint8  `db:"level"               valid:"required,matches(^(1|2|3)$)"       json:"level"       gqlgen:"Level"`                //
	VldrGroup   uint16 `db:"vldr_group"         valid:"required,matches(^([0-9]{4})$)"    json:"vldr_group"  gqlgen:"VldrGroup"`             //
	Operator    string `db:"operator"         valid:"required,length(1|10),alphanum"    json:"operator"    gqlgen:"Operator"`                //
	Values      string `db:"values"            valid:"required,length(1|150)"            json:"values"      gqlgen:"Values"`                 //
	Bitwise     uint8  `db:"bitwise"                valid:"required,int,matches(^(4|5|6|7)$)" json:"bitwise"     gqlgen:"Bitwise"`           //
	Func        string `db:"func"            valid:"required,length(1|75)"             json:"func"        gqlgen:"Func"`                     //
	Message     string `db:"message"        valid:"required,runelength(1|150)"        json:"message"     gqlgen:"Message"`                   //
	Description string `db:"description"        valid:"required,runelength(1|255)"        json:"description" gqlgen:"Description"`           //
	Element     string `db:"element"         valid:"required,length(1|50)"             json:"element"     gqlgen:"Element"`                  //
	Version     int    `db:"version"                                valid:"-"                                 json:"version"     gqlgen:"-"` //
	UpdateAt    uint   `db:"update_at"                  valid:"-"                                 json:"update_at"   gqlgen:"UpdateAt"`      //
	CreateAt    uint   `db:"create_at"                  valid:"-"                                 json:"create_at"   gqlgen:"CreateAt"`      //
}

// TableName 结构体到数据库表名称的映射
func (m *Rule) TableName() string {
	return "mm_rules"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Rule) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Rule) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// String 结构体输出到字符串的默认方式
func (m *Rule) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, group: %d, description: %s, level: %d, vldr_group: %d, operator: %s, values: %s, bitwise: %d",
		m.UUID,
		m.Name,
		m.Group,
		m.Description,
		m.Level,
		m.VldrGroup,
		m.Operator,
		m.Values,
		m.Bitwise,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Rule) IsNode() {}

// 创建时间
func (m *Rule) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Rule) GetUpdateAt() *uint {
	return &m.UpdateAt
}
