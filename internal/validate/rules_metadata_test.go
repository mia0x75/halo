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

// metadataContext holds the columns/indexes the metadata-backed rules need:
//
//	t1     id INT, body TEXT           indexes idx_id(id), idx_multi(id,body)
//	users  id INT
//	tt     id INT, created/updated TIMESTAMP
//	single id INT, created TIMESTAMP
//
// A separate context keeps these extra columns from disturbing the cases that
// rely on unitContext's minimal t1.
func metadataContext() *validate.Context {
	col := func(name string, t meta.ColumnType) *meta.Column {
		return &meta.Column{Name: name, Type: t}
	}
	return &validate.Context{
		Ticket: &domain.Ticket{Database: "starwars"},
		Databases: []models.Database{
			{Name: "starwars", Charset: "utf8mb4", Collate: "utf8mb4_unicode_ci"},
			{Name: "db", Charset: "utf8mb4", Collate: "utf8mb4_unicode_ci"},
		},
		Tables: map[string][]*meta.Table{
			"starwars": {
				{
					Name: "t1",
					Columns: map[string]*meta.Column{
						"id":   col("id", "INT"),
						"body": col("body", "TEXT"),
					},
					Indexes: map[string]*meta.Index{
						"idx_id":    {Name: "idx_id", Cols: []string{"id"}},
						"idx_multi": {Name: "idx_multi", Cols: []string{"id", "body"}},
					},
				},
				{
					Name:    "users",
					Columns: map[string]*meta.Column{"id": col("id", "INT")},
					Indexes: map[string]*meta.Index{},
				},
				{
					Name: "tt",
					Columns: map[string]*meta.Column{
						"id":      col("id", "INT"),
						"created": col("created", "TIMESTAMP"),
						"updated": col("updated", "TIMESTAMP"),
					},
					Indexes: map[string]*meta.Index{},
				},
				{
					Name: "single",
					Columns: map[string]*meta.Column{
						"id":      col("id", "INT"),
						"created": col("created", "TIMESTAMP"),
					},
					Indexes: map[string]*meta.Index{},
				},
			},
		},
	}
}

// TestRulesMetadata_Unit covers the index-overlay, timestamp-count and rename
// rules, all of which were empty handlers. The seed disables most of them
// (even Bitwise), so the rule under test is force-enabled here.
func TestRulesMetadata_Unit(t *testing.T) {
	rules := seedRules(t)
	ctx := metadataContext()
	p := parser.New()
	assert := assert.New(t)

	cases := []RuleTestCase{
		// CIX-L3-004 索引覆盖 (t1 已有 idx_id(id)/idx_multi(id,body))
		{"CIX-L3-004", "create index idx_full on t1(id);", false},
		{"CIX-L3-004", "create index idx_prefix on t1(id);", false},
		{"CIX-L3-004", "create index idx_body on t1(body);", true},

		// MTB-L3-014 覆盖索引 (alter add index)
		{"MTB-L3-014", "alter table t1 add index idx_dup(id);", false},
		{"MTB-L3-014", "alter table t1 add index idx_prefix(id);", false},
		{"MTB-L3-014", "alter table t1 add index idx_body(body);", true},

		// MTB-L3-007 仅允许一个时间戳列
		{"MTB-L3-007", "alter table tt add column x int;", false},
		{"MTB-L3-007", "alter table single add column x int;", true},
		{"MTB-L3-007", "alter table single add column created2 timestamp;", false},
		{"MTB-L3-007", "alter table tt drop column updated;", true},

		// RTB-L3-001 源库必须已存在
		{"RTB-L3-001", "rename table users to users2;", true},
		{"RTB-L3-001", "rename table mock.users to users2;", false},

		// RTB-L3-002 源表必须已存在
		{"RTB-L3-002", "rename table users to users2;", true},
		{"RTB-L3-002", "rename table nosuch to users2;", false},

		// RTB-L3-003 目标库必须已存在
		{"RTB-L3-003", "rename table users to users2;", true},
		{"RTB-L3-003", "rename table users to db.users2;", true},
		{"RTB-L3-003", "rename table users to mock.users2;", false},

		// RTB-L3-004 目标表必须不存在
		{"RTB-L3-004", "rename table users to users2;", true},
		{"RTB-L3-004", "rename table users to t1;", false},
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
		rule := *mr
		rule.Bitwise |= 1
		ctx.Rules = []*validate.Rule{{Rule: rule}}
		validate.RunRegistry(ctx)

		valid := s.Violations.Marshal() == ""
		assert.Equal(c.Valid, valid, "%s: %s", c.RuleName, c.Text)
	}
}
