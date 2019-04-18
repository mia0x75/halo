package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Option 系统选项
type Option struct {
	Name        string `xorm:"'name' notnull varchar(50) pk <-"         valid:"required,length(1|50),alphanum"  json:"group"       gqlgen:"Group"`       //
	UUID        string `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                               json:"uuid"        gqlgen:"UUID"`        //
	Value       string `xorm:"'value' notnull tinytext"                 valid:"required,length(1|255),alphanum" json:"value"       gqlgen:"Value"`       //
	Writable    uint8  `xorm:"'writable' notnull tinyint <-"            valid:"required,range(0|1)"             json:"writable"    gqlgen:"Writable"`    //
	Description string `xorm:"'description' notnull varchar(75) <-"     valid:"required,runelength(1|75)"       json:"description" gqlgen:"Description"` //
	Element     string `xorm:"'element' notnull varchar(15) <-"         valid:"required,runelength(1|15)"       json:"element"     gqlgen:"Element"`     //
	Version     int    `xorm:"'version'"                                valid:"-"                               json:"version"     gqlgen:"-"`           //
	UpdateAt    uint   `xorm:"'update_at' notnull int"                  valid:"-"                               json:"update_at"   gqlgen:"UpdateAt"`    //
	CreateAt    uint   `xorm:"'create_at' notnull int"                  valid:"-"                               json:"create_at"   gqlgen:"CreateAt"`    //
}

// TableName 结构体到数据库表名称的映射
func (m *Option) TableName() string {
	return "mm_options"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Option) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Option) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Option) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Option) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, value: %s, description: %s, writable: %d, element: %s",
		m.UUID,
		m.Name,
		m.Value,
		m.Description,
		m.Writable,
		m.Element,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Option) IsNode() {}
