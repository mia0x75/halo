package validate

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"

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

// Validate 规则组的审核入口
func (v *ViewCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if cv, ok := node.(*ast.CreateViewStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.cv = cv
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// ViewNameQualified 视图名标识符规则
// RULE: CVW-L2-001
func (v *ViewCreateVldr) ViewNameQualified(s *models.Statement, r *models.Rule) {
	viewName := v.cv.ViewName.Name.O
	if err := Match(r, viewName, viewName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ViewNameLowerCaseRequired 视图名大小写规则
// RULE: CVW-L2-002
func (v *ViewCreateVldr) ViewNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	viewName := v.cv.ViewName.Name.O
	if err := Match(r, viewName, viewName); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ViewNameMaxLength 视图名长度规则
// RULE: CVW-L2-003
func (v *ViewCreateVldr) ViewNameMaxLength(s *models.Statement, r *models.Rule) {
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

// ViewNamePrefixRequired 视图名前缀规则
// RULE: CVW-L2-004
func (v *ViewCreateVldr) ViewNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	viewName := strings.TrimSpace(v.cv.ViewName.Name.O)
	if len(viewName) > 0 {
		if err := Match(r, viewName, r.Values); err != nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, viewName, r.Values),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}

}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CVW-L3-001
func (v *ViewCreateVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
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

// TargetViewDoesNotExist 目标视图是否存在
// RULE: CVW-L3-002
func (v *ViewCreateVldr) TargetViewDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// ViewAlterVldr 修改视图语句相关的审核规则
type ViewAlterVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *ViewAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ViewAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ViewAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MVW-L3-001
func (v *ViewAlterVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetViewDoesNotExist 目标视图是否存在
// RULE: MVW-L3-002
func (v *ViewAlterVldr) TargetViewDoesNotExist(s *models.Statement, r *models.Rule) {
}

// ViewDropVldr 删除视图语句相关的审核规则
type ViewDropVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *ViewDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ViewDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ViewDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DVW-L3-001
func (v *ViewDropVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetViewDoesNotExist 目标视图是否存在
// RULE: DVW-L3-002
func (v *ViewDropVldr) TargetViewDoesNotExist(s *models.Statement, r *models.Rule) {
}
