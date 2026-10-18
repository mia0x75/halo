package validate

import (
	"fmt"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"
)

// SelectVldr 查询语句审核
type SelectVldr struct {
	vldr

	sd *ast.SelectStmt
}

// Enabled 当前规则组是否生效
func (v *SelectVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *SelectVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if sd, ok := node.(*ast.SelectStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.sd = sd
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// WithoutWhereNotAllowed 禁止没有WHERE的查询
// RULE: SEL-L2-001
func (v *SelectVldr) WithoutWhereNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.sd.Where == nil {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// WithoutLimitNotAllowed 禁止没有LIMIT的查询
// RULE: SEL-L2-002
func (v *SelectVldr) WithoutLimitNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.sd.Limit == nil {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// UseWildcardNotAllowed 禁止SELECT STAR
// RULE: SEL-L2-003
func (v *SelectVldr) UseWildcardNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, field := range v.sd.Fields.Fields {
		if field.WildCard != nil {
			c := &Clause{
				Description: r.Message,
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
			break
		}
	}
}

// UseExplicitLockNotAllowed 禁止指定锁的类型
// RULE: SEL-L2-004
func (v *SelectVldr) UseExplicitLockNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.sd.LockTp != ast.SelectLockNone {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// TargetDatabaseDoesNotExist 目标数据库必须已存在
// RULE: SEL-L3-001
func (v *SelectVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
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

// TargetTableDoesNotExist 目标表必须已存在
// RULE: SEL-L3-002
func (v *SelectVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
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

// TargetColumnDoesNotDoesNotExist 目标列必须已存在，从抽象语法数解出来所有的关联字段信息，包括FieldList/Where/OrderBy/GroupBy
// RULE: SEL-L3-003
func (v *SelectVldr) TargetColumnDoesNotDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	refs := v.tableRefs(v.sd.From)
	reported := map[string]bool{}
	for _, field := range v.sd.Fields.Fields {
		for _, c := range collectedColumns(field.Expr) {
			name, bad := missingColumn(refs, c)
			if !bad || reported[name] {
				continue
			}
			reported[name] = true
			appendViolation(s, r, name)
		}
	}
}

// ReturnBlobOrTextNotAllowed 是否允许返回BLOB/TEXT列
// RULE: SEL-L3-004
func (v *SelectVldr) ReturnBlobOrTextNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	refs := v.tableRefs(v.sd.From)
	reported := map[string]bool{}
	for _, field := range v.sd.Fields.Fields {
		// Wildcard columns are not expanded; SEL-L2-003 already forbids them.
		ref, ok := field.Expr.(*ast.ColumnNameExpr)
		if !ok || ref.Name == nil {
			continue
		}
		table, ok := lookupColumn(refs, ref.Name)
		if !ok {
			continue
		}
		col := table.GetColumn(ref.Name.Name.O)
		if col == nil || !col.Type.IsLargeObject() {
			continue
		}
		name := ref.Name.Name.O
		if reported[name] {
			continue
		}
		reported[name] = true
		appendViolation(s, r, name)
	}
}
