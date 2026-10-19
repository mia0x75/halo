package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// EventCreateVldr 创建事件语句相关的审核规则
type EventCreateVldr struct {
	vldr

	ce *ast.CreateEventStmt
}

// Enabled 当前规则组是否生效
func (v *EventCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *EventCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		ce, ok := s.StmtNode.(*ast.CreateEventStmt)
		if !ok {
			continue
		}
		v.ce = ce
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// EventNameQualified 事件名标识符规则
// RULE: CEV-L2-001
func (v *EventCreateVldr) EventNameQualified(s *Subject, r *Rule) {
	nameRegexRule(s, r, refOf(v.ce.Name).Name, true)
}

// EventNameLowerCaseRequired 事件名大小写规则
// RULE: CEV-L2-002
func (v *EventCreateVldr) EventNameLowerCaseRequired(s *Subject, r *Rule) {
	nameRegexRule(s, r, refOf(v.ce.Name).Name, false)
}

// EventNameMaxLength 事件名长度规则
// RULE: CEV-L2-003
func (v *EventCreateVldr) EventNameMaxLength(s *Subject, r *Rule) {
	nameMaxLengthRule(s, r, refOf(v.ce.Name).Name)
}

// EventNamePrefixRequired 事件名前缀规则
// RULE: CEV-L2-004
func (v *EventCreateVldr) EventNamePrefixRequired(s *Subject, r *Rule) {
	namePrefixRule(s, r, refOf(v.ce.Name).Name)
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CEV-L3-001
func (v *EventCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.ce.Name).Schema)
}

// TargetEventDoesNotExist 目标事件是否存在（创建时要求不存在）
// RULE: CEV-L3-002
func (v *EventCreateVldr) TargetEventDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.ce.Name)
	v.requireObjectAbsent(s, r, ref.Schema, meta.ObjectEvent, ref.Name)
}

// EventAlterVldr 修改事件语句相关的审核规则
type EventAlterVldr struct {
	vldr

	ae *ast.AlterEventStmt
}

// Enabled 当前规则组是否生效
func (v *EventAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *EventAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		ae, ok := s.StmtNode.(*ast.AlterEventStmt)
		if !ok {
			continue
		}
		v.ae = ae
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MEV-L3-001
func (v *EventAlterVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.ae.Name).Schema)
}

// TargetEventDoesNotExist 目标事件是否存在（修改时要求已存在）
// RULE: MEV-L3-002
func (v *EventAlterVldr) TargetEventDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.ae.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectEvent, ref.Name)
}

// EventDropVldr 删除事件语句相关的审核规则
type EventDropVldr struct {
	vldr

	de *ast.DropEventStmt
}

// Enabled 当前规则组是否生效
func (v *EventDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *EventDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		de, ok := s.StmtNode.(*ast.DropEventStmt)
		if !ok {
			continue
		}
		v.de = de
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DEV-L3-001
func (v *EventDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.de.Name).Schema)
}

// TargetEventDoesNotExist 目标事件是否存在（删除时要求已存在）
// RULE: DEV-L3-002
func (v *EventDropVldr) TargetEventDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.de.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectEvent, ref.Name)
}
