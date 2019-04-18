package validate

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/models"
)

// ViewCreateVldr 创建视图语句相关的审核规则
type ViewCreateVldr struct {
	vldr

	cv *ast.CreateViewStmt
}

// Call 利用反射方法动态调用审核函数
func (v *ViewCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ViewCreateVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ViewCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ViewCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if cv, ok := node.(*ast.CreateViewStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.cv = cv
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// ViewCreateViewNameQualified 视图名标识符规则
// RULE: CVW-L2-001
func (v *ViewCreateVldr) ViewCreateViewNameQualified(s *models.Statement, r *models.Rule) {
	viewName := v.cv.ViewName.Name.O
	if err := Match(r, viewName, viewName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ViewCreateViewNameLowerCaseRequired 视图名大小写规则
// RULE: CVW-L2-002
func (v *ViewCreateVldr) ViewCreateViewNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	viewName := v.cv.ViewName.Name.O
	if err := Match(r, viewName, viewName); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ViewCreateViewNameMaxLength 视图名长度规则
// RULE: CVW-L2-003
func (v *ViewCreateVldr) ViewCreateViewNameMaxLength(s *models.Statement, r *models.Rule) {
	viewName := v.cv.ViewName.Name.O
	threshold, _ := strconv.Atoi(r.Values)
	if len(viewName) > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, viewName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ViewCreateViewNamePrefixRequired 视图名前缀规则
// RULE: CVW-L2-004
func (v *ViewCreateVldr) ViewCreateViewNamePrefixRequired(s *models.Statement, r *models.Rule) {
	// v.ViewCreateViewNameQualified(stmt, node, rule, result)
}

// ViewCreateTargetDatabaseExists 目标库是否存在
// RULE: CVW-L3-001
func (v *ViewCreateVldr) ViewCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// ViewCreateTargetViewExists 目标视图是否存在
// RULE: CVW-L3-002
func (v *ViewCreateVldr) ViewCreateTargetViewExists(s *models.Statement, r *models.Rule) {
}

// ViewAlterVldr 修改视图语句相关的审核规则
type ViewAlterVldr struct {
	vldr

	// av *ast.AlterViewStmt
}

// Call 利用反射方法动态调用审核函数
func (v *ViewAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ViewAlterVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ViewAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ViewAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	// for _, s := range v.stmts {
	// 	// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
	// 	node := s.StmtNode
	// 	if av, ok := node.(*ast.AlterViewStmt); !ok {
	// 		// 类型断言不成功
	// 		continue
	// 	} else {
	// 		v.av = av
	// 	}
	// 	for _, r := range v.rules {
	// 		if r.Bitwise&1 != 1 {
	// 			continue
	// 		}
	// 		v.Call(r.Func, s, r)
	// 	}
	// }
}

// ViewAlterTargetDatabaseExists 目标库是否存在
// RULE: MVW-L3-001
func (v *ViewAlterVldr) ViewAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// ViewAlterTargetViewExists 目标视图是否存在
// RULE: MVW-L3-002
func (v *ViewAlterVldr) ViewAlterTargetViewExists(s *models.Statement, r *models.Rule) {
}

// ViewDropVldr 删除视图语句相关的审核规则
type ViewDropVldr struct {
	vldr

	// dv *ast.DropViewStmt
}

// Call 利用反射方法动态调用审核函数
func (v *ViewDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ViewDropVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ViewDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ViewDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	// for _, s := range v.stmts {
	// 	// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
	// 	node := s.StmtNode
	// 	if dv, ok := node.(*ast.DropViewStmt); !ok {
	// 		// 类型断言不成功
	// 		continue
	// 	} else {
	// 		v.dv = dv
	// 	}
	// 	for _, r := range v.rules {
	// 		if r.Bitwise&1 != 1 {
	// 			continue
	// 		}
	// 		v.Call(r.Func, s, r)
	// 	}
	// }
}

// ViewDropTargetDatabaseExists 目标库是否存在
// RULE: DVW-L3-001
func (v *ViewDropVldr) ViewDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// ViewDropTargetViewExists 目标视图是否存在
// RULE: DVW-L3-002
func (v *ViewDropVldr) ViewDropTargetViewExists(s *models.Statement, r *models.Rule) {
}
