package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// FuncCreateVldr 创建函数语句相关的审核规则
type FuncCreateVldr struct {
	vldr

	cf *ast.CreateFunctionStmt
}

// Enabled 当前规则组是否生效
func (v *FuncCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *FuncCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		cf, ok := s.StmtNode.(*ast.CreateFunctionStmt)
		if !ok {
			continue
		}
		v.cf = cf
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// FuncNameQuilified 函数名标识符规则
// RULE: CFU-L2-001
func (v *FuncCreateVldr) FuncNameQuilified(s *Subject, r *Rule) {
	ref := refOf(v.cf.Name)
	nameRegexRule(s, r, ref.Name, true)
}

// FuncNameLowerCaseRequired 函数名大小写规则
// RULE: CFU-L2-002
func (v *FuncCreateVldr) FuncNameLowerCaseRequired(s *Subject, r *Rule) {
	ref := refOf(v.cf.Name)
	nameRegexRule(s, r, ref.Name, false)
}

// FuncNameMaxLength 函数名长度规则
// RULE: CFU-L2-003
func (v *FuncCreateVldr) FuncNameMaxLength(s *Subject, r *Rule) {
	ref := refOf(v.cf.Name)
	nameMaxLengthRule(s, r, ref.Name)
}

// FuncNamePrefixRequired 函数名前缀规则
// RULE: CFU-L2-004
func (v *FuncCreateVldr) FuncNamePrefixRequired(s *Subject, r *Rule) {
	ref := refOf(v.cf.Name)
	namePrefixRule(s, r, ref.Name)
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CFU-L3-001
func (v *FuncCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.cf.Name)
	v.requireDatabase(s, r, ref.Schema)
}

// TargetFuncDoesNotExist 目标函数是否存在（创建时要求不存在）
// RULE: CFU-L3-002
func (v *FuncCreateVldr) TargetFuncDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.cf.Name)
	v.requireObjectAbsent(s, r, ref.Schema, meta.ObjectFunction, ref.Name)
}

// FuncAlterVldr 修改函数语句相关的审核规则
type FuncAlterVldr struct {
	vldr

	af *ast.AlterFunctionStmt
}

// Enabled 当前规则组是否生效
func (v *FuncAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *FuncAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		af, ok := s.StmtNode.(*ast.AlterFunctionStmt)
		if !ok {
			continue
		}
		v.af = af
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MFU-L3-001
func (v *FuncAlterVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.af.Name)
	v.requireDatabase(s, r, ref.Schema)
}

// TargetFuncDoesNotExist 目标函数是否存在（修改时要求已存在）
// RULE: MFU-L3-002
func (v *FuncAlterVldr) TargetFuncDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.af.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectFunction, ref.Name)
}

// FuncDropVldr 删除函数语句相关的审核规则
type FuncDropVldr struct {
	vldr

	df *ast.DropFunctionStmt
}

// Enabled 当前规则组是否生效
func (v *FuncDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *FuncDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		df, ok := s.StmtNode.(*ast.DropFunctionStmt)
		if !ok {
			continue
		}
		v.df = df
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DFU-L3-001
func (v *FuncDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.df.Name)
	v.requireDatabase(s, r, ref.Schema)
}

// TargetFuncDoesNotExist 目标函数是否存在（删除时要求已存在）
// RULE: DFU-L3-002
func (v *FuncDropVldr) TargetFuncDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.df.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectFunction, ref.Name)
}
