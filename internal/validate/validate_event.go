package validate

import (
	"sync"
)

// EventCreateVldr 创建事件语句相关的审核规则
type EventCreateVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *EventCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *EventCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// EventNameQualified 事件名标识符规则
// RULE: CEV-L2-001
func (v *EventCreateVldr) createEventEventNameQualified(s *Subject, r *Rule) {
}

// EventNameLowerCaseRequired 事件名大小写规则
// RULE: CEV-L2-002
func (v *EventCreateVldr) createEventEventNameLowerCaseRequired(s *Subject, r *Rule) {
}

// EventNameMaxLength 事件名长度规则
// RULE: CEV-L2-003
func (v *EventCreateVldr) createEventEventNameMaxLength(s *Subject, r *Rule) {
}

// EventNamePrefixRequired 事件名前缀规则
// RULE: CEV-L2-004
func (v *EventCreateVldr) createEventEventNamePrefixRequired(s *Subject, r *Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CEV-L3-001
func (v *EventCreateVldr) createEventTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetEventDoesNotExist 目标事件是否存在
// RULE: CEV-L3-002
func (v *EventCreateVldr) createEventTargetEventDoesNotExist(s *Subject, r *Rule) {
}

// EventAlterVldr 修改事件语句相关的审核规则
type EventAlterVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *EventAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *EventAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: MEV-L3-001
func (v *EventAlterVldr) alterEventTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetEventDoesNotExist 目标事件是否存在
// RULE: MEV-L3-002
func (v *EventAlterVldr) alterEventTargetEventDoesNotExist(s *Subject, r *Rule) {
}

// EventDropVldr 删除事件语句相关的审核规则
type EventDropVldr struct {
	vldr
}

// Enabled 当前规则组是否生效
func (v *EventDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *EventDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DEV-L3-001
func (v *EventDropVldr) dropEventTargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetEventDoesNotExist 目标事件是否存在
// RULE: DEV-L3-002
func (v *EventDropVldr) dropEventTargetEventDoesNotExist(s *Subject, r *Rule) {
}
