package validate

import (
	"sync"

	"github.com/mia0x75/halo/models"
)

// TriggerCreateVldr 创建触发器语句相关的审核规则
type TriggerCreateVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *TriggerCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TriggerCreateVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TriggerCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TriggerCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TriggerCreateTriggerNameQualified 触发器名标识符规则
// RULE: CTG-L2-001
func (v *TriggerCreateVldr) TriggerCreateTriggerNameQualified(s *models.Statement, r *models.Rule) {
}

// TriggerCreateTriggerNameLowerCaseRequired 触发器名大小写规则
// RULE: CTG-L2-002
func (v *TriggerCreateVldr) TriggerCreateTriggerNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// TriggerCreateTriggerNameMaxLength 触发器名长度规则
// RULE: CTG-L2-003
func (v *TriggerCreateVldr) TriggerCreateTriggerNameMaxLength(s *models.Statement, r *models.Rule) {
}

// TriggerCreateTriggerPrefixRequired 触发器名前缀规则
// RULE: CTG-L2-004
func (v *TriggerCreateVldr) TriggerCreateTriggerPrefixRequired(s *models.Statement, r *models.Rule) {
}

// TriggerCreateTargetDatabaseExists 目标库是否存在
// RULE: CTG-L3-001
func (v *TriggerCreateVldr) TriggerCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TriggerCreateTargetTableExists 目标表是否存在
// RULE: CTG-L3-002
func (v *TriggerCreateVldr) TriggerCreateTargetTableExists(s *models.Statement, r *models.Rule) {
}

// TriggerCreateTargetTriggerExists 目标触发器是否存在
// RULE: CTG-L3-003
func (v *TriggerCreateVldr) TriggerCreateTargetTriggerExists(s *models.Statement, r *models.Rule) {
}

// TriggerAlterTargetDatabaseExists 目标库是否存在
// RULE: MTG-L3-001
func (v *TriggerCreateVldr) TriggerAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TriggerAlterVldr 修改触发器语句相关的审核规则
type TriggerAlterVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *TriggerAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TriggerAlterVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TriggerAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TriggerAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TriggerAlterTargetDatabaseExists 目标库是否存在
// RULE:DTG-L3-001
func (v *TriggerAlterVldr) TriggerAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TriggerAlterTargetTableExists 目标表是否存在
// RULE: MTG-L3-002
func (v *TriggerAlterVldr) TriggerAlterTargetTableExists(s *models.Statement, r *models.Rule) {
}

// TriggerAlterTargetTriggerExists 目标触发器是否存在
// RULE: MTG-L3-003
func (v *TriggerAlterVldr) TriggerAlterTargetTriggerExists(s *models.Statement, r *models.Rule) {
}

// TriggerDropVldr 删除触发器语句相关的审核规则
type TriggerDropVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *TriggerDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TriggerDropVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TriggerDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TriggerDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TriggerDropTargetDatabaseExists 目标库是否存在
// RULE:DTG-L3-001
func (v *TriggerDropVldr) TriggerDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TriggerDropTargetTableExists 目标表是否存在
// RULE:DTG-L3-002
func (v *TriggerDropVldr) TriggerDropTargetTableExists(s *models.Statement, r *models.Rule) {
}

// TriggerDropTargetTriggerExists 目标触发器是否存在
// RULE:DTG-L3-003
func (v *TriggerDropVldr) TriggerDropTargetTriggerExists(s *models.Statement, r *models.Rule) {
}
