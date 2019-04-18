package validate

import (
	"sync"

	"github.com/mia0x75/halo/models"
)

// ProcCreateVldr 创建存储过程语句相关的审核规则
type ProcCreateVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *ProcCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ProcCreateVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ProcCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ProcCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// ProcCreateProcNameQualified 存储过程名标识符规则
// RULE: CSP-L2-001
func (v *ProcCreateVldr) ProcCreateProcNameQualified(s *models.Statement, r *models.Rule) {
}

// ProcCreateProcNameLowerCaseRequired 存储过程名大小写规则
// RULE: CSP-L2-002
func (v *ProcCreateVldr) ProcCreateProcNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// ProcCreateProcNameMaxLength 存储过程名长度规则
// RULE: CSP-L2-003
func (v *ProcCreateVldr) ProcCreateProcNameMaxLength(s *models.Statement, r *models.Rule) {
}

// ProcCreateProcNamePrefixRequired 存储过程名前缀规则
// RULE: CSP-L2-004
func (v *ProcCreateVldr) ProcCreateProcNamePrefixRequired(s *models.Statement, r *models.Rule) {
}

// ProcCreateTargetDatabaseExists 目标库是否存在
// RULE: CSP-L3-001
func (v *ProcCreateVldr) ProcCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// ProcCreateTargetProcExists 目标存储过程是否存在
// RULE: CSP-L3-002
func (v *ProcCreateVldr) ProcCreateTargetProcExists(s *models.Statement, r *models.Rule) {
}

// ProcAlterVldr 修改存储过程语句相关的审核规则
type ProcAlterVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *ProcAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ProcAlterVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ProcAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ProcAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// ProcAlterTargetDatabaseExists 目标库是否存在
// RULE: MSP-L3-001
func (v *ProcAlterVldr) ProcAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// ProcAlterTargetProcExists 目标存储过程是否存在
// RULE: MSP-L3-002
func (v *ProcAlterVldr) ProcAlterTargetProcExists(s *models.Statement, r *models.Rule) {
}

// ProcDropVldr 删除存储过程语句相关的审核规则
type ProcDropVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *ProcDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ProcDropVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ProcDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ProcDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// ProcDropTargetDatabaseExists 目标库是否存在
// RULE: DSP-L3-001
func (v *ProcDropVldr) ProcDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// ProcDropTargetProcExists 目标存储过程是否存在
// RULE: DSP-L3-002
func (v *ProcDropVldr) ProcDropTargetProcExists(s *models.Statement, r *models.Rule) {
}
