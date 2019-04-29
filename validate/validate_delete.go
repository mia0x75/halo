package validate

import (
	"fmt"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/models"
)

// DeleteVldr 删除数据语句相关的审核规则
type DeleteVldr struct {
	vldr

	dd *ast.DeleteStmt
}

// Call 利用反射方法动态调用审核函数
func (v *DeleteVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *DeleteVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *DeleteVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
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
			v.Call(r.Func, s, r)
		}
	}
}

// WithoutWhereNotAllowed 不允许没有WHERE的删除
// RULE: DEL-L2-001
func (v *DeleteVldr) WithoutWhereNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.dd.Where == nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// RowsLimit 单次删除的最大行数
// RULE: DEL-L3-001
func (v *DeleteVldr) RowsLimit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DEL-L3-002
func (v *DeleteVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, ti.Database),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// TargetTableDoesNotDoesNotExist 目标表是否存在
// RULE: DEL-L3-003
func (v *DeleteVldr) TargetTableDoesNotDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if ti.Table == nil {
			continue
		}
		if v.TableInfo(ti.Database, ti.Table.Name) == nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", ti.Database, ti.Table.Name)),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// TargetColumnDoesNotDoesNotExist 条件过滤列是否存在
// RULE: DEL-L3-004
func (v *DeleteVldr) TargetColumnDoesNotDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO:
}
