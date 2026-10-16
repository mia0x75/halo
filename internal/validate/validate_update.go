package validate

import (
	"fmt"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"
)

// UpdateVldr 数据更新语句相关的审核规则
type UpdateVldr struct {
	vldr

	ud *ast.UpdateStmt
}

// Enabled 当前规则组是否生效
func (v *UpdateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *UpdateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ud, ok := node.(*ast.UpdateStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ud = ud
			v.Walk(v.ud)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// WithoutWhereNotAllowed 是否允许没有WHERE的更新
// RULE: UPD-L2-001
func (v *UpdateVldr) WithoutWhereNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.ud.Where == nil {
		c := &Clause{
			Description: r.Message,
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: UPD-L3-001
func (v *UpdateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
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
// RULE: UPD-L3-002
func (v *UpdateVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
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
// RULE: UPD-L3-003
func (v *UpdateVldr) TargetColumnDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}

// MaxAllowedUpdateRows 允许单次更新的最大行数
// RULE: UPD-L3-004
func (v *UpdateVldr) MaxAllowedUpdateRows(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}
