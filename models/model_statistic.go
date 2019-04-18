package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Statistic 统计模型
type Statistic struct {
	Group    string  `xorm:"'group' notnull char(25) pk"              valid:"-"                              json:"group"        gqlgen:"Group"`    //
	Key      string  `xorm:"'key' notnull varchar(50) pk"             valid:"required,length(1|50),alphanum" json:"name"         gqlgen:"Key"`      //
	UUID     string  `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                              json:"uuid"         gqlgen:"UUID"`     //
	Value    float64 `xorm:"'value' notnull decimal(18,4)"            valid:"required"                       json:"value"        gqlgen:"Value"`    //
	Version  int     `xorm:"'version'"                                valid:"-"                              json:"version"      gqlgen:"-"`        //
	UpdateAt uint    `xorm:"'update_at' notnull int"                  valid:"-"                              json:"update_at"    gqlgen:"UpdateAt"` //
	CreateAt uint    `xorm:"'create_at' notnull int"                  valid:"-"                              json:"create_at"    gqlgen:"CreateAt"` //
}

// TableName 结构体到数据库表名称的映射
func (m *Statistic) TableName() string {
	return "mm_statistics"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Statistic) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Statistic) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Statistic) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Statistic) String() string {
	return fmt.Sprintf("uuid: %s, group: %s, name: %s, value: %f",
		m.UUID,
		m.Group,
		m.Key,
		m.Value,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Statistic) IsNode() {}
