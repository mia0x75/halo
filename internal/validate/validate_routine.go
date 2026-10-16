package validate

import (
	"sync"
)

// ProcCreateVldr 创建存储过程语句相关的审核规则
type ProcCreateVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *ProcCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ProcCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// ProcNameQualified 存储过程名标识符规则
// RULE: CSP-L2-001
func (v *ProcCreateVldr) createProcProcNameQualified(s *Subject, r *Rule) {
}

// ProcNameLowerCaseRequired 存储过程名大小写规则
// RULE: CSP-L2-002
func (v *ProcCreateVldr) createProcProcNameLowerCaseRequired(s *Subject, r *Rule) {
}

// ProcNameMaxLength 存储过程名长度规则
// RULE: CSP-L2-003
func (v *ProcCreateVldr) createProcProcNameMaxLength(s *Subject, r *Rule) {
}

// ProcNamePrefixRequired 存储过程名前缀规则
// RULE: CSP-L2-004
func (v *ProcCreateVldr) createProcProcNamePrefixRequired(s *Subject, r *Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CSP-L3-001
func (v *ProcCreateVldr) createProcTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetProcDoesNotExist 目标存储过程是否存在
// RULE: CSP-L3-002
func (v *ProcCreateVldr) createProcTargetProcDoesNotExist(s *Subject, r *Rule) {
}

// ProcAlterVldr 修改存储过程语句相关的审核规则
type ProcAlterVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *ProcAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ProcAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MSP-L3-001
func (v *ProcAlterVldr) alterProcTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetProcDoesNotExist 目标存储过程是否存在
// RULE: MSP-L3-002
func (v *ProcAlterVldr) alterProcTargetProcDoesNotExist(s *Subject, r *Rule) {
}

// ProcDropVldr 删除存储过程语句相关的审核规则
type ProcDropVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *ProcDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *ProcDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DSP-L3-001
func (v *ProcDropVldr) dropProcTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetProcDoesNotExist 目标存储过程是否存在
// RULE: DSP-L3-002
func (v *ProcDropVldr) dropProcTargetProcDoesNotExist(s *Subject, r *Rule) {
}
