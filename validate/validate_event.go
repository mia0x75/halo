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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *EventCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *EventCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// EventCreateEventNameQualified 事件名标识符规则
// RULE: CEV-L2-001
func (v *EventCreateVldr) EventCreateEventNameQualified(s *models.Statement, r *models.Rule) {
}

// EventCreateEventNameLowerCaseRequired 事件名大小写规则
// RULE: CEV-L2-002
func (v *EventCreateVldr) EventCreateEventNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// EventCreateEventNameMaxLength 事件名长度规则
// RULE: CEV-L2-003
func (v *EventCreateVldr) EventCreateEventNameMaxLength(s *models.Statement, r *models.Rule) {
}

// EventCreateEventNamePrefixRequired 事件名前缀规则
// RULE: CEV-L2-004
func (v *EventCreateVldr) EventCreateEventNamePrefixRequired(s *models.Statement, r *models.Rule) {
}

// EventCreateTargetDatabaseExists 目标库是否存在
// RULE: CEV-L3-001
func (v *EventCreateVldr) EventCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// EventCreateTargetEventExists 目标事件是否存在
// RULE: CEV-L3-002
func (v *EventCreateVldr) EventCreateTargetEventExists(s *models.Statement, r *models.Rule) {
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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *EventAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *EventAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// EventAlterTargetDatabaseExists 目标库是否存在
// RULE: MEV-L3-001
func (v *EventAlterVldr) EventAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// EventAlterTargetEventExists 目标事件是否存在
// RULE: MEV-L3-002
func (v *EventAlterVldr) EventAlterTargetEventExists(s *models.Statement, r *models.Rule) {
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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *EventDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *EventDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// EventDropTargetDatabaseExists 目标库是否存在
// RULE: DEV-L3-001
func (v *EventDropVldr) EventDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// EventDropTargetEventExists 目标事件是否存在
// RULE: DEV-L3-002
func (v *EventDropVldr) EventDropTargetEventExists(s *models.Statement, r *models.Rule) {
}
