package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"

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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *DeleteVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *DeleteVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if dd, ok := node.(*ast.DeleteStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.dd = dd
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// DeleteWithoutWhereEnabled 不允许没有WHERE的删除
// RULE: DEL-L2-001
func (v *DeleteVldr) DeleteWithoutWhereEnabled(s *models.Statement, r *models.Rule) {
	if v.dd.Where == nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DeleteRowsLimit 单次删除的最大行数
// RULE: DEL-L3-001
func (v *DeleteVldr) DeleteRowsLimit(s *models.Statement, r *models.Rule) {
}

// DeleteTargetDatabaseExists 目标库是否存在
// RULE: DEL-L3-002
func (v *DeleteVldr) DeleteTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// DeleteTargetTableExists 目标表是否存在
// RULE: DEL-L3-003
func (v *DeleteVldr) DeleteTargetTableExists(s *models.Statement, r *models.Rule) {
}

// DeleteFilterColumnExists 条件过滤列是否存在
// RULE: DEL-L3-004
func (v *DeleteVldr) DeleteFilterColumnExists(s *models.Statement, r *models.Rule) {
}
