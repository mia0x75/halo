package validate

import (
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// TriggerCreateVldr 创建触发器语句相关的审核规则
type TriggerCreateVldr struct {
	vldr

	ct *ast.CreateTriggerStmt
}

// Enabled 当前规则组是否生效
func (v *TriggerCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TriggerCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		ct, ok := s.StmtNode.(*ast.CreateTriggerStmt)
		if !ok {
			continue
		}
		v.ct = ct
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// triggerDB resolves the schema of a trigger statement, falling back to the
// table's schema when the trigger name is unqualified.
func triggerDB(trigger, table *ast.TableName) string {
	if db := refOf(trigger).Schema; db != "" {
		return db
	}
	return refOf(table).Schema
}

// TriggerNameQualified 触发器名标识符规则
// RULE: CTG-L2-001
func (v *TriggerCreateVldr) TriggerNameQualified(s *Subject, r *Rule) {
	nameRegexRule(s, r, refOf(v.ct.Name).Name, true)
}

// TriggerNameLowerCaseRequired 触发器名大小写规则
// RULE: CTG-L2-002
func (v *TriggerCreateVldr) TriggerNameLowerCaseRequired(s *Subject, r *Rule) {
	nameRegexRule(s, r, refOf(v.ct.Name).Name, false)
}

// TriggerNameMaxLength 触发器名长度规则
// RULE: CTG-L2-003
func (v *TriggerCreateVldr) TriggerNameMaxLength(s *Subject, r *Rule) {
	nameMaxLengthRule(s, r, refOf(v.ct.Name).Name)
}

// TriggerPrefixRequired 触发器名前缀规则
// RULE: CTG-L2-004
func (v *TriggerCreateVldr) TriggerPrefixRequired(s *Subject, r *Rule) {
	namePrefixRule(s, r, refOf(v.ct.Name).Name)
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: CTG-L3-001
func (v *TriggerCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, triggerDB(v.ct.Name, v.ct.Table))
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: CTG-L3-002
func (v *TriggerCreateVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
	table := refOf(v.ct.Table)
	v.requireTable(s, r, table.Schema, table.Name)
}

// TargetTriggerDoesNotExist 目标触发器是否存在（创建时要求不存在）
// RULE: CTG-L3-003
func (v *TriggerCreateVldr) TargetTriggerDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.ct.Name)
	v.requireObjectAbsent(s, r, triggerDB(v.ct.Name, v.ct.Table), meta.ObjectTrigger, ref.Name)
}

// TriggerAlterVldr 修改触发器语句相关的审核规则
//
// MySQL has no ALTER TRIGGER statement, so no subject ever reaches this group;
// the MTG rules are matched by the seed but cannot fire. The handlers stay as
// documented no-ops until a dialect exposes ALTER TRIGGER.
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
// RULE: MTG-L3-001
func (v *TriggerAlterVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: MTG-L3-002
func (v *TriggerAlterVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在
// RULE: MTG-L3-003
func (v *TriggerAlterVldr) TargetTriggerDoesNotExist(s *Subject, r *Rule) {
}

// TriggerDropVldr 删除触发器语句相关的审核规则
type TriggerDropVldr struct {
	vldr

	dt *ast.DropTriggerStmt
}

// Enabled 当前规则组是否生效
func (v *TriggerDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TriggerDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		dt, ok := s.StmtNode.(*ast.DropTriggerStmt)
		if !ok {
			continue
		}
		v.dt = dt
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: DTG-L3-001
func (v *TriggerDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	v.requireDatabase(s, r, refOf(v.dt.Name).Schema)
}

// TargetTableDoesNotExist 目标表是否存在
//
// MySQL's DROP TRIGGER names only the trigger, so the AST carries no table and
// this rule cannot be evaluated; it stays a documented no-op.
// RULE: DTG-L3-002
func (v *TriggerDropVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
}

// TargetTriggerDoesNotExist 目标触发器是否存在（删除时要求已存在）
// RULE: DTG-L3-003
func (v *TriggerDropVldr) TargetTriggerDoesNotExist(s *Subject, r *Rule) {
	ref := refOf(v.dt.Name)
	v.requireObjectPresent(s, r, ref.Schema, meta.ObjectTrigger, ref.Name)
}
