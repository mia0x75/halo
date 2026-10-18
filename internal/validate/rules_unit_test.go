package validate_test

import (
	"testing"

	"github.com/mia0x75/parser"
	"github.com/stretchr/testify/assert"

	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/validate"
)

// TestRules_Unit runs every rule case without a database: rules come from the
// seed, metadata is synthesised, and only the addressed rule is executed. It
// mirrors TestGeneralRules but never touches the metadata or target database.
func TestRules_Unit(t *testing.T) {
	rules := seedRules(t)
	ctx := unitContext()
	p := parser.New()
	assert := assert.New(t)

	cases := []RuleTestCase{}
	cases = append(cases, tableCases...)
	cases = append(cases, viewCases...)
	cases = append(cases, updateCasess...)
	cases = append(cases, replaceCases...)
	cases = append(cases, insertCases...)
	cases = append(cases, indexCases...)
	cases = append(cases, deleteCases...)
	cases = append(cases, databaseCases...)
	cases = append(cases, miscCases...)

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
		// The seed ships some implemented rules disabled (even Bitwise); force
		// the rule on so a case can exercise its handler regardless of seed state.
		rule := *mr
		rule.Bitwise |= 1
		ctx.Rules = []*validate.Rule{{Rule: rule}}
		validate.RunRegistry(ctx)

		valid := s.Violations.Marshal() == ""
		assert.Equal(c.Valid, valid, "%s: %s", c.RuleName, c.Text)
	}
}

// TestMergeStatement_Unit runs the multi-statement (Misc) cases without a DB.
func TestMergeStatement_Unit(t *testing.T) {
	rules := seedRules(t)
	ctx := unitContext()
	p := parser.New()
	assert := assert.New(t)

	cases := []RuleTestCase{
		{"MSC-L1-004", "CREATE TABLE t1 (c1 INT); CREATE INDEX index_1 ON t1(c1);", false},
		{"MSC-L1-004", "CREATE TABLE t1 (c1 INT); ALTER TABLE t1 ADD c2 INT;", false},
		{"MSC-L1-004", "ALTER TABLE t1 ADD c1 INT; ALTER TABLE t1 ADD c2 INT;", false},
		{"MSC-L1-004", "ALTER TABLE t1 ADD c1 INT; ALTER TABLE t1 DROP INDEX index_2;", false},
		{"MSC-L1-004", "CREATE INDEX index_1 ON t1(c1); ALTER TABLE t1 DROP INDEX index_2;", false},
		{"MSC-L1-008", "SELECT 1;", true},
		{"MSC-L1-008", "SELECT * FROM t1; ALTER TABLE t1 ADD c1 INT;", false},
		{"MSC-L1-008", "ALTER TABLE t1 ADD c1 INT; SELECT * FROM t1;", false},
		{"MSC-L1-004", "SELECT 1;", true},
	}

	for _, c := range cases {
		mr, ok := rules[c.RuleName]
		if !ok {
			t.Errorf("rule %s not found in seed", c.RuleName)
			continue
		}
		nodes, _, err := p.Parse(c.Text, "", "")
		if err != nil {
			t.Errorf("语法错误: %s", err.Error())
			continue
		}
		subjects := make([]*validate.Subject, 0, len(nodes))
		for _, node := range nodes {
			subjects = append(subjects, &validate.Subject{
				Content:    c.Text,
				Violations: &models.Violations{},
				StmtNode:   node,
			})
		}
		ctx.Subjects = subjects
		ctx.Rules = []*validate.Rule{{Rule: *mr}}
		validate.RunRegistry(ctx)

		valid := true
		for _, s := range subjects {
			if s.Violations.Marshal() != "" {
				valid = false
				break
			}
		}
		assert.Equal(c.Valid, valid, "%s: %s", c.RuleName, c.Text)
	}
}
