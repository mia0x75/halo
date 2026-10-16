package validate

import (
	"fmt"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"
)

// DeleteVldr 删除数据语句相关的审核规则
type DeleteVldr struct {
	vldr

	dd *ast.DeleteStmt
}

// Enabled 当前规则组是否生效
func (v *DeleteVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *DeleteVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if dd, ok := node.(*ast.DeleteStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.dd = dd
			v.Walk(v.dd)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// WithoutWhereNotAllowed 不允许没有WHERE的删除
// RULE: DEL-L2-001
func (v *DeleteVldr) WithoutWhereNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.dd.Where == nil {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// MaxAllowedDeleteRows 单次删除的最大行数
// RULE: DEL-L3-001
func (v *DeleteVldr) MaxAllowedDeleteRows(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DEL-L3-002
func (v *DeleteVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
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

// TargetTableDoesNotDoesNotExist 目标表是否存在
// RULE: DEL-L3-003
func (v *DeleteVldr) TargetTableDoesNotDoesNotExist(s *Subject, r *Rule) {
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

// TargetColumnDoesNotDoesNotExist 条件过滤列是否存在
// RULE: DEL-L3-004
func (v *DeleteVldr) TargetColumnDoesNotDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}
