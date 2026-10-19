package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// ProcCreateVldr 创建存储过程语句相关的审核规则
type ProcCreateVldr struct {
	vldr

	cp *ast.CreateProcedureStmt
}

// Enabled 当前规则组是否生效
func (v *ProcCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ProcCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		cp, ok := s.StmtNode.(*ast.CreateProcedureStmt)
		if !ok {
			continue
		}
		v.cp = cp
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// ProcNameQualified 存储过程名标识符规则
// RULE: CSP-L2-001
func (v *ProcCreateVldr) ProcNameQualified(s *Subject, r *Rule) {
	nameRegexRule(s, r, refOf(v.cp.Name).Name, true)
}

// ProcNameLowerCaseRequired 存储过程名大小写规则
// RULE: CSP-L2-002
func (v *ProcCreateVldr) ProcNameLowerCaseRequired(s *Subject, r *Rule) {
	nameRegexRule(s, r, refOf(v.cp.Name).Name, false)
}

// ProcNameMaxLength 存储过程名长度规则
// RULE: CSP-L2-003
func (v *ProcCreateVldr) ProcNameMaxLength(s *Subject, r *Rule) {
	nameMaxLengthRule(s, r, refOf(v.cp.Name).Name)
}

// ProcNamePrefixRequired 存储过程名前缀规则
// RULE: CSP-L2-004
func (v *ProcCreateVldr) ProcNamePrefixRequired(s *Subject, r *Rule) {
	namePrefixRule(s, r, refOf(v.cp.Name).Name)
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CSP-L3-001
func (v *ProcCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.cp.Name).Schema)
}

// TargetProcDoesNotExist 目标存储过程是否存在（创建时要求不存在）
// RULE: CSP-L3-002
func (v *ProcCreateVldr) TargetProcDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.cp.Name)
	v.requireObjectAbsent(s, r, ref.Schema, meta.ObjectProcedure, ref.Name)
}

// ProcAlterVldr 修改存储过程语句相关的审核规则
type ProcAlterVldr struct {
	vldr

	ap *ast.AlterProcedureStmt
}

// Enabled 当前规则组是否生效
func (v *ProcAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ProcAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		ap, ok := s.StmtNode.(*ast.AlterProcedureStmt)
		if !ok {
			continue
		}
		v.ap = ap
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MSP-L3-001
func (v *ProcAlterVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.ap.Name).Schema)
}

// TargetProcDoesNotExist 目标存储过程是否存在（修改时要求已存在）
// RULE: MSP-L3-002
func (v *ProcAlterVldr) TargetProcDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.ap.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectProcedure, ref.Name)
}

// ProcDropVldr 删除存储过程语句相关的审核规则
type ProcDropVldr struct {
	vldr

	dp *ast.DropProcedureStmt
}

// Enabled 当前规则组是否生效
func (v *ProcDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ProcDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		dp, ok := s.StmtNode.(*ast.DropProcedureStmt)
		if !ok {
			continue
		}
		v.dp = dp
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DSP-L3-001
func (v *ProcDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.dp.Name).Schema)
}

// TargetProcDoesNotExist 目标存储过程是否存在（删除时要求已存在）
// RULE: DSP-L3-002
func (v *ProcDropVldr) TargetProcDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.dp.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectProcedure, ref.Name)
}
