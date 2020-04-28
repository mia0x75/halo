package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Rule 系统审核规则
type Rule struct {
	Name        string `xorm:"'name' notnull char(10) pk"               valid:"required,length(10|10)"            json:"name"        gqlgen:"Name"`        //
	UUID        string `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                 json:"uuid"        gqlgen:"UUID"`        //
	Group       uint8  `xorm:"'group' notnull tinyint <-"               valid:"required,range(0|255)"             json:"group"       gqlgen:"Group"`       //
	Level       uint8  `xorm:"'level' notnull tinyint <-"               valid:"required,matches(^(1|2|3)$)"       json:"level"       gqlgen:"Level"`       //
	VldrGroup   uint16 `xorm:"'vldr_group' notnull smallint <-"         valid:"required,matches(^([0-9]{4})$)"    json:"vldr_group"  gqlgen:"VldrGroup"`   //
	Operator    string `xorm:"'operator' notnull varchar(5) <-"         valid:"required,length(1|10),alphanum"    json:"operator"    gqlgen:"Operator"`    //
	Values      string `xorm:"'values' notnull varchar(150)"            valid:"required,length(1|150)"            json:"values"      gqlgen:"Values"`      //
	Bitwise     uint8  `xorm:"'bitwise' notnull tinyint"                valid:"required,int,matches(^(4|5|6|7)$)" json:"bitwise"     gqlgen:"Bitwise"`     //
	Func        string `xorm:"'func' notnull varchar(75) <-"            valid:"required,length(1|75)"             json:"func"        gqlgen:"Func"`        //
	Message     string `xorm:"'message' notnull varchar(150) <-"        valid:"required,runelength(1|150)"        json:"message"     gqlgen:"Message"`     //
	Description string `xorm:"'description' notnull tinytext <-"        valid:"required,runelength(1|255)"        json:"description" gqlgen:"Description"` //
	Element     string `xorm:"'element' notnull varchar(50) <-"         valid:"required,length(1|50)"             json:"element"     gqlgen:"Element"`     //
	Version     int    `xorm:"'version'"                                valid:"-"                                 json:"version"     gqlgen:"-"`           //
	UpdateAt    uint   `xorm:"'update_at' notnull int"                  valid:"-"                                 json:"update_at"   gqlgen:"UpdateAt"`    //
	CreateAt    uint   `xorm:"'create_at' notnull int"                  valid:"-"                                 json:"create_at"   gqlgen:"CreateAt"`    //
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

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Rule) AfterSet(colName string, _ xorm.Cell) {
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
