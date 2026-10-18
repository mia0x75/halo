package validate_test

import (
	"strings"
	"testing"

	"github.com/mia0x75/parser"
	"github.com/stretchr/testify/assert"

	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/validate"
)

// insertManyValues builds a single INSERT with n value rows, used to exercise
// the statement row-count rule.
func insertManyValues(n int) string {
	var b strings.Builder
	b.WriteString("insert into t1(id) values")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("(1)")
	}
	b.WriteString(";")
	return b.String()
}

// TestRulesExtra_Unit covers rule groups the original case tables did not reach
// (index DDL, fulltext index, rename table, select guards, and the table-alter
// metadata rules). It shares the database-free harness of TestRules_Unit.
//
// Rules whose seed bitwise flag is even are disabled by the engine and are
// intentionally omitted; rules whose Walk support is incomplete (INSERT/SELECT
// schema-qualified targets, views) cannot be asserted negatively either.
func TestRulesExtra_Unit(t *testing.T) {
	rules := seedRules(t)
	ctx := unitContext()
	p := parser.New()
	assert := assert.New(t)

	cases := []RuleTestCase{
		// CIX-L3-005 索引名是否重复 (t1 already has idx_id)
		{"CIX-L3-005", "CREATE INDEX idx_id ON t1(id);", false},
		{"CIX-L3-005", "CREATE INDEX index_new ON t1(id);", true},

		// CIX-L3-007 禁止在 BLOB/TEXT 列上建索引 (t1.body is TEXT)
		{"CIX-L3-007", "CREATE INDEX index_2 ON t1(body);", false},
		{"CIX-L3-007", "CREATE INDEX index_2 ON t1(id);", true},

		// INS-L2-005 列数与值数匹配
		{"INS-L2-005", "insert into t1(a,b) values(1,2);", true},
		{"INS-L2-005", "insert into t1(a,b) values(1);", false},

		// RPL-L2-005 列数与值数匹配
		{"RPL-L2-005", "replace into t1(a,b) values(1,2);", true},
		{"RPL-L2-005", "replace into t1(a,b) values(1);", false},

		// SEL-L2-001 禁止没有 WHERE
		{"SEL-L2-001", "select id from t1 where id=1;", true},
		{"SEL-L2-001", "select id from t1;", false},

		// SEL-L2-002 禁止没有 LIMIT
		{"SEL-L2-002", "select id from t1 limit 1;", true},
		{"SEL-L2-002", "select id from t1;", false},

		// SEL-L2-003 禁止 SELECT *
		{"SEL-L2-003", "select id from t1;", true},
		{"SEL-L2-003", "select * from t1;", false},

		// SEL-L2-004 禁止显式锁
		{"SEL-L2-004", "select id from t1 where id=1;", true},
		{"SEL-L2-004", "select id from t1 where id=1 for update;", false},

		// RIX-L3-001/002/003 目标库/表/索引存在
		{"RIX-L3-001", "drop index index_2 on t1;", true},
		{"RIX-L3-001", "drop index index_2 on mock.t1;", false},
		{"RIX-L3-002", "drop index index_2 on t1;", true},
		{"RIX-L3-002", "drop index index_2 on t2;", false},
		{"RIX-L3-003", "drop index idx_id on t1;", true},
		{"RIX-L3-003", "drop index index_missing on t1;", false},

		// RTB-L2-001 源表与目标表相同
		{"RTB-L2-001", "rename table t1 to t9;", true},
		{"RTB-L2-001", "rename table t1 to t1;", false},

		// RTB-L2-002 目标表名标识符
		{"RTB-L2-002", "rename table t1 to t9;", true},
		{"RTB-L2-002", "rename table t1 to `9t`;", false},

		// RTB-L2-003 目标表名必须小写
		{"RTB-L2-003", "rename table t1 to t9;", true},
		{"RTB-L2-003", "rename table t1 to T9;", false},

		// RTB-L2-004 目标表名最大长度 (20)
		{"RTB-L2-004", "rename table t1 to t1234567890123456789;", true},
		{"RTB-L2-004", "rename table t1 to t12345678901234567890;", false},

		// MTB-L2-033 禁止全文索引
		{"MTB-L2-033", "alter table t1 add fulltext index ft_1(body);", false},
		{"MTB-L2-033", "alter table t1 add index index_1(id);", true},

		// MTB-L2-035 全文索引名标识符
		{"MTB-L2-035", "alter table t1 add fulltext index ft_1(body);", true},
		{"MTB-L2-035", "alter table t1 add fulltext index `ft 1`(body);", false},

		// MTB-L2-036 全文索引名必须小写
		{"MTB-L2-036", "alter table t1 add fulltext index ft_1(body);", true},
		{"MTB-L2-036", "alter table t1 add fulltext index FT_1(body);", false},

		// MTB-L2-037 全文索引名最大长度 (10)
		{"MTB-L2-037", "alter table t1 add fulltext index ft_1234567(body);", true},
		{"MTB-L2-037", "alter table t1 add fulltext index ft_12345678(body);", false},

		// MTB-L2-038 全文索引名前缀 (^ft_[1-9])
		{"MTB-L2-038", "alter table t1 add fulltext index ft_1(body);", true},
		{"MTB-L2-038", "alter table t1 add fulltext index ft_0(body);", false},

		// MTB-L3-024 修改表注解不可为空
		{"MTB-L3-024", "alter table t1 comment = 'x';", true},
		{"MTB-L3-024", "alter table t1 comment = '';", false},

		// MTB-L3-001/002 目标库/表存在
		{"MTB-L3-001", "alter table t1 add c1 int;", true},
		{"MTB-L3-001", "alter table mock.t1 add c1 int;", false},
		{"MTB-L3-002", "alter table t1 add c1 int;", true},
		{"MTB-L3-002", "alter table t2 add c1 int;", false},

		// MTB-L3-005 新增列不得与现有列重名 (t1 has id, body)
		{"MTB-L3-005", "alter table t1 add id int;", false},
		{"MTB-L3-005", "alter table t1 add c9 int;", true},

		// MTB-L3-008 修改/删除的列必须存在
		{"MTB-L3-008", "alter table t1 modify c9 int;", false},
		{"MTB-L3-008", "alter table t1 modify id int;", true},

		// MTB-L3-013 索引名不得重复
		{"MTB-L3-013", "alter table t1 add index idx_id (id);", false},
		{"MTB-L3-013", "alter table t1 add index index_new (id);", true},

		// MTB-L3-017 删除的索引必须存在
		{"MTB-L3-017", "alter table t1 drop index index_missing;", false},
		{"MTB-L3-017", "alter table t1 drop index idx_id;", true},

		// MTB-L3-015 索引列必须存在
		{"MTB-L3-015", "alter table t1 add index index_1 (c9);", false},
		{"MTB-L3-015", "alter table t1 add index index_1 (id);", true},

		// MTB-L3-016 禁止在 BLOB/TEXT 列上索引
		{"MTB-L3-016", "alter table t1 add index index_1 (body);", false},
		{"MTB-L3-016", "alter table t1 add index index_1 (id);", true},

		// MTB-L3-004 AFTER 的参照列必须存在
		{"MTB-L3-004", "alter table t1 add c1 int after c9;", false},
		{"MTB-L3-004", "alter table t1 add c1 int after id;", true},

		// MTB-L3-021 改表名时新表名不得已存在 (t1 存在)
		{"MTB-L3-021", "alter table t1 rename to t1;", true},
		{"MTB-L3-021", "alter table t1 rename to t9;", false},

		// MTB-L2-012 列的字符集与排序规则必须匹配
		{"MTB-L2-012", "alter table t1 add c1 varchar(1) charset utf8mb4 collate utf8mb4_bin;", true},
		{"MTB-L2-012", "alter table t1 add c1 varchar(1) charset utf8mb4 collate latin1_general_ci;", false},

		// MTB-L2-013 NOT NULL 列必须有默认值
		{"MTB-L2-013", "alter table t1 add c1 int not null default 1;", true},
		{"MTB-L2-013", "alter table t1 add c1 int not null;", false},

		// MTB-L3-022 全文索引必须不存在 (t1 already has idx_id, but it is a
		// normal index; the rule keys on the fulltext index name)
		{"MTB-L3-022", "alter table t1 add fulltext index ft_new(body);", true},
		{"MTB-L3-022", "alter table t1 add fulltext index idx_id(body);", false},

		// MTB-L3-025 索引数量上限 (threshold 5, t1 has 1; adding 5 exceeds)
		{"MTB-L3-025", "alter table t1 add index i1(id);", true},
		{"MTB-L3-025", "alter table t1 add index i1(id),add index i2(id),add index i3(id),add index i4(id),add index i5(id);", false},

		// INS-L3-006 单语句行数上限 (threshold 1000)
		{"INS-L3-006", "insert into t1(id) values(1),(2);", true},
		{"INS-L3-006", insertManyValues(1001), false},
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
		ctx.Rules = []*validate.Rule{{Rule: *mr}}
		validate.RunRegistry(ctx)

		valid := s.Violations.Marshal() == ""
		assert.Equal(c.Valid, valid, "%s: %s", c.RuleName, c.Text)
	}
}
