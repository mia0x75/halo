package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mia0x75/parser/ast"
	"xorm.io/xorm"
)

// Statement 工单分解出来的单个SQL语句
type Statement struct {
	TicketID     uint         `xorm:"'ticket_id' notnull int pk"               valid:"required,int,range(0|4294967295)"  json:"ticket_id"     gqlgen:"-"`            //
	Sequence     uint16       `xorm:"'sequence' notnull smallint pk"           valid:"required,int,range(0|65535)"       json:"sequence"      gqlgen:"Sequence"`     //
	UUID         string       `xorm:"'uuid' notnull char(36) unique(unique_1)" valid:"-"                                 json:"uuid"          gqlgen:"UUID"`         //
	Content      string       `xorm:"'content' notnull text"                   valid:"required,length(1|65535),ascii"    json:"content"       gqlgen:"Content"`      //
	Type         uint8        `xorm:"'type' notnull tinyint"                   valid:"required,matches(^([1-9]?[0-9])$)" json:"type"          gqlgen:"Type"`         // 类型 0-99
	Status       uint8        `xorm:"'status' notnull tinyint"                 valid:"required,matches(^([1-9]?[0-9])$)" json:"status"        gqlgen:"Status"`       // 状态 0-99
	Report       string       `xorm:"'report' notnull json"                    valid:"required,length(1|65535)"          json:"report"        gqlgen:"Report"`       //
	Plan         string       `xorm:"'plan' notnull json"                      valid:"required,length(1|65535)"          json:"plan"          gqlgen:"Plan"`         //
	Results      string       `xorm:"'results' text"                           valid:"length(1|65535)"                   json:"results"       gqlgen:"Results"`      //
	RowsAffected uint         `xorm:"'rows_affected' notnull int"              valid:"required,int,range(0|4294967295)"  json:"rows_affected" gqlgen:"RowsAffected"` //
	Version      int          `xorm:"'version'"                                valid:"-"                                 json:"version"       gqlgen:"-"`            //
	UpdateAt     uint         `xorm:"'update_at' notnull int"                  valid:"-"                                 json:"update_at"     gqlgen:"UpdateAt"`     //
	CreateAt     uint         `xorm:"'create_at' notnull int"                  valid:"-"                                 json:"create_at"     gqlgen:"CreateAt"`     //
	StmtNode     ast.StmtNode `xorm:"-"                                        valid:"-"                                 json:"-"             gqlgen:"-"`            //
	Violations   *Violations  `xorm:"-"                                        valid:"-"                                 json:"-"             gqlgen:"-"`            //
}

// TableName 结构体到数据库表名称的映射
func (m *Statement) TableName() string {
	return "mm_statements"
}

// BeforeInsert ORM在执行数据插入前会调用该方法
func (m *Statement) BeforeInsert() {
	m.UUID = uuid.New().String()
	m.Content = strings.TrimSpace(m.Content)
	m.CreateAt = uint(time.Now().Unix())
}

// BeforeUpdate ORM在执行数据更新前会调用该方法
func (m *Statement) BeforeUpdate() {
	m.Content = strings.TrimSpace(m.Content)
	m.UpdateAt = uint(time.Now().Unix())
}

// AfterSet ORM在执行数据更新后会调用该方法
func (m *Statement) AfterSet(colName string, _ xorm.Cell) {
}

// String 结构体输出到字符串的默认方式
func (m *Statement) String() string {
	return fmt.Sprintf("uuid: %s, content: %s, type: %d, status: %d, report: %s",
		m.UUID,
		m.Content,
		m.Type,
		m.Status,
		m.Report,
	)
}

// IsNode GraphQL的基类需要实现的接口，暂时不动
func (Statement) IsNode() {}

// 创建时间
func (m *Statement) GetCreateAt() uint {
	return m.CreateAt
}

// 最后一次修改时间
func (m *Statement) GetUpdateAt() *uint {
	return &m.UpdateAt
}

// Violations 单独一条语句不通过的所有的信息描述
type Violations struct {
	sync.Mutex
	clauses []*Clause
}

// Add 增加一个描述
func (v *Violations) Add(level uint8, description string) {
	v.Lock()
	defer v.Unlock()
	v.clauses = append(v.clauses, &Clause{
		Level:       level,
		Description: description,
	})
}

// Append 增加一个描述
func (v *Violations) Append(clause *Clause) {
	v.Lock()
	defer v.Unlock()
	v.clauses = append(v.clauses, clause)
}

// Marshal 把结构体的内容最终生成报告
func (v *Violations) Marshal() string {
	v.Lock()
	defer v.Unlock()
	if v.clauses != nil {
		if len(v.clauses) > 0 {
			if bs, err := json.Marshal(v.clauses); err == nil {
				return string(bs)
			}
		}
	}
	return ""
}

// Clauses 对于Violations上私有属性，通过方法返回
func (v *Violations) Clauses() []*Clause {
	v.Lock()
	defer v.Unlock()
	return v.clauses
}

// Clause 语句审核不通过的描述结构体
type Clause struct {
	Level       uint8
	Description string
}
