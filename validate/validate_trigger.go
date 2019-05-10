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

// Validate 规则组的审核入口
func (v *TriggerCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TriggerNameQualified 触发器名标识符规则
// RULE: CTG-L2-001
func (v *TriggerCreateVldr) TriggerNameQualified(s *models.Statement, r *models.Rule) {
}

// TriggerNameLowerCaseRequired 触发器名大小写规则
// RULE: CTG-L2-002
func (v *TriggerCreateVldr) TriggerNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
}

// TriggerNameMaxLength 触发器名长度规则
// RULE: CTG-L2-003
func (v *TriggerCreateVldr) TriggerNameMaxLength(s *models.Statement, r *models.Rule) {
}

// TriggerPrefixRequired 触发器名前缀规则
// RULE: CTG-L2-004
func (v *TriggerCreateVldr) TriggerPrefixRequired(s *models.Statement, r *models.Rule) {
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CTG-L3-001
func (v *TriggerCreateVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: CTG-L3-002
func (v *TriggerCreateVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE: CTG-L3-003
func (v *TriggerCreateVldr) TargetTriggerDoesNotExist(s *models.Statement, r *models.Rule) {
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

// Validate 规则组的审核入口
func (v *TriggerAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE:DTG-L3-001
func (v *TriggerAlterVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: MTG-L3-002
func (v *TriggerAlterVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE: MTG-L3-003
func (v *TriggerAlterVldr) TargetTriggerDoesNotExist(s *models.Statement, r *models.Rule) {
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

// Validate 规则组的审核入口
func (v *TriggerDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE:DTG-L3-001
func (v *TriggerDropVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE:DTG-L3-002
func (v *TriggerDropVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE:DTG-L3-003
func (v *TriggerDropVldr) TargetTriggerDoesNotExist(s *models.Statement, r *models.Rule) {
}
