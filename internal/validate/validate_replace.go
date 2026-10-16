package validate

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// ReplaceVldr rule handlers. REPLACE reuses the INSERT AST (*ast.InsertStmt with
// IsReplace set), so the checks mirror InsertVldr's but without the
// "skip when IsReplace" guards: here IsReplace is always true and the rules are
// meant to fire.

// ExplicitColumnRequired 是否要求显式列申明
// RULE: RPL-L2-001
func (v *ReplaceVldr) ExplicitColumnRequired(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if len(v.id.Columns) == 0 {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// UsingSelectNotAllowed 是否允许REPLACE...SELECT
// RULE: RPL-L2-002
func (v *ReplaceVldr) UsingSelectNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.id.Select != nil {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// ColumnsValuesMustMatch 列类型、值是否匹配
// RULE: RPL-L2-005
func (v *ReplaceVldr) ColumnsValuesMustMatch(s *Subject, r *Rule) {
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
// RULE: RPL-L3-001
func (v *ReplaceVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
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
// RULE: RPL-L3-002
func (v *ReplaceVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
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
// RULE: RPL-L3-003
func (v *ReplaceVldr) TargetColumnDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}

// ValueForNotNullColumnRequired 非空列是否有值
// RULE: RPL-L3-004
func (v *ReplaceVldr) ValueForNotNullColumnRequired(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}
