package validate

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"
)

// InsertVldr 插入数据语句相关的审核规则
type InsertVldr struct {
	vldr

	id *ast.InsertStmt
}

// Enabled 当前规则组是否生效
func (v *InsertVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *InsertVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if id, ok := node.(*ast.InsertStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.id = id
			v.Walk(v.id)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// ExplicitColumnRequired 是否要求显式列申明
// RULE: INS-L2-001
func (v *InsertVldr) ExplicitColumnRequired(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.id.IsReplace {
		return
	}

	if len(v.id.Columns) == 0 {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// UsingSelectNotAllowed 是否允许INSERT...SELECT
// RULE: INS-L2-002
func (v *InsertVldr) UsingSelectNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.id.IsReplace {
		return
	}

	if v.id.Select != nil {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// RowsLimit 单语句允许操作的最大行数
// RULE: INS-L2-004
func (v *InsertVldr) RowsLimit(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	thredhold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	if len(v.id.Lists) > thredhold {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, thredhold),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// ColumnsValuesMustMatch 列类型、值是否匹配
// RULE: INS-L2-005
func (v *InsertVldr) ColumnsValuesMustMatch(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	columnNum := len(v.id.Columns)
	for _, list := range v.id.Lists {
		if len(list) != columnNum {
			c := &Clause{
				Description: r.Message,
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: INS-L3-001
func (v *InsertVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, ti.Database),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: INS-L3-002
func (v *InsertVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if ti.Table == nil {
			continue
		}
		if v.TableInfo(ti.Database, ti.Table.Name) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", ti.Database, ti.Table.Name)),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetColumnDoesNotExist 目标列是否存在
// RULE: INS-L3-003
func (v *InsertVldr) TargetColumnDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}

// ValueForNotNullColumnRequired 非空列是否有值
// RULE: INS-L3-004
func (v *InsertVldr) ValueForNotNullColumnRequired(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}

// ReplaceVldr 插入数据语句相关的审核规则
type ReplaceVldr struct {
	vldr

	id *ast.InsertStmt
}

// Enabled 当前规则组是否生效
func (v *ReplaceVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ReplaceVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if id, ok := node.(*ast.InsertStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.id = id
			v.Walk(v.id)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}
