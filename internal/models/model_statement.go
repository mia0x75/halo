package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mia0x75/parser/ast"
)

// Statement 工单分解出来的单个SQL语句
type Statement struct {
	TicketID     uint         `db:"ticket_id"               valid:"required,int,range(0|4294967295)"  json:"ticket_id"     gqlgen:"-"`                 //
	Sequence     uint16       `db:"sequence"           valid:"required,int,range(0|65535)"       json:"sequence"      gqlgen:"Sequence"`               //
	UUID         string       `db:"uuid" valid:"-"                                 json:"uuid"          gqlgen:"UUID"`                                 //
	Content      string       `db:"content"                   valid:"required,length(1|65535),ascii"    json:"content"       gqlgen:"Content"`         //
	Type         uint8        `db:"type"                   valid:"required,matches(^([1-9]?[0-9])$)" json:"type"          gqlgen:"Type"`               // 类型 0-99
	Status       uint8        `db:"status"                 valid:"required,matches(^([1-9]?[0-9])$)" json:"status"        gqlgen:"Status"`             // 状态 0-99
	Report       string       `db:"report"                    valid:"required,length(1|65535)"          json:"report"        gqlgen:"Report"`          //
	Plan         string       `db:"plan"                      valid:"required,length(1|65535)"          json:"plan"          gqlgen:"Plan"`            //
	Results      string       `db:"results"                           valid:"length(1|65535)"                   json:"results"       gqlgen:"Results"` //
	RowsAffected uint         `db:"rows_affected"              valid:"required,int,range(0|4294967295)"  json:"rows_affected" gqlgen:"RowsAffected"`   //
	Version      int          `db:"version"                                valid:"-"                                 json:"version"       gqlgen:"-"`  //
	UpdateAt     uint         `db:"update_at"                  valid:"-"                                 json:"update_at"     gqlgen:"UpdateAt"`       //
	CreateAt     uint         `db:"create_at"                  valid:"-"                                 json:"create_at"     gqlgen:"CreateAt"`       //
	StmtNode     ast.StmtNode `valid:"-"                                 json:"-"             gqlgen:"-"`                                              //
	Violations   *Violations  `valid:"-"                                 json:"-"             gqlgen:"-"`                                              //
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
	Level uint8
	// Description 是命中提示
	Description string
	// Rule 是命中规则的名称，由规则分发处回填
	Rule string `json:",omitempty"`
}

// Len 返回当前记录的命中条数
func (v *Violations) Len() int {
	v.Lock()
	defer v.Unlock()
	return len(v.clauses)
}

// TagRule 给 offset 之后新增的命中补上规则名。规则分发处知道当前规则的名称，
// 而具体的规则处理器只负责产出描述，两者在此汇合。
func (v *Violations) TagRule(rule string, offset int) {
	v.Lock()
	defer v.Unlock()
	for i := offset; i < len(v.clauses); i++ {
		if v.clauses[i].Rule == "" {
			v.clauses[i].Rule = rule
		}
	}
}
