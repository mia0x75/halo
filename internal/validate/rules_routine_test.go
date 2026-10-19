package validate_test

import (
	"testing"

	"github.com/mia0x75/parser"
	"github.com/stretchr/testify/assert"

	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/validate"
)

// routineCases cover the storage-program groups (functions, procedures,
// triggers, events) plus ALTER VIEW. The unit harness reuses unitContext, whose
// halo_* objects provide the "already exists" branch; the integration harness
// creates the same objects on the live target.
var routineCases = []RuleTestCase{
	// CFU 新建函数
	{"CFU-L2-001", "create function fn_ok() returns int deterministic return 1;", true},
	{"CFU-L2-001", "create function `fn bad`() returns int deterministic return 1;", false},
	{"CFU-L2-002", "create function fn_ok() returns int deterministic return 1;", true},
	{"CFU-L2-002", "create function fn_Bad() returns int deterministic return 1;", false},
	{"CFU-L2-003", "create function fn_ok() returns int deterministic return 1;", true},
	{"CFU-L2-003", "create function fn_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa() returns int deterministic return 1;", false},
	{"CFU-L2-004", "create function fn_ok() returns int deterministic return 1;", true},
	{"CFU-L2-004", "create function ok_fn() returns int deterministic return 1;", false},
	{"CFU-L3-001", "create function starwars.fn_ok() returns int deterministic return 1;", true},
	{"CFU-L3-001", "create function nodb.fn_x() returns int deterministic return 1;", false},
	{"CFU-L3-002", "create function starwars.fn_new() returns int deterministic return 1;", true},
	{"CFU-L3-002", "create function starwars.halo_fn() returns int deterministic return 1;", false},

	// MFU 修改函数
	{"MFU-L3-001", "alter function starwars.halo_fn comment 'x';", true},
	{"MFU-L3-001", "alter function nodb.fn_x comment 'x';", false},
	{"MFU-L3-002", "alter function starwars.halo_fn comment 'x';", true},
	{"MFU-L3-002", "alter function starwars.no_fn comment 'x';", false},

	// DFU 删除函数
	{"DFU-L3-001", "drop function starwars.halo_fn;", true},
	{"DFU-L3-001", "drop function nodb.fn_x;", false},
	{"DFU-L3-002", "drop function starwars.halo_fn;", true},
	{"DFU-L3-002", "drop function starwars.no_fn;", false},

	// CSP 新建存储过程
	{"CSP-L2-001", "create procedure sp_ok() select 1;", true},
	{"CSP-L2-001", "create procedure `sp bad`() select 1;", false},
	{"CSP-L2-002", "create procedure sp_ok() select 1;", true},
	{"CSP-L2-002", "create procedure sp_Bad() select 1;", false},
	{"CSP-L2-003", "create procedure sp_ok() select 1;", true},
	{"CSP-L2-003", "create procedure sp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa() select 1;", false},
	{"CSP-L2-004", "create procedure sp_ok() select 1;", true},
	{"CSP-L2-004", "create procedure ok_sp() select 1;", false},
	{"CSP-L3-001", "create procedure starwars.sp_ok() select 1;", true},
	{"CSP-L3-001", "create procedure nodb.sp_x() select 1;", false},
	{"CSP-L3-002", "create procedure starwars.sp_new() select 1;", true},
	{"CSP-L3-002", "create procedure starwars.halo_proc() select 1;", false},

	// MSP 修改存储过程
	{"MSP-L3-001", "alter procedure starwars.halo_proc comment 'x';", true},
	{"MSP-L3-001", "alter procedure nodb.sp_x comment 'x';", false},
	{"MSP-L3-002", "alter procedure starwars.halo_proc comment 'x';", true},
	{"MSP-L3-002", "alter procedure starwars.no_proc comment 'x';", false},

	// DSP 删除存储过程
	{"DSP-L3-001", "drop procedure starwars.halo_proc;", true},
	{"DSP-L3-001", "drop procedure nodb.sp_x;", false},
	{"DSP-L3-002", "drop procedure starwars.halo_proc;", true},
	{"DSP-L3-002", "drop procedure starwars.no_proc;", false},

	// CTG 新建触发器
	{"CTG-L2-001", "create trigger tg_ok before insert on t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L2-001", "create trigger `tg bad` before insert on t1 for each row insert into t1(id) values(1);", false},
	{"CTG-L2-002", "create trigger tg_ok before insert on t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L2-002", "create trigger tg_Bad before insert on t1 for each row insert into t1(id) values(1);", false},
	{"CTG-L2-003", "create trigger tg_ok before insert on t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L2-003", "create trigger tg_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa before insert on t1 for each row insert into t1(id) values(1);", false},
	{"CTG-L2-004", "create trigger tg_ok before insert on t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L2-004", "create trigger ok_tg before insert on t1 for each row insert into t1(id) values(1);", false},
	{"CTG-L3-001", "create trigger tg_ok before insert on starwars.t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L3-001", "create trigger tg_ok before insert on nodb.t1 for each row insert into t1(id) values(1);", false},
	{"CTG-L3-002", "create trigger tg_ok before insert on starwars.t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L3-002", "create trigger tg_ok before insert on starwars.notable for each row insert into t1(id) values(1);", false},
	{"CTG-L3-003", "create trigger starwars.tg_new before insert on starwars.t1 for each row insert into t1(id) values(1);", true},
	{"CTG-L3-003", "create trigger starwars.halo_trg before insert on starwars.t1 for each row insert into t1(id) values(1);", false},

	// DTG 删除触发器
	{"DTG-L3-001", "drop trigger starwars.tg_x;", true},
	{"DTG-L3-001", "drop trigger nodb.tg_x;", false},
	{"DTG-L3-003", "drop trigger starwars.halo_trg;", true},
	{"DTG-L3-003", "drop trigger starwars.no_trg;", false},

	// CEV 新建事件
	{"CEV-L2-001", "create event ev_ok on schedule every 1 day do select 1;", true},
	{"CEV-L2-001", "create event `ev bad` on schedule every 1 day do select 1;", false},
	{"CEV-L2-002", "create event ev_ok on schedule every 1 day do select 1;", true},
	{"CEV-L2-002", "create event ev_Bad on schedule every 1 day do select 1;", false},
	{"CEV-L2-003", "create event ev_ok on schedule every 1 day do select 1;", true},
	{"CEV-L2-003", "create event ev_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa on schedule every 1 day do select 1;", false},
	{"CEV-L2-004", "create event ev_ok on schedule every 1 day do select 1;", true},
	{"CEV-L2-004", "create event ok_ev on schedule every 1 day do select 1;", false},
	{"CEV-L3-001", "create event ev_ok on schedule every 1 day do select 1;", true},
	{"CEV-L3-001", "create event nodb.ev_x on schedule every 1 day do select 1;", false},
	{"CEV-L3-002", "create event starwars.ev_new on schedule every 1 day do select 1;", true},
	{"CEV-L3-002", "create event starwars.halo_ev on schedule every 1 day do select 1;", false},

	// MEV 修改事件
	{"MEV-L3-001", "alter event starwars.halo_ev disable;", true},
	{"MEV-L3-001", "alter event nodb.ev_x disable;", false},
	{"MEV-L3-002", "alter event starwars.halo_ev disable;", true},
	{"MEV-L3-002", "alter event starwars.no_ev disable;", false},

	// DEV 删除事件
	{"DEV-L3-001", "drop event starwars.halo_ev;", true},
	{"DEV-L3-001", "drop event nodb.ev_x;", false},
	{"DEV-L3-002", "drop event starwars.halo_ev;", true},
	{"DEV-L3-002", "drop event starwars.no_ev;", false},
}

// TestRulesRoutine_Unit drives the storage-program cases without a database.
func TestRulesRoutine_Unit(t *testing.T) {
	rules := seedRules(t)
	ctx := unitContext()
	p := parser.New()
	assert := assert.New(t)

	for _, c := range routineCases {
		mr, ok := rules[c.RuleName]
		if !ok {
			t.Errorf("rule %s not found in seed", c.RuleName)
			continue
		}
		stmt, err := p.ParseOneStmt(c.Text, "", "")
		if err != nil {
			t.Errorf("语法错误: %s: %s", c.RuleName, err.Error())
			continue
		}
		s := &validate.Subject{Content: c.Text, Violations: &models.Violations{}, StmtNode: stmt}
		ctx.Subjects = []*validate.Subject{s}
		rule := *mr
		rule.Bitwise |= 1
		ctx.Rules = []*validate.Rule{{Rule: rule}}
		validate.RunRegistry(ctx)

		valid := s.Violations.Marshal() == ""
		assert.Equal(c.Valid, valid, "%s: %s", c.RuleName, c.Text)
	}
}

// TestRulesRoutine_Integration drives the same cases against the live target so
// the object metadata loader is exercised end to end.
func TestRulesRoutine_Integration(t *testing.T) {
	a := buildApp(t)
	ensureFixture(t, a)
	ctx := buildContext(t)
	p := parser.New()

	for _, c := range routineCases {
		valid, ok := runCase(t, ctx, p, c)
		if !ok {
			continue
		}
		if c.Valid != valid {
			t.Errorf("%s: %s: want valid=%v, got %v", c.RuleName, c.Text, c.Valid, valid)
		}
	}
}
