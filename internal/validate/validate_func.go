package validate

import (
	"sync"
)

// FuncCreateVldr 创建函数语句相关的审核规则
type FuncCreateVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *FuncCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *FuncCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// FuncNameQuilified 函数名标识符规则
// RULE: CFU-L2-001
func (v *FuncCreateVldr) createFuncFuncNameQuilified(s *Subject, r *Rule) {
}

// FuncNameLowerCaseRequired 函数名大小写规则
// RULE: CFU-L2-002
func (v *FuncCreateVldr) createFuncFuncNameLowerCaseRequired(s *Subject, r *Rule) {
}

// FuncNameMaxLength 函数名长度规则
// RULE: CFU-L2-003
func (v *FuncCreateVldr) createFuncFuncNameMaxLength(s *Subject, r *Rule) {
}

// FuncNamePrefixRequired 函数名前缀规则
// RULE: CFU-L2-004
func (v *FuncCreateVldr) createFuncFuncNamePrefixRequired(s *Subject, r *Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CFU-L3-001
func (v *FuncCreateVldr) createFuncTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetFuncDoesNotExist 目标函数是否存在
// RULE: CFU-L3-002
func (v *FuncCreateVldr) createFuncTargetFuncDoesNotExist(s *Subject, r *Rule) {
}

// FuncAlterVldr 修改函数语句相关的审核规则
type FuncAlterVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *FuncAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *FuncAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MFU-L3-001
func (v *FuncAlterVldr) alterFuncTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetFuncDoesNotExist 目标函数是否存在
// RULE: MFU-L3-002
func (v *FuncAlterVldr) alterFuncTargetFuncDoesNotExist(s *Subject, r *Rule) {
}

// FuncDropVldr 删除函数语句相关的审核规则
type FuncDropVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *FuncDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *FuncDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DFU-L3-001
func (v *FuncDropVldr) dropFuncTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetFuncDoesNotExist 目标函数是否存在
// RULE: DFU-L3-002
func (v *FuncDropVldr) dropFuncTargetFuncDoesNotExist(s *Subject, r *Rule) {
}
