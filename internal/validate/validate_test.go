package validate

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/driver"

	"github.com/mia0x75/halo/internal/platform/config"
)

// 保留，不可以删除
var _ = driver.ValueExpr{}

// TestMain parses the config flag outside init(), so the testing package's own
// flags are already registered when Parse runs. HALO_TEST_CONFIG takes priority
// so the integration tests can point at a throwaway database.
func TestMain(m *testing.M) {
	cfg := flag.String("c", "/etc/halo/cfg.json", "configuration file")
	flag.Parse()
	path := *cfg
	if env := os.Getenv("HALO_TEST_CONFIG"); env != "" {
		path = env
	}
	if _, err := os.Stat(path); err == nil {
		config.ParseConfig(path)
	}
	os.Exit(m.Run())
}

func TestSingleSQL(t *testing.T) {
	p := parser.New()
	stmt, _ := p.ParseOneStmt("alter table t1 drop index idx_c, add unique index `  `(id);", "", "")
	for _, spec := range stmt.(*ast.AlterTableStmt).Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq &&
			c.Tp != ast.ConstraintUniqIndex &&
			c.Tp != ast.ConstraintUniqKey {
			continue
		}
		fmt.Printf("Name:<%d>", len(strings.TrimSpace(c.Name)))
	}
}

// TestTableCreateTableCharsetCollateMatch 建表是校验规则与字符集必须匹配
// RULE:CTB-L2-003
func TestTableCreateTableCharsetCollateMatch(t *testing.T) {
}

// TestTableCreateFromSelectEnabled 是否允许查询语句建表
// RULE:CTB-L2-009
func TestTableCreateFromSelectEnabled(t *testing.T) {
}

// TestTableCreateColumnCharsetCollateMatch 列的字符集与排序规则必须匹配
// RULE:CTB-L2-019
func TestTableCreateColumnCharsetCollateMatch(t *testing.T) {
}

// TestTableCreateTargetDatabaseExists 目标库是否存在
// RULE:CTB-L3-001
func TestTableCreateTargetDatabaseExists(t *testing.T) {
}

// TestTableCreateTargetTableExists 目标表是否存在
// RULE:CTB-L3-002
func TestTableCreateTargetTableExists(t *testing.T) {
}

// TestUpdateTargetDatabaseExists 目标库是否存在
// RULE:UPD-L3-001
func TestUpdateTargetDatabaseExists(t *testing.T) {
}

// TestUpdateTargetTableExists 目标表是否存在
// RULE:UPD-L3-002
func TestUpdateTargetTableExists(t *testing.T) {
}

// TestUpdateTargetColumnExists 目标列是否存在
// RULE:UPD-L3-003
func TestUpdateTargetColumnExists(t *testing.T) {
}

// TestUpdateFilterColumnExists 条件过滤列是否存在
// RULE:UPD-L3-004
func TestUpdateFilterColumnExists(t *testing.T) {
}

// TestUpdateRowsLimit 单次更新允许的最大行数
// RULE:UPD-L3-005
func TestUpdateRowsLimit(t *testing.T) {
}

// TestViewCreateTargetDatabaseExists 目标库是否存在
// RULE:CVW-L3-001
func TestViewCreateTargetDatabaseExists(t *testing.T) {
}

// TestViewCreateTargetViewExists 目标视图是否存在
// RULE:CVW-L3-002
func TestViewCreateTargetViewExists(t *testing.T) {
}

// TestViewAlterTargetDatabaseExists 目标库是否存在
// RULE:VWA-L3-001
func TestViewAlterTargetDatabaseExists(t *testing.T) {
}

// TestViewAlterTargetViewExists 目标视图是否存在
// RULE:VWA-L3-002
func TestViewAlterTargetViewExists(t *testing.T) {
}

// TestViewDropTargetDatabaseExists 目标库是否存在
// RULE:VWD-L3-001
func TestViewDropTargetDatabaseExists(t *testing.T) {
}

// TestViewDropTargetViewExists 目标视图是否存在
// RULE:VWD-L3-002
func TestViewDropTargetViewExists(t *testing.T) {
}

// TestInsertExplicitColumnRequired 是否要求显式列申明
// RULE:INS-L2-001
func TestInsertExplicitColumnRequired(t *testing.T) {
}

// TestInsertUsingSelectEnabled 是否允许INSERT...SELECT
// RULE:INS-L2-002
func TestInsertUsingSelectEnabled(t *testing.T) {
}

// TestInsertMergeRequired 是否合并INSERT
// RULE:INS-L2-003
func TestInsertMergeRequired(t *testing.T) {
}

// TestInsertRowsLimit 单语句允许操作的最大行数
// RULE:INS-L2-004
func TestInsertRowsLimit(t *testing.T) {
}

// TestInsertColumnValueMatch 列类型、值是否匹配
// RULE:INS-L2-005
func TestInsertColumnValueMatch(t *testing.T) {
}

// TestInsertTargetDatabaseExists 目标库是否存在
// RULE:INS-L3-001
func TestInsertTargetDatabaseExists(t *testing.T) {
}

// TestInsertTargetTableExists 目标表是否存在
// RULE:INS-L3-002
func TestInsertTargetTableExists(t *testing.T) {
}

// TestInsertTargetColumnExists 目标列是否存在
// RULE:INS-L3-003
func TestInsertTargetColumnExists(t *testing.T) {
}

// TestInsertValueForNotNullColumnRequired 非空列是否有值
// RULE:INS-L3-004
func TestInsertValueForNotNullColumnRequired(t *testing.T) {
}

// TestIndexCreateTargetDatabaseExists 添加索引的表所属库是否存在
// RULE:CIX-L3-001
func TestIndexCreateTargetDatabaseExists(t *testing.T) {
}

// TestIndexCreateTargetTableExists 条件索引的表是否存在
// RULE:CIX-L3-002
func TestIndexCreateTargetTableExists(t *testing.T) {
}

// TestIndexCreateTargetColumnExists 添加索引的列是否存在
// RULE:CIX-L3-003
func TestIndexCreateTargetColumnExists(t *testing.T) {
}

// TestIndexCreateTargetIndexExists 索引内容是否重复
// RULE:CIX-L3-004
func TestIndexCreateTargetIndexExists(t *testing.T) {
}

// TestIndexCreateTargetNameExists 索引名是否重复
// RULE:CIX-L3-005
func TestIndexCreateTargetNameExists(t *testing.T) {
}

// TestIndexCreateIndexCountLimit 最多能建多少个索引
// RULE:CIX-L3-006
func TestIndexCreateIndexCountLimit(t *testing.T) {
}

// TestIndexCreateIndexBlobColumnEnabled 是否允许在BLOB列上建索引
// RULE:CIX-L3-007
func TestIndexCreateIndexBlobColumnEnabled(t *testing.T) {
}

// TestIndexDropTargetDatabaseExists 目标库是否存在
// RULE:CIX-L3-008
func TestIndexDropTargetDatabaseExists(t *testing.T) {
}

// TestIndexDropTargetTableExists 索引的表是否存在
// RULE:CIX-L3-009
func TestIndexDropTargetTableExists(t *testing.T) {
}

// TestIndexDropTargetIndexExists 删除的索引是否存在
// RULE:CIX-L3-010
func TestIndexDropTargetIndexExists(t *testing.T) {
}

// TestDeleteRowsLimit 单语句允许操作的最大行数
// RULE:DEL-L3-001
func TestDeleteRowsLimit(t *testing.T) {
}

// TestDeleteTargetDatabaseExists 目标库是否存在
// RULE:DEL-L3-002
func TestDeleteTargetDatabaseExists(t *testing.T) {
}

// TestDeleteTargetTableExists 目标表是否存在
// RULE:DEL-L3-003
func TestDeleteTargetTableExists(t *testing.T) {
}

// TestDeleteFilterColumnExists 条件过滤列是否存在
// RULE:DEL-L3-004
func TestDeleteFilterColumnExists(t *testing.T) {
}

// TestDatabaseCreateCharsetCollateMatch 建库的字符集与排序规则必须匹配
// RULE:CDB-L2-001
func TestDatabaseCreateCharsetCollateMatch(t *testing.T) {
}

// TestDatabaseCreateTargetDatabaseExists 目标库是否存在
// RULE:CDB-L2-002
func TestDatabaseCreateTargetDatabaseExists(t *testing.T) {
}

// TestDatabaseAlterAvailableCharsets 可用的字符集
// RULE:MDB-L2-001
func TestDatabaseAlterAvailableCharsets(t *testing.T) {
}

// TestDatabaseAlterAvailableCollates 可用的排序规则
// RULE:MDB-L2-002
func TestDatabaseAlterAvailableCollates(t *testing.T) {
}

// TestDatabaseAlterCharsetCollateMatch 字符集与排序规则必须匹配
// RULE:MDB-L2-003
func TestDatabaseAlterCharsetCollateMatch(t *testing.T) {
}

// TestDatabaseAlterTargetDatabaseExists 目标库是否存在
// RULE:MDB-L3-001
func TestDatabaseAlterTargetDatabaseExists(t *testing.T) {
}

// TestDatabaseDropTargetDatabaseExists 目标库是否存在
// RULE:DDB-L3-001
func TestDatabaseDropTargetDatabaseExists(t *testing.T) {
}

// TestReplaceUsingSelectEnabled 是否允许REPLACE...SELECT
// RULE:RPL-L2-002
func TestReplaceUsingSelectEnabled(t *testing.T) {
}

// TestReplaceMergeRequired 是否合并REPLACE
// RULE:RPL-L2-003
func TestReplaceMergeRequired(t *testing.T) {
}

// TestReplaceRowsLimit 单语句允许操作的最大行数
// RULE:RPL-L2-004
func TestReplaceRowsLimit(t *testing.T) {
}

// TestReplaceColumnValueMatch 列类型、值是否匹配
// RULE:RPL-L2-005
func TestReplaceColumnValueMatch(t *testing.T) {
}

// TestReplaceTargetDatabaseExists 目标库是否存在
// RULE:RPL-L3-001
func TestReplaceTargetDatabaseExists(t *testing.T) {
}

// TestReplaceTargetTableExists 目标表是否存在
// RULE:RPL-L3-002
func TestReplaceTargetTableExists(t *testing.T) {
}

// TestReplaceTargetColumnExists 目标列是否存在
// RULE:RPL-L3-003
func TestReplaceTargetColumnExists(t *testing.T) {
}

// TestReplaceValueForNotNullColumnRequired 非空列是否有值
// RULE:RPL-L3-004
func TestReplaceValueForNotNullColumnRequired(t *testing.T) {
}

// TestKeywordEnabled 是否允许使用关键字
// RULE:MSC-L1-008
func TestKeywordEnabled(t *testing.T) {
}

// TestSplitTicket 是否允许一个工单中包含多个SQL
// RULE:MSC-L0-001
func TestSplitTicket(t *testing.T) {
}
