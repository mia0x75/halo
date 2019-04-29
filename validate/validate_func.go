package validate

import (
	"sync"

	"github.com/mia0x75/halo/models"
)

// FuncCreateVldr 创建函数语句相关的审核规则
type FuncCreateVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *FuncCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
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
func (v *FuncCreateVldr) FuncNameQuilified(s *models.Statement, r *models.Rule) {
}

// FuncNameLowerCaseRequired 函数名大小写规则
// RULE: CFU-L2-002
func (v *FuncCreateVldr) FuncNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// FuncNameMaxLength 函数名长度规则
// RULE: CFU-L2-003
func (v *FuncCreateVldr) FuncNameMaxLength(s *models.Statement, r *models.Rule) {
}

// FuncNamePrefixRequired 函数名前缀规则
// RULE: CFU-L2-004
func (v *FuncCreateVldr) FuncNamePrefixRequired(s *models.Statement, r *models.Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CFU-L3-001
func (v *FuncCreateVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetFuncDoesNotExist 目标函数是否存在
// RULE: CFU-L3-002
func (v *FuncCreateVldr) TargetFuncDoesNotExist(s *models.Statement, r *models.Rule) {
}

// FuncAlterVldr 修改函数语句相关的审核规则
type FuncAlterVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *FuncAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
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
func (v *FuncAlterVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetFuncDoesNotExist 目标函数是否存在
// RULE: MFU-L3-002
func (v *FuncAlterVldr) TargetFuncDoesNotExist(s *models.Statement, r *models.Rule) {
}

// FuncDropVldr 删除函数语句相关的审核规则
type FuncDropVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *FuncDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
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
func (v *FuncDropVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetFuncDoesNotExist 目标函数是否存在
// RULE: DFU-L3-002
func (v *FuncDropVldr) TargetFuncDoesNotExist(s *models.Statement, r *models.Rule) {
}
