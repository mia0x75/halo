package validate

import (
	"sync"

	"github.com/mia0x75/halo/models"
)

// EventCreateVldr 创建事件语句相关的审核规则
type EventCreateVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *EventCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
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
func (v *EventCreateVldr) EventNameQualified(s *models.Statement, r *models.Rule) {
}

// EventNameLowerCaseRequired 事件名大小写规则
// RULE: CEV-L2-002
func (v *EventCreateVldr) EventNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// EventNameMaxLength 事件名长度规则
// RULE: CEV-L2-003
func (v *EventCreateVldr) EventNameMaxLength(s *models.Statement, r *models.Rule) {
}

// EventNamePrefixRequired 事件名前缀规则
// RULE: CEV-L2-004
func (v *EventCreateVldr) EventNamePrefixRequired(s *models.Statement, r *models.Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CEV-L3-001
func (v *EventCreateVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetEventDoesNotExist 目标事件是否存在
// RULE: CEV-L3-002
func (v *EventCreateVldr) TargetEventDoesNotExist(s *models.Statement, r *models.Rule) {
}

// EventAlterVldr 修改事件语句相关的审核规则
type EventAlterVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *EventAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
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
func (v *EventAlterVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetEventDoesNotExist 目标事件是否存在
// RULE: MEV-L3-002
func (v *EventAlterVldr) TargetEventDoesNotExist(s *models.Statement, r *models.Rule) {
}

// EventDropVldr 删除事件语句相关的审核规则
type EventDropVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *EventDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
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
func (v *EventDropVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetEventDoesNotExist 目标事件是否存在
// RULE: DEV-L3-002
func (v *EventDropVldr) TargetEventDoesNotExist(s *models.Statement, r *models.Rule) {
}
