package models

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Edge 根据图论，Edge是连接不同对象的线
type Edge struct {
	EdgeID       uint   `xorm:"'edge_id' notnull int pk autoincr"            valid:"-"                                json:"edge_id"       gqlgen:"-"`            //
	UUID         string `xorm:"'uuid' notnull char(36) unique(unique_1)"     valid:"-"                                json:"uuid"          gqlgen:"UUID"`         //
	Type         uint   `xorm:"'type' notnull int unique(unique_2)"          valid:"required,int,range(0|4294967295)" json:"type"          gqlgen:"Type"`         //
	AncestorID   uint   `xorm:"'ancestor_id' notnull int unique(unique_2)"   valid:"required,int,range(0|4294967295)" json:"ancestor_id"   gqlgen:"AncestorID"`   //
	DescendantID uint   `xorm:"'descendant_id' notnull int unique(unique_2)" valid:"required,int,range(0|4294967295)" json:"descendant_id" gqlgen:"DescendantID"` //
	Version      int    `xorm:"'version'"                                    valid:"-"                                json:"version"       gqlgen:"-"`            //
	UpdateAt     uint   `xorm:"'update_at' notnull int"                      valid:"-"                                json:"update_at"     gqlgen:"UpdateAt"`     //
	CreateAt     uint   `xorm:"'create_at' notnull int"                      valid:"-"                                json:"create_at"     gqlgen:"CreateAt"`     //
}

// TableName 结构体到数据库表名称的映射
func (m *Edge) TableName() string {
	return "mm_edges"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Edge) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Edge) BeforeUpdate() {
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Edge) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Edge) String() string {
	return fmt.Sprintf("uuid: %s, type: %d, ancestor_id: %d, descendant_id: %d",
		m.UUID,
		m.Type,
		m.AncestorID,
		m.DescendantID,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Edge) IsNode() {}
