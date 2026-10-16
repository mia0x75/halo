package validate

import (
	"sync"
)

// TriggerCreateVldr 创建触发器语句相关的审核规则
type TriggerCreateVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *TriggerCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TriggerCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TriggerNameQualified 触发器名标识符规则
// RULE: CTG-L2-001
func (v *TriggerCreateVldr) createTriggerTriggerNameQualified(s *Subject, r *Rule) {
}

// TriggerNameLowerCaseRequired 触发器名大小写规则
// RULE: CTG-L2-002
func (v *TriggerCreateVldr) createTriggerTriggerNameLowerCaseRequired(s *Subject, r *Rule) {
}

// TriggerNameMaxLength 触发器名长度规则
// RULE: CTG-L2-003
func (v *TriggerCreateVldr) createTriggerTriggerNameMaxLength(s *Subject, r *Rule) {
}

// TriggerPrefixRequired 触发器名前缀规则
// RULE: CTG-L2-004
func (v *TriggerCreateVldr) createTriggerTriggerPrefixRequired(s *Subject, r *Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CTG-L3-001
func (v *TriggerCreateVldr) createTriggerTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: CTG-L3-002
func (v *TriggerCreateVldr) createTriggerTargetTableDoesNotExist(s *Subject, r *Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE: CTG-L3-003
func (v *TriggerCreateVldr) createTriggerTargetTriggerDoesNotExist(s *Subject, r *Rule) {
}

// TriggerAlterVldr 修改触发器语句相关的审核规则
type TriggerAlterVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *TriggerAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TriggerAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE:DTG-L3-001
func (v *TriggerAlterVldr) alterTriggerTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: MTG-L3-002
func (v *TriggerAlterVldr) alterTriggerTargetTableDoesNotExist(s *Subject, r *Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE: MTG-L3-003
func (v *TriggerAlterVldr) alterTriggerTargetTriggerDoesNotExist(s *Subject, r *Rule) {
}

// TriggerDropVldr 删除触发器语句相关的审核规则
type TriggerDropVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *TriggerDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TriggerDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE:DTG-L3-001
func (v *TriggerDropVldr) dropTriggerTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE:DTG-L3-002
func (v *TriggerDropVldr) dropTriggerTargetTableDoesNotExist(s *Subject, r *Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE:DTG-L3-003
func (v *TriggerDropVldr) dropTriggerTargetTriggerDoesNotExist(s *Subject, r *Rule) {
}
