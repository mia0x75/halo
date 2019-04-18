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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *FuncCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *FuncCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// FuncCreateFuncNameQuilified 函数名标识符规则
// RULE: CFU-L2-001
func (v *FuncCreateVldr) FuncCreateFuncNameQuilified(s *models.Statement, r *models.Rule) {
}

// FuncCreateFuncNameLowerCaseRequired 函数名大小写规则
// RULE: CFU-L2-002
func (v *FuncCreateVldr) FuncCreateFuncNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// FuncCreateFuncNameMaxLength 函数名长度规则
// RULE: CFU-L2-003
func (v *FuncCreateVldr) FuncCreateFuncNameMaxLength(s *models.Statement, r *models.Rule) {
}

// FuncCreateFuncNamePrefixRequired 函数名前缀规则
// RULE: CFU-L2-004
func (v *FuncCreateVldr) FuncCreateFuncNamePrefixRequired(s *models.Statement, r *models.Rule) {
}

// FuncCreateTargetDatabaseExists 目标库是否存在
// RULE: CFU-L3-001
func (v *FuncCreateVldr) FuncCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// FuncCreateTargetFuncExists 目标函数是否存在
// RULE: CFU-L3-002
func (v *FuncCreateVldr) FuncCreateTargetFuncExists(s *models.Statement, r *models.Rule) {
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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *FuncAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *FuncAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// FuncAlterTargetDatabaseExists 目标库是否存在
// RULE: MFU-L3-001
func (v *FuncAlterVldr) FuncAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// FuncAlterTargetFuncExists 目标函数是否存在
// RULE: MFU-L3-002
func (v *FuncAlterVldr) FuncAlterTargetFuncExists(s *models.Statement, r *models.Rule) {
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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *FuncDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *FuncDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// FuncDropTargetDatabaseExists 目标库是否存在
// RULE: DFU-L3-001
func (v *FuncDropVldr) FuncDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// FuncDropTargetFuncExists 目标函数是否存在
// RULE: DFU-L3-002
func (v *FuncDropVldr) FuncDropTargetFuncExists(s *models.Statement, r *models.Rule) {
}
