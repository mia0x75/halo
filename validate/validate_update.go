package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/models"
)

// UpdateVldr 数据更新语句相关的审核规则
type UpdateVldr struct {
	vldr

	ud *ast.UpdateStmt
}

// Call 利用反射方法动态调用审核函数
func (v *UpdateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *UpdateVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *UpdateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *UpdateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ud, ok := node.(*ast.UpdateStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ud = ud
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// UpdateWithoutWhereEnabled 是否允许没有WHERE的更新
// RULE: UPD-L2-001
func (v *UpdateVldr) UpdateWithoutWhereEnabled(s *models.Statement, r *models.Rule) {
	if v.ud.Where == nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// UpdateTargetDatabaseExists 目标库是否存在
// RULE: UPD-L3-001
func (v *UpdateVldr) UpdateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// UpdateTargetTableExists 目标表是否存在
// RULE: UPD-L3-002
func (v *UpdateVldr) UpdateTargetTableExists(s *models.Statement, r *models.Rule) {
}

// UpdateTargetColumnExists 目标列是否存在
// RULE: UPD-L3-003
func (v *UpdateVldr) UpdateTargetColumnExists(s *models.Statement, r *models.Rule) {
}

// UpdateFilterColumnExists 条件过滤列是否存在
// RULE: UPD-L3-004
func (v *UpdateVldr) UpdateFilterColumnExists(s *models.Statement, r *models.Rule) {
}

// UpdateRowsLimit 允许单次更新的最大行数
// RULE: UPD-L3-005
func (v *UpdateVldr) UpdateRowsLimit(s *models.Statement, r *models.Rule) {
}
