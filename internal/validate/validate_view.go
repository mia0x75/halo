package validate

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// ViewCreateVldr 创建视图语句相关的审核规则
type ViewCreateVldr struct {
	vldr

	cv *ast.CreateViewStmt
}

// Enabled 当前规则组是否生效
func (v *ViewCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ViewCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if cv, ok := node.(*ast.CreateViewStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.cv = cv
			v.Walk(cv)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// ViewNameQualified 视图名标识符规则
// RULE: CVW-L2-001
func (v *ViewCreateVldr) ViewNameQualified(s *Subject, r *Rule) {
	viewName := v.cv.ViewName.Name.O
	if err := Match(r, viewName, viewName, r.Parameters); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// ViewNameLowerCaseRequired 视图名大小写规则
// RULE: CVW-L2-002
func (v *ViewCreateVldr) ViewNameLowerCaseRequired(s *Subject, r *Rule) {
	viewName := v.cv.ViewName.Name.O
	if err := Match(r, viewName, viewName); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// ViewNameMaxLength 视图名长度规则
// RULE: CVW-L2-003
func (v *ViewCreateVldr) ViewNameMaxLength(s *Subject, r *Rule) {
	viewName := v.cv.ViewName.Name.O
	threshold, _ := strconv.Atoi(r.Parameters)
	if len(viewName) > threshold {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, viewName, threshold),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// ViewNamePrefixRequired 视图名前缀规则
// RULE: CVW-L2-004
func (v *ViewCreateVldr) ViewNamePrefixRequired(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	viewName := strings.TrimSpace(v.cv.ViewName.Name.O)
	if len(viewName) > 0 {
		if err := Match(r, viewName, r.Parameters); err != nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, viewName, r.Parameters),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}

}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CVW-L3-001
func (v *ViewCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
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

// TargetViewDoesNotExist 目标视图是否存在（创建时要求不存在）
// RULE: CVW-L3-002
func (v *ViewCreateVldr) TargetViewDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ref := refOf(v.cv.ViewName)
	v.requireObjectAbsent(s, r, ref.Schema, meta.ObjectView, ref.Name)
}

// ViewAlterVldr 修改视图语句相关的审核规则
type ViewAlterVldr struct {
	vldr

	av *ast.AlterViewStmt
}

// Enabled 当前规则组是否生效
func (v *ViewAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ViewAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		av, ok := s.StmtNode.(*ast.AlterViewStmt)
		if !ok {
			continue
		}
		v.av = av
		v.Walk(av)
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MVW-L3-001
func (v *ViewAlterVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &Clause{Description: fmt.Sprintf(r.Message, ti.Database), Level: uint8(r.Level)}
			s.Violations.Append(c)
		}
	}
}

// TargetViewDoesNotExist 目标视图是否存在（修改时要求已存在）
// RULE: MVW-L3-002
func (v *ViewAlterVldr) TargetViewDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ref := refOf(v.av.ViewName)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectView, ref.Name)
}

// ViewDropVldr 删除视图语句相关的审核规则
type ViewDropVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *ViewDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ViewDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		node := s.StmtNode
		dt, ok := node.(*ast.DropTableStmt)
		if !ok || !dt.IsView {
			continue
		}
		v.Walk(dt)
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DVW-L3-001
func (v *ViewDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &Clause{Description: fmt.Sprintf(r.Message, ti.Database), Level: uint8(r.Level)}
			s.Violations.Append(c)
		}
	}
}

// TargetViewDoesNotExist 目标视图是否存在（删除时要求已存在）
// RULE: DVW-L3-002
func (v *ViewDropVldr) TargetViewDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if ti.Table == nil {
			continue
		}
		v.requireObjectPresent(s, r, ti.Database, meta.ObjectView, ti.Table.Name)
	}
}
