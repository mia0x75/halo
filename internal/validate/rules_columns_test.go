package validate_test

import (
	"testing"

	"github.com/mia0x75/parser"
	"github.com/stretchr/testify/assert"

	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/validate"
)

// columnContext synthesises a table whose columns exercise the nullability,
// default and type branches of the column-level rules. It is separate from
// unitContext so the "table does not exist" cases there stay intact.
//
// users:
//
//	id     INT     NOT NULL auto_increment
//	name   VARCHAR NOT NULL (no default -> must be provided on insert)
//	status INT     NOT NULL DEFAULT 0
//	bio    TEXT    NULL
//	note   TEXT    NULL
func columnContext() *validate.Context {
	col := func(name string, t meta.ColumnType, nullable bool, def string, auto bool) *meta.Column {
		return &meta.Column{Name: name, Type: t, Nullable: nullable, Default: def, AutoIncrement: auto}
	}
	return &validate.Context{
		Ticket: &domain.Ticket{Database: "starwars"},
		Databases: []models.Database{
			{Name: "starwars", Charset: "utf8mb4", Collate: "utf8mb4_unicode_ci"},
		},
		Tables: map[string][]*meta.Table{
			"starwars": {
				{
					Name:    "users",
					Engine:  "InnoDB",
					Charset: "utf8mb4",
					Collate: "utf8mb4_unicode_ci",
					Columns: map[string]*meta.Column{
						"id":     col("id", "INT", false, "", true),
						"name":   col("name", "VARCHAR", false, "", false),
						"status": col("status", "INT", false, "0", false),
						"bio":    col("bio", "TEXT", true, "", false),
						"note":   col("note", "TEXT", true, "", false),
					},
					Indexes: map[string]*meta.Index{},
				},
			},
		},
	}
}

// TestRulesColumns_Unit covers the column-level rules that previously had empty
// handlers. It shares the database-free harness of TestRules_Unit.
func TestRulesColumns_Unit(t *testing.T) {
	rules := seedRules(t)
	ctx := columnContext()
	p := parser.New()
	assert := assert.New(t)

	cases := []RuleTestCase{
		// INS-L3-003 插入列必须存在 (users: id,name,status,bio,note)
		{"INS-L3-003", "insert into users(id,name) values(1,'a');", true},
		{"INS-L3-003", "insert into users(c9) values(1);", false},

		// INS-L3-004 非空列必须提供值 (name 非空无默认, id 自增, status 有默认)
		{"INS-L3-004", "insert into users(id,name) values(1,'a');", true},
		{"INS-L3-003", "insert into users(id,name,status) values(1,'a',1);", true},
		{"INS-L3-004", "insert into users(id,status) values(1,1);", false},

		// RPL-L3-003 替换列必须存在
		{"RPL-L3-003", "replace into users(id,name) values(1,'a');", true},
		{"RPL-L3-003", "replace into users(c9) values(1);", false},

		// RPL-L3-004 非空列必须有值
		{"RPL-L3-004", "replace into users(id,name) values(1,'a');", true},
		{"RPL-L3-004", "insert into users(id) values(1);", false},

		// UPD-L3-003 SET/WHERE 引用的列必须存在
		{"UPD-L3-003", "update users set name='a' where id=1;", true},
		{"UPD-L3-003", "update users set c9=1 where id=1;", false},
		{"UPD-L3-003", "update users set name='a' where c9=1;", false},

		// UPD-L3-005 单次更新最大行数 (threshold 1000)
		{"UPD-L3-005", "update users set name='a' where id=1 limit 10;", true},
		{"UPD-L3-005", "update users set name='a' where id=1 limit 1001;", false},
		{"UPD-L3-005", "update users set name='a' where id=1;", false},

		// DEL-L3-001 单次删除最大行数 (threshold 1000)
		{"DEL-L3-001", "delete from users where id=1 limit 10;", true},
		{"DEL-L3-001", "delete from users where id=1 limit 1001;", false},
		{"DEL-L3-001", "delete from users where id=1;", false},

		// DEL-L3-004 WHERE 引用的列必须存在
		{"DEL-L3-004", "delete from users where id=1;", true},
		{"DEL-L3-004", "delete from users where c9=1;", false},

		// SEL-L3-003 返回列与过滤列必须存在 (WHERE / JOIN ON / GROUP BY / ORDER BY / HAVING)
		{"SEL-L3-003", "select id from users where id=1;", true},
		{"SEL-L3-003", "select c9 from users where id=1;", false},
		{"SEL-L3-003", "select id from users where c9=1;", false},
		{"SEL-L3-003", "select id from users where id=1 order by c9;", false},
		{"SEL-L3-003", "select id from users where id=1 group by c9;", false},
		{"SEL-L3-003", "select id from users where id=1 having c9 > 1;", false},
		{"SEL-L3-003", "select id as x from users where id=1 order by x;", true},

		// SEL-L3-004 禁止返回 BLOB/TEXT (bio, note 为 TEXT)
		{"SEL-L3-004", "select id, name from users where id=1;", true},
		{"SEL-L3-004", "select bio from users where id=1;", false},
		{"SEL-L3-004", "select id, note from users where id=1;", false},
	}

	for _, c := range cases {
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
		// The seed leaves several of these rules disabled (even Bitwise); the
		// handler itself is under test, so force the rule on.
		rule := *mr
		rule.Bitwise |= 1
		ctx.Rules = []*validate.Rule{{Rule: rule}}
		validate.RunRegistry(ctx)

		valid := s.Violations.Marshal() == ""
		assert.Equal(c.Valid, valid, "%s: %s", c.RuleName, c.Text)
	}
}
