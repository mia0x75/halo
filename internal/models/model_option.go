package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Option 系统选项
type Option struct {
	Name        string `db:"name"         valid:"required,length(1|50),alphanum"  json:"group"       gqlgen:"Group"`                                    //
	UUID        string `db:"uuid" valid:"-"                               json:"uuid"        gqlgen:"UUID"`                                             //
	Parameter   string `db:"parameter"                 valid:"required,length(1|255),alphanum" json:"parameter"       gqlgen:"Parameter"`               //
	Writable    uint8  `db:"writable"            valid:"required,range(0|1)"             json:"writable"    gqlgen:"Writable"`                          //
	Description string `db:"description"     valid:"required,runelength(1|75)"       json:"description" gqlgen:"Description"`                           //
	Element     string `db:"element"         valid:"required,runelength(1|15)"       json:"element"     gqlgen:"Element"`                               //
	GroupName   string `db:"group_name"                                valid:"-"                                 json:"groupName"   gqlgen:"GroupName"` //
	SortOrder   int    `db:"sort_order"                                valid:"-"                                 json:"sortOrder"   gqlgen:"SortOrder"` //
	DependsOn   string `db:"depends_on"                                valid:"-"                                 json:"dependsOn"   gqlgen:"DependsOn"` //
	MinValue    *int   `db:"min_value"                                 valid:"-"                                 json:"minValue"    gqlgen:"MinValue"`  //
	MaxValue    *int   `db:"max_value"                                 valid:"-"                                 json:"maxValue"    gqlgen:"MaxValue"`  //
	Version     int    `db:"version"                                valid:"-"                               json:"version"     gqlgen:"-"`              //
	UpdateAt    uint   `db:"update_at"                  valid:"-"                               json:"update_at"   gqlgen:"UpdateAt"`                   //
	CreateAt    uint   `db:"create_at"                  valid:"-"                               json:"create_at"   gqlgen:"CreateAt"`                   //
}

// Secret reports whether the value must never be returned in clear text.
func (m *Option) Secret() bool {
	return m != nil && m.Element == "password"
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

// String 结构体输出到字符串的默认方式
func (m *Option) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, parameter: %s, description: %s, writable: %d, element: %s",
		m.UUID,
		m.Name,
		m.Parameter,
		m.Description,
		m.Writable,
		m.Element,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Option) IsNode() {}

// 创建时间
func (m *Option) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Option) GetUpdateAt() *uint {
	return &m.UpdateAt
}
