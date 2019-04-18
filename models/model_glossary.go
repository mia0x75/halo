package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Glossary 系统字典表
type Glossary struct {
	Group       string `xorm:"'group' notnull varchar(25) pk unique(unique_1) <-" valid:"required,length(1|25)"      json:"group"       gqlgen:"Catalog"`     //
	Key         uint   `xorm:"'key' notnull tinyint pk <-"                        valid:"required,range(0|255)"      json:"key"         gqlgen:"Iota"`        //
	Value       string `xorm:"'value' notnull varchar(50) unique(unique_1) <-"    valid:"required,runelength(1|50)"  json:"value"       gqlgen:"Name"`        //
	UUID        string `xorm:"'uuid' notnull char(36) unique(unique_2)"           valid:"-"                          json:"uuid"        gqlgen:"UUID"`        //
	Description string `xorm:"'description' notnull varchar(100) <-"              valid:"required,runelength(1|100)" json:"description" gqlgen:"Description"` //
	Version     int    `xorm:"'version'"                                          valid:"-"                          json:"version"     gqlgen:"-"`           //
	UpdateAt    uint   `xorm:"'update_at' notnull int"                            valid:"-"                          json:"update_at"   gqlgen:"UpdateAt"`    //
	CreateAt    uint   `xorm:"'create_at' notnull int"                            valid:"-"                          json:"create_at"   gqlgen:"CreateAt"`    //
}

// TableName 结构体到数据库表名称的映射
func (m *Glossary) TableName() string {
	return "mm_glossaries"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Glossary) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Glossary) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Glossary) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Glossary) String() string {
	return fmt.Sprintf("group: %s, key: %d, uuid: %s, value: %s, description: %s",
		m.Group,
		m.Key,
		m.UUID,
		m.Value,
		m.Description,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Glossary) IsNode() {}
