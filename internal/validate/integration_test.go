package validate_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mia0x75/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/crypto"
	"github.com/mia0x75/halo/internal/platform/db"
	"github.com/mia0x75/halo/internal/validate"
)

// testClusterUUID is the seeded target cluster the rule cases were authored
// against (see scripts/test/seed_meta_test.sql). Without it (or without a live
// metadata database) the cases skip.
const testClusterUUID = "fd577136-0000-0000-0000-000000000001"

// testClusterPassword is the target password the local test containers use
// (mysql/target seed). The metadata seed only stores a placeholder, so the
// tests encrypt this value at runtime with the active config key.
const testClusterPassword = "test"

// withRealPassword replaces the seeded placeholder password with an encrypted
// copy of testClusterPassword, so the executor can authenticate to the target.
func withRealPassword(t *testing.T, cluster *models.Cluster) {
	t.Helper()
	enc, err := crypto.EncryptAES([]byte(testClusterPassword), config.Config().Secret.Crypto)
	require.NoError(t, err)
	cluster.Password = enc
}

// ensureFixture creates the schema-aware fixtures the table/index cases read.
// The cases assume a pre-existing starwars users/orders pair (target seed); the
// extra tables are created here (not seeded, per the test-data contract) and
// dropped on cleanup.
func ensureFixture(t *testing.T, a *app.App) {
	t.Helper()
	cluster, err := a.Repos.Cluster.GetByUUID(context.Background(), testClusterUUID)
	if errors.Is(err, domain.ErrNotFound) {
		t.Skip("test cluster not seeded; run: make test-seed-meta")
	}
	require.NoError(t, err)
	cluster.Database = "starwars"

	fixtures := []struct{ create, drop string }{
		{
			create: "CREATE TABLE IF NOT EXISTS starwars.t1 (id INT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='t1'",
			drop:   "DROP TABLE IF EXISTS starwars.t1",
		},
		{
			// TEXT payload exercises the blob/text index guard (CIX-L3-007).
			create: "CREATE TABLE IF NOT EXISTS starwars.large_t (id INT, payload TEXT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='large_t'",
			drop:   "DROP TABLE IF EXISTS starwars.large_t",
		},
		{
			// Five indexes trip the per-table index cap (CIX-L3-006, limit 5).
			create: "CREATE TABLE IF NOT EXISTS starwars.idx_many (id INT, KEY index_1(id), KEY index_2(id), KEY index_3(id), KEY index_4(id), KEY index_5(id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='idx_many'",
			drop:   "DROP TABLE IF EXISTS starwars.idx_many",
		},
		// User objects exercise the existence rules (CVW/DVW/CFU/MFU/DFU/CTG/
		// DTG/CEV/MEV/DEV/CSP/MSP/DSP). Routines use a single-statement body so
		// the executor's statement splitter does not need a client delimiter.
		{
			create: "CREATE OR REPLACE VIEW starwars.halo_v AS SELECT id FROM starwars.t1",
			drop:   "DROP VIEW IF EXISTS starwars.halo_v",
		},
		{
			create: "CREATE FUNCTION starwars.halo_fn() RETURNS INT DETERMINISTIC RETURN 1",
			drop:   "DROP FUNCTION IF EXISTS starwars.halo_fn",
		},
		{
			create: "CREATE PROCEDURE starwars.halo_proc() SELECT 1",
			drop:   "DROP PROCEDURE IF EXISTS starwars.halo_proc",
		},
		{
			create: "CREATE TRIGGER starwars.halo_trg BEFORE INSERT ON starwars.t1 FOR EACH ROW SET NEW.id = NEW.id",
			drop:   "DROP TRIGGER IF EXISTS starwars.halo_trg",
		},
		{
			create: "CREATE EVENT starwars.halo_ev ON SCHEDULE EVERY 1 DAY DO SELECT 1",
			drop:   "DROP EVENT IF EXISTS starwars.halo_ev",
		},
	}

	exec := agent.LocalExecutor{}
	for _, f := range fixtures {
		_, err = exec.Execute(context.Background(), cluster, testClusterPassword, f.create)
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		for _, f := range fixtures {
			_, _ = exec.Execute(context.Background(), cluster, testClusterPassword, f.drop)
		}
	})
}

// replaceManyValues mirrors insertManyValues for the REPLACE row-count rule.
func replaceManyValues(n int) string {
	var b strings.Builder
	b.WriteString("replace into t1(id) values")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("(1)")
	}
	b.WriteString(";")
	return b.String()
}

type RuleTestCase struct {
	RuleName string
	Text     string
	Valid    bool
}

// buildApp boots the production dependency container, which also loads the
// caches, and skips when no test database is configured or seeded.
func buildApp(t *testing.T) *app.App {
	t.Helper()
	cfg := config.Config()
	if cfg == nil || cfg.Database.Addr == "" {
		t.Skip("no test config: set HALO_TEST_CONFIG to a config file")
	}
	a, err := app.New(context.Background(), db.Config{
		Addr:           cfg.Database.Addr,
		MaxIdle:        cfg.Database.MaxIdle,
		MaxConnections: cfg.Database.MaxConnections,
		WaitTimeout:    cfg.Database.WaitTimeout,
		Driver:         cfg.Database.Driver,
	})
	if err != nil {
		t.Skipf("metadata database unavailable: %v", err)
	}
	t.Cleanup(func() { _ = a.Close() })
	if cache.ClustersMap.Count() == 0 {
		t.Skip("metadata database has no clusters seeded")
	}
	return a
}

// buildContext mirrors the engine's Run setup: it resolves the seeded target
// cluster and loads its metadata so the schema aware rules can run.
func buildContext(t *testing.T) *validate.Context {
	t.Helper()
	cluster := cache.ClustersMap.Any(func(elem *models.Cluster) bool {
		return elem.UUID == testClusterUUID
	})
	if cluster == nil {
		t.Skipf("test cluster %s not found; run: make test-seed-meta", testClusterUUID)
	}
	withRealPassword(t, cluster)
	ticket := &models.Ticket{
		UUID:       "9f25e9ee-3ca3-4fe6-9747-39ea90be59b1",
		UserID:     1,
		ReviewerID: 1,
		Database:   "starwars",
	}
	ctx := &validate.Context{
		Cluster: bridge.ClusterToDomain(cluster),
		Ticket:  bridge.TicketToDomain(ticket),
	}
	passwd := func(c *models.Cluster) string {
		bs, _ := crypto.DecryptAES(c.Password, config.Config().Secret.Crypto)
		return string(bs)
	}
	exec := agent.LocalExecutor{}
	target := ctx.Cluster
	if tables, err := exec.LoadTables(context.Background(), target, passwd(cluster)); err == nil {
		ctx.Tables = tables
	}
	if databases, err := exec.Databases(context.Background(), target, passwd(cluster)); err == nil {
		ctx.Databases = make([]models.Database, 0, len(databases))
		for _, d := range databases {
			ctx.Databases = append(ctx.Databases, models.Database{Name: d.Name, Charset: d.Charset, Collate: d.Collate})
		}
	}
	if objects, err := exec.LoadObjects(context.Background(), target, passwd(cluster)); err == nil {
		ctx.Objects = objects
	}
	return ctx
}

// runCase drives a single rule through the public engine. Only the addressed
// rule is injected, so the result is equivalent to calling that rule alone.
func runCase(t *testing.T, ctx *validate.Context, p *parser.Parser, c RuleTestCase) (bool, bool) {
	t.Helper()
	mr := cache.RulesMap.Any(func(elem *models.Rule) bool {
		return elem.Name == c.RuleName
	})
	if mr == nil {
		t.Errorf("rule %s not found", c.RuleName)
		return false, false
	}
	stmt, err := p.ParseOneStmt(c.Text, "", "")
	if err != nil {
		t.Errorf("语法错误: %s: %s", c.RuleName, err.Error())
		return false, false
	}
	s := &validate.Subject{Content: c.Text, Violations: &models.Violations{}, StmtNode: stmt}
	ctx.Subjects = []*validate.Subject{s}
	rule := bridge.RuleToDomain(mr)
	// The seed ships some implemented rules disabled (even Bitwise); force the
	// rule on so a case can exercise its handler regardless of seed state.
	rule.Bitwise |= 1
	ctx.Rules = []*validate.Rule{{Rule: *rule}}
	validate.RunRegistry(ctx)
	return s.Violations.Marshal() == "", true
}

// ruleCaseTables gathers every authored case table so a rule can be tested by
// name without duplicating its cases.
var ruleCaseTables = [][]RuleTestCase{
	tableCases,
	viewCases,
	updateCasess,
	replaceCases,
	insertCases,
	indexCases,
	deleteCases,
	databaseCases,
	miscCases,
}

// casesForRule returns every case authored for a single rule.
func casesForRule(name string) []RuleTestCase {
	var out []RuleTestCase
	for _, table := range ruleCaseTables {
		for _, c := range table {
			if c.RuleName == name {
				out = append(out, c)
			}
		}
	}
	return out
}

// runRuleCases boots the live fixture and drives every case authored for name.
func runRuleCases(t *testing.T, name string) {
	t.Helper()
	cases := casesForRule(name)
	if len(cases) == 0 {
		t.Skipf("no test cases authored for %s", name)
	}
	a := buildApp(t)
	ensureFixture(t, a)
	ctx := buildContext(t)
	p := parser.New()
	for _, c := range cases {
		valid, ok := runCase(t, ctx, p, c)
		if !ok {
			continue
		}
		if c.Valid != valid {
			t.Errorf("%s: %s: want valid=%v, got %v", name, c.Text, c.Valid, valid)
		}
	}
}

func TestRun_Integration(t *testing.T) {
	a := buildApp(t)

	cluster, err := a.Repos.Cluster.GetByUUID(context.Background(), testClusterUUID)
	if errors.Is(err, domain.ErrNotFound) {
		t.Skip("test cluster not seeded; run: make test-seed-meta")
	}
	require.NoError(t, err)
	enc, err := crypto.EncryptAES([]byte(testClusterPassword), config.Config().Secret.Crypto)
	require.NoError(t, err)
	cluster.Password = enc

	stmts := []*models.Statement{
		{Sequence: 1, UUID: "8b0e2f31-0000-4000-8000-000000000001", Content: "SELECT 1;", Type: 1, Status: 1, Violations: &models.Violations{}},
	}
	validate.Run(stmts, bridge.ClusterFromDomain(cluster), &models.Ticket{Database: "starwars"})
	require.NotNil(t, stmts[0].Violations)
}

var tableCases = []RuleTestCase{
	{"CTB-L2-001", "create table t1 (id INT);", false},                                                 // TableCreateAvailableCharsets - 没有指定字符集
	{"CTB-L2-001", "create table t1 (id INT) CHARSET = 'utf8';", false},                                // TableCreateAvailableCharsets - 不允许
	{"CTB-L2-001", "create table t1 (id INT) CHARSET = 'UTF8';", false},                                // TableCreateAvailableCharsets - 不允许（大写）
	{"CTB-L2-001", "create table t1 (id INT) CHARSET = 'utf8mb4';", true},                              // TableCreateAvailableCharsets - 允许
	{"CTB-L2-001", "create table t1 (id INT) CHARSET = 'UTF8MB4';", true},                              // TableCreateAvailableCharsets - 允许（大写）
	{"CTB-L2-001", "create table t2 as select * from t1;", true},                                       //
	{"CTB-L2-002", "create table t1 (id INT);", false},                                                 // TableCreateAvailableCollates - 未指定
	{"CTB-L2-002", "create table t1 (id INT) collate = 'utf8mb4_general_ci';", true},                   // TableCreateAvailableCollates - 允许
	{"CTB-L2-002", "create table t1 (id INT) collate = 'UTF8MB4_GENERAL_CI';", true},                   // TableCreateAvailableCollates - 允许（大写）
	{"CTB-L2-002", "create table t1 (id INT) default collate = 'utf8mb4_bin';", true},                  // TableCreateAvailableCollates - 允许
	{"CTB-L2-002", "create table t1 (id INT) collate = 'latin1_general_ci';", false},                   // TableCreateAvailableCollates - 不允许
	{"CTB-L2-002", "create table t2 as select * from t1;", true},                                       //
	{"CTB-L2-003", "create table t1 (id INT) CHARSET = 'utf8mb4' collate = 'utf8_bin';", false},        // 字符集和排序规则匹配检查 - 不匹配
	{"CTB-L2-003", "create table t1 (id INT) CHARSET = 'utf8' collate = 'utf8mb4_unicode_ci';", false}, // 字符集和排序规则匹配检查 - 不匹配
	{"CTB-L2-003", "create table t1 (id INT) CHARSET = 'utf8mb4' collate = 'utf8mb4_bin';", true},      // 字符集和排序规则匹配检查 - 匹配
	{"CTB-L2-003", "create table t2 as select * from t1;", true},                                       //
	{"CTB-L2-003", "CREATE TABLE t2 (c1 VARCHAR(2)) CHARSET 'utf8mb4';", true},
	{"CTB-L2-003", "CREATE TABLE t2 (c1 VARCHAR(2)) COLLATE 'utf8mb4_bin';", true},
	{"CTB-L2-004", "create table t1 (id INT);", false},                                                                         // TableCreateAvailableEngines
	{"CTB-L2-004", "create table t1 (id INT) engine = 'innodb';", true},                                                        // TableCreateAvailableEngines
	{"CTB-L2-004", "create table t1 (id INT) ENGINE 'INNODB';", true},                                                          // TableCreateAvailableEngines
	{"CTB-L2-004", "CREATE TABLE t1 (id INT) CHARSET='latin1' ENGINE 'innodb';", true},                                         // TableCreateAvailableEngines
	{"CTB-L2-004", "create table t1 (id INT) engine = 'csv';", false},                                                          // TableCreateAvailableEngines
	{"CTB-L2-004", "create table t1 (id INT) engine = 'rocksdb';", true},                                                       // TableCreateAvailableEngines
	{"CTB-L2-004", "create table t2 as select * from t1;", true},                                                               //
	{"CTB-L2-005", "create table t1 (id INT);", true},                                                                          // TableCreateTableNameQualified
	{"CTB-L2-005", "create table `t 1` (id INT);", false},                                                                      // TableCreateTableNameQualified
	{"CTB-L2-005", "create table T1 (id INT);", true},                                                                          // TableCreateTableNameQualified
	{"CTB-L2-005", "create table t_1 (id INT);", true},                                                                         // TableCreateTableNameQualified
	{"CTB-L2-005", "create table `t=1` (id INT);", false},                                                                      // TableCreateTableNameQualified
	{"CTB-L2-006", "create table t1 (id INT);", true},                                                                          // TableCreateTableNameLowerCaseRequired
	{"CTB-L2-006", "create table T1 (id INT);", false},                                                                         // TableCreateTableNameLowerCaseRequired
	{"CTB-L2-007", "create table T1 (id INT);", true},                                                                          // TableCreateTableNameMaxLength
	{"CTB-L2-007", "create table T1234567890123456789 (id INT);", true},                                                        // TableCreateTableNameMaxLength
	{"CTB-L2-007", "create table T12345678901234567890 (id INT);", false},                                                      // TableCreateTableNameMaxLength
	{"CTB-L2-007", "create table t2 as select * from t1;", true},                                                               //
	{"CTB-L2-008", "create table t1 (id INT);", false},                                                                         // TableCreateTableCommentRequired
	{"CTB-L2-008", "create table t1 (id INT) COMMENT = '';", false},                                                            // TableCreateTableCommentRequired
	{"CTB-L2-008", "create table t1 (id INT) comment = 't1';", true},                                                           // TableCreateTableCommentRequired
	{"CTB-L2-008", "create table t1 (id INT) comment 't1';", true},                                                             // TableCreateTableCommentRequired
	{"CTB-L2-008", "create table t2 as select * from t1;", true},                                                               //
	{"CTB-L2-009", "create table t2 as select * from t1;", false},                                                              //
	{"CTB-L2-009", "create table t1 (id INT);", true},                                                                          // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (id INT);", true},                                                                          // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (Id INT);", true},                                                                          // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (ID INT);", true},                                                                          // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (ID_123 INT);", true},                                                                      // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (_Id INT);", false},                                                                        // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (_9 INT);", false},                                                                         // TableCreateColumnNameQualified
	{"CTB-L2-010", "create table t1 (`id 1` INT);", false},                                                                     // TableCreateColumnNameQualified
	{"CTB-L2-010", "CREATE TABLE t2 (` ` VARCHAR(2));", false},                                                                 //
	{"CTB-L2-011", "create table t1 (id INT);", true},                                                                          // TableCreateColumnNameLowerCaseRequired
	{"CTB-L2-011", "create table t_1 (id INT);", true},                                                                         // TableCreateColumnNameLowerCaseRequired
	{"CTB-L2-011", "create table t1 (Id INT);", false},                                                                         // TableCreateColumnNameLowerCaseRequired
	{"CTB-L2-011", "CREATE TABLE t2 (` ` VARCHAR(2));", true},                                                                  //
	{"CTB-L2-012", "create table t1 (Id_12345678901234567 INT);", true},                                                        // TableCreateColumnNameMaxLength
	{"CTB-L2-012", "create table t1 (Id INT);", true},                                                                          // TableCreateColumnNameMaxLength
	{"CTB-L2-012", "create table t1 (Id_123456789012345678 INT);", false},                                                      // TableCreateColumnNameMaxLength
	{"CTB-L2-013", "create table t1 (Id INT);", true},                                                                          // TableCreateColumnNameDuplicate
	{"CTB-L2-013", "create table t1 (Id INT, ID TINYINT);", false},                                                             // TableCreateColumnNameDuplicate
	{"CTB-L2-013", "create table t1 (Id INT, Id TINYINT);", false},                                                             // TableCreateColumnNameDuplicate
	{"CTB-L2-015", "create table t1 (Id INT);", true},                                                                          // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id TINYINT);", true},                                                                      // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id SMALLINT);", true},                                                                     // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id MEDIUMINT);", true},                                                                    // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id BIGINT);", true},                                                                       // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id VARCHAR(1));", true},                                                                   // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id CHAR(1));", true},                                                                      // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id DATE);", true},                                                                         // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id TIME);", true},                                                                         // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id DATETIME);", true},                                                                     // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id YEAR);", true},                                                                         // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id DECIMAL(10,1));", true},                                                                // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id FLOAT);", false},                                                                       // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id FLOAT(3,1));", false},                                                                  // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id DOUBLE);", false},                                                                      // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id DOUBLE(10,1));", false},                                                                // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id REAL);", false},                                                                        // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id BINARY(10));", true},                                                                   // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id VARBINARY(10));", true},                                                                // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id TINYBLOB);", true},                                                                     // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id BLOB);", true},                                                                         // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id MEDIUMBLOB);", true},                                                                   // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id LONGBLOB);", true},                                                                     // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id TINYTEXT);", true},                                                                     // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id TEXT);", true},                                                                         // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id MEDIUMTEXT);", true},                                                                   // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id LONGTEXT);", true},                                                                     // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id TIMESTAMP);", true},                                                                    // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id BIT);", false},                                                                         // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id SET('a','b'));", false},                                                                // TableCreateColumnUnwantedTypes
	{"CTB-L2-015", "create table t1 (Id ENUM('a','b'));", false},                                                               // TableCreateColumnUnwantedTypes
	{"CTB-L2-016", "create table t1 (Id INT);", false},                                                                         // TableCreateColumnCommentRequired
	{"CTB-L2-016", "create table t1 (Id INT COMMENT '');", false},                                                              // TableCreateColumnCommentRequired
	{"CTB-L2-016", "create table t1 (Id INT COMMENT 'Id');", true},                                                             // TableCreateColumnCommentRequired
	{"CTB-L2-017", "create table t1 (Id VARCHAR(1));", true},                                                                   // TableCreateColumnAvailableCharsets
	{"CTB-L2-017", "create table t1 (Id VARCHAR(1) CHARSET 'utf8');", false},                                                   // TableCreateColumnAvailableCharsets
	{"CTB-L2-017", "create table t1 (Id VARCHAR(1) CHARSET 'utf8mb4');", true},                                                 // TableCreateColumnAvailableCharsets
	{"CTB-L2-018", "create table t1 (Id VARCHAR(1));", true},                                                                   // TableCreateColumnAvailableCollates
	{"CTB-L2-018", "create table t1 (Id VARCHAR(1) COLLATE 'utf8_general_ci');", false},                                        // TableCreateColumnAvailableCollates
	{"CTB-L2-018", "create table t1 (Id VARCHAR(1) COLLATE 'utf8mb4_bin');", true},                                             // TableCreateColumnAvailableCollates
	{"CTB-L2-019", "create table t1 (Id VARCHAR(1) COLLATE 'utf8mb4_bin');", true},                                             //
	{"CTB-L2-019", "create table t1 (Id VARCHAR(1) CHARSET 'latin1');", true},                                                  //
	{"CTB-L2-019", "create table t1 (Id VARCHAR(1) CHARSET 'utf8mb4' COLLATE 'utf8mb4_bin');", true},                           //
	{"CTB-L2-019", "create table t1 (Id VARCHAR(1) CHARSET 'latin1' COLLATE 'latin1_general_ci');", true},                      //
	{"CTB-L2-019", "create table t1 (Id VARCHAR(1) CHARSET 'utf8mb4' COLLATE 'latin1_general_ci');", false},                    //
	{"CTB-L2-020", "create table t1 (Id INT);", true},                                                                          // TableCreateColumnNotNullWithDefaultRequired
	{"CTB-L2-020", "create table t1 (Id INT NOT NULL);", false},                                                                // TableCreateColumnNotNullWithDefaultRequired
	{"CTB-L2-020", "create table t1 (Id INT AUTO_INCREMENT NOT NULL);", true},                                                  // TableCreateColumnNotNullWithDefaultRequired
	{"CTB-L2-020", "create table t1 (Id INT NOT NULL DEFAULT 1);", true},                                                       // TableCreateColumnNotNullWithDefaultRequired
	{"CTB-L2-021", "create table t1 (Id TINYINT AUTO_INCREMENT NOT NULL);", false},                                             // TableCreateColumnAutoIncAvailableTypes
	{"CTB-L2-021", "create table t1 (Id SMALLINT AUTO_INCREMENT NOT NULL);", false},                                            // TableCreateColumnAutoIncAvailableTypes
	{"CTB-L2-021", "create table t1 (Id INT AUTO_INCREMENT NOT NULL);", true},                                                  // TableCreateColumnAutoIncAvailableTypes
	{"CTB-L2-021", "create table t1 (Id BIGINT AUTO_INCREMENT NOT NULL);", true},                                               // TableCreateColumnAutoIncAvailableTypes
	{"CTB-L2-022", "create table t1 (Id TINYINT AUTO_INCREMENT NOT NULL);", false},                                             // TableCreateColumnAutoIncIsUnsigned
	{"CTB-L2-022", "create table t1 (Id INT AUTO_INCREMENT NOT NULL);", false},                                                 // TableCreateColumnAutoIncIsUnsigned
	{"CTB-L2-022", "create table t1 (Id INT UNSIGNED AUTO_INCREMENT NOT NULL);", true},                                         // TableCreateColumnAutoIncIsUnsigned
	{"CTB-L2-023", "create table t1 (Id INT AUTO_INCREMENT NOT NULL);", false},                                                 // TableCreateColumnAutoIncMustPrimaryKey
	{"CTB-L2-023", "create table t1 (Id INT AUTO_INCREMENT PRIMARY KEY);", true},                                               // TableCreateColumnAutoIncMustPrimaryKey
	{"CTB-L2-023", "create table t1 (Id INT AUTO_INCREMENT, PRIMARY KEY(Id));", true},                                          // TableCreateColumnAutoIncMustPrimaryKey
	{"CTB-L2-024", "create table t1 (ts1 TIMESTAMP);", true},                                                                   // TableCreateTimestampColumnCountLimit
	{"CTB-L2-024", "create table t1 (ts1 TIMESTAMP, ts2 TIMESTAMP);", false},                                                   // TableCreateTimestampColumnCountLimit
	{"CTB-L2-025", "create table t1 (a INT, b INT, c INT, KEY i1(a,b,c));", true},                                              // TableCreateIndexMaxColumnLimit
	{"CTB-L2-025", "create table t1 (a INT, b INT, c INT, d INT, KEY i1(a,b,c,d));", false},                                    // TableCreateIndexMaxColumnLimit
	{"CTB-L2-026", "create table t1 (Id INT );", false},                                                                        // TableCreatePrimaryKeyRequired
	{"CTB-L2-026", "create table t1 (Id INT NOT NULL PRIMARY KEY);", true},                                                     // TableCreatePrimaryKeyRequired
	{"CTB-L2-026", "create table t1 (Id INT, PRIMARY KEY(Id));", true},                                                         // TableCreatePrimaryKeyRequired
	{"CTB-L2-033", "create table t1 (Id INT, KEY idx_id (Id));", true},                                                         // TableCreateIndexNameQualified
	{"CTB-L2-033", "create table t1 (Id INT, KEY idx_Id (Id));", true},                                                         // TableCreateIndexNameQualified
	{"CTB-L2-033", "create table t1 (Id INT, KEY `idx id` (Id));", false},                                                      // TableCreateIndexNameQualified
	{"CTB-L2-033", "create table t1 (Id INT, KEY _idx_id (Id));", false},                                                       // TableCreateIndexNameQualified
	{"CTB-L2-033", "CREATE TABLE t2 (c1 INT, KEY ` ` (c1));", true},                                                            //
	{"CTB-L2-034", "create table t1 (Id INT, KEY idx_id (Id));", true},                                                         // TableCreateIndexNameLowerCaseRequired
	{"CTB-L2-034", "create table t1 (Id INT, KEY idx_1 (Id));", true},                                                          // TableCreateIndexNameLowerCaseRequired
	{"CTB-L2-034", "create table t1 (Id INT, KEY _idx_id (Id));", true},                                                        // TableCreateIndexNameLowerCaseRequired
	{"CTB-L2-034", "create table t1 (Id INT, KEY idx_Id (Id));", false},                                                        // TableCreateIndexNameLowerCaseRequired
	{"CTB-L2-034", "CREATE TABLE t2 (c1 INT, KEY ` ` (c1));", true},                                                            //
	{"CTB-L2-035", "create table t1 (Id INT, KEY index_1 (Id));", true},                                                        // TableCreateIndexNameMaxLength
	{"CTB-L2-035", "create table t1 (Id INT, KEY index_1234 (Id));", true},                                                     // TableCreateIndexNameMaxLength
	{"CTB-L2-035", "create table t1 (Id INT, KEY index_12345 (Id));", false},                                                   // TableCreateIndexNameMaxLength
	{"CTB-L2-036", "create table t1 (Id INT, KEY index_1 (Id));", true},                                                        // TableCreateIndexNamePrefixRequired
	{"CTB-L2-036", "create table t1 (Id INT, KEY idx_1 (Id));", false},                                                         // TableCreateIndexNamePrefixRequired
	{"CTB-L2-036", "CREATE TABLE t2 (c1 INT, KEY ` ` (c1));", true},                                                            //
	{"CTB-L2-032", "create table t1 (Id INT, KEY(Id));", false},                                                                // TableCreateIndexNameExplicit
	{"CTB-L2-032", "create table t1 (Id INT, KEY idx_id (Id));", true},                                                         // TableCreateIndexNameExplicit
	{"CTB-L2-037", "create table t1 (Id INT, UNIQUE(Id));", false},                                                             // TableCreateUniqueNameExplicit
	{"CTB-L2-037", "create table t1 (Id INT, UNIQUE unq_id (Id));", true},                                                      // TableCreateUniqueNameExplicit
	{"CTB-L2-038", "create table t1 (Id INT, UNIQUE unq_id (Id));", true},                                                      // TableCreateUniqueNameQualified
	{"CTB-L2-038", "create table t1 (Id INT, UNIQUE unq_Id (Id));", true},                                                      // TableCreateUniqueNameQualified
	{"CTB-L2-038", "create table t1 (Id INT, UNIQUE `unq id` (Id));", false},                                                   // TableCreateUniqueNameQualified
	{"CTB-L2-038", "create table t1 (Id INT, UNIQUE _unq_id (Id));", false},                                                    // TableCreateUniqueNameQualified
	{"CTB-L2-038", "CREATE TABLE t2 (c1 INT, UNIQUE KEY ` ` (c1));", true},                                                     //
	{"CTB-L2-039", "create table t1 (Id INT, UNIQUE unq_id (Id));", true},                                                      // TableCreateUniqueNameLowerCaseRequired
	{"CTB-L2-039", "create table t1 (Id INT, UNIQUE unq_1 (Id));", true},                                                       // TableCreateUniqueNameLowerCaseRequired
	{"CTB-L2-039", "create table t1 (Id INT, UNIQUE _unq_id (Id));", true},                                                     // TableCreateUniqueNameLowerCaseRequired
	{"CTB-L2-039", "create table t1 (Id INT, UNIQUE unq_Id (Id));", false},                                                     // TableCreateUniqueNameLowerCaseRequired
	{"CTB-L2-039", "CREATE TABLE t2 (c1 INT, UNIQUE KEY ` ` (c1));", true},                                                     //
	{"CTB-L2-040", "create table t1 (Id INT, UNIQUE unique_1 (Id));", true},                                                    // TableCreateUniqueNameMaxLength
	{"CTB-L2-040", "create table t1 (Id INT, UNIQUE unique_123 (Id));", true},                                                  // TableCreateUniqueNameMaxLength
	{"CTB-L2-040", "create table t1 (Id INT, UNIQUE unique_1234 (Id));", false},                                                // TableCreateUniqueNameMaxLength
	{"CTB-L2-041", "create table t1 (Id INT, UNIQUE unique_1 (Id));", true},                                                    // TableCreateUniqueNamePrefixRequired
	{"CTB-L2-041", "create table t1 (Id INT, UNIQUE unq_1 (Id));", false},                                                      // TableCreateUniqueNamePrefixRequired
	{"CTB-L2-041", "CREATE TABLE t2 (c1 INT, UNIQUE KEY ` ` (c1));", true},                                                     //
	{"CTB-L2-042", "create table t1 (Id INT, FOREIGN KEY `fk_a` (`id`) REFERENCES `t2` (`id`));", false},                       // TableCreateForeignKeyEnabled
	{"CTB-L2-043", "create table t1 (Id INT, FOREIGN KEY (`id`) REFERENCES `t2` (`id`));", false},                              // TableCreateForeignKeyNameExplicit
	{"CTB-L2-043", "create table t1 (Id INT, FOREIGN KEY `fk_a` (`id`) REFERENCES `t2` (`id`));", true},                        // TableCreateForeignKeyNameExplicit
	{"CTB-L2-044", "create table t1 (Id INT, FOREIGN KEY `fk_a` (`id`) REFERENCES `t2` (`id`));", true},                        // TableCreateForeignKeyNameQualified
	{"CTB-L2-044", "create table t1 (Id INT, FOREIGN KEY `fk_1` (`id`) REFERENCES `t2` (`id`));", true},                        // TableCreateForeignKeyNameQualified
	{"CTB-L2-044", "create table t1 (Id INT, FOREIGN KEY `fk D` (`id`) REFERENCES `t2` (`id`));", false},                       // TableCreateForeignKeyNameQualified
	{"CTB-L2-045", "create table t1 (Id INT, FOREIGN KEY `fk_a` (`id`) REFERENCES `t2` (`id`));", true},                        // TableCreateForeignKeyNameLowerCaseRequired
	{"CTB-L2-045", "create table t1 (Id INT, FOREIGN KEY `fk_A` (`id`) REFERENCES `t2` (`id`));", false},                       // TableCreateForeignKeyNameLowerCaseRequired
	{"CTB-L2-046", "create table t1 (Id INT, FOREIGN KEY `fk_a` (`id`) REFERENCES `t2` (`id`));", true},                        // TableCreateForeignKeyNameMaxLength
	{"CTB-L2-046", "create table t1 (Id INT, FOREIGN KEY `fk_ab12345678901234567890` (`id`) REFERENCES `t2` (`id`));", true},   // TableCreateForeignKeyNameMaxLength
	{"CTB-L2-046", "create table t1 (Id INT, FOREIGN KEY `fk_ab123456789012345678901` (`id`) REFERENCES `t2` (`id`));", false}, // TableCreateForeignKeyNameMaxLength
	{"CTB-L2-047", "create table t1 (Id INT, FOREIGN KEY `fk_a` (`id`) REFERENCES `t2` (`id`));", true},                        // TableCreateForeignKeyNamePrefixRequired
	{"CTB-L2-047", "create table t1 (Id INT, FOREIGN KEY `k_a` (`id`) REFERENCES `t2` (`id`));", false},                        // TableCreateForeignKeyNamePrefixRequired
	{"CTB-L2-048", "create table t1 (id INT,KEY i1(id),KEY i2(id),KEY i3(id),KEY i4(id),KEY i5(id));", true},                   // TableCreateIndexCountLimit
	{"CTB-L2-048", "create table t1 (id INT PRIMARY KEY,KEY i1(id),KEY i2(id),KEY i3(id),KEY i4(id));", true},                  // TableCreateIndexCountLimit
	{"CTB-L2-048", "create table t1 (id INT,KEY i1(id),KEY i2(id),KEY i3(id),KEY i4(id),KEY i5(id),KEY i6(id));", false},       // TableCreateIndexCountLimit
	{"CTB-L2-048", "CREATE TABLE mock_t1 (`id` INT AUTO_INCREMENT PRIMARY KEY,`k` INT,`l` INT,`m` INT,`n` INT,`o` INT,`p` INT,`q` INT,`r` INT,`s` INT,`t` INT UNIQUE KEY,`u` INT COMMENT 'u' UNIQUE KEY,`v` INT COMMENT 'v' UNIQUE KEY,`w` INT COMMENT 'w' UNIQUE KEY,`x` INT COMMENT 'x' UNIQUE KEY);", false}, // TableCreateIndexCountLimit
	{"CTB-L2-049", "CREATE TABLE mock_t2 LIKE t1;", false},
	{"CTB-L2-049", "CREATE TABLE t2 (c1 INT);", true},
	{"CTB-L2-050", "create table t1 (id INT AUTO_INCREMENT, id2 INT AUTO_INCREMENT);", false},       //
	{"CTB-L2-050", "create table t1 (id INT AUTO_INCREMENT);", true},                                //
	{"CTB-L2-050", "create table t1 (id INT);", true},                                               //
	{"CTB-L2-051", "create table t1 (id INT);", true},                                               //
	{"CTB-L2-051", "create table t1 (id INT PRIMARY KEY);", true},                                   //
	{"CTB-L2-051", "create table t1 (id INT,PRIMARY KEY (id));", true},                              //
	{"CTB-L2-051", "create table t1 (id INT PRIMARY KEY, id2 INT PRIMARY KEY);", false},             //
	{"CTB-L2-051", "create table t1 (id INT PRIMARY KEY, id2 INT, PRIMARY KEY (id2));", false},      //
	{"CTB-L2-051", "create table t1 (id INT,id2 INT, PRIMARY KEY (id), PRIMARY KEY (id2));", false}, //
	{"CTB-L3-001", "CREATE TABLE mock.t1(id INT);", false},
	{"CTB-L3-001", "CREATE TABLE starwars.t1(id INT);", true},
	{"CTB-L3-001", "CREATE TABLE t1(id INT);", true},
	{"CTB-L3-002", "CREATE TABLE t1(id INT);", false},
	{"CTB-L3-002", "CREATE TABLE t2(id INT);", true},

	{"MTB-L2-001", "alter table t1 CONVERT TO CHARACTER SET 'utf8mb4';", true},                                //
	{"MTB-L2-001", "alter table t1 CONVERT TO CHARACTER SET 'latin1';", false},                                //
	{"MTB-L2-001", "ALTER table db.blacklist_model_audit_snapshot CONVERT TO CHARACTER SET 'utf8mb4';", true}, //
	{"MTB-L2-002", "alter table t1 convert to character set utf8 collate utf8_unicode_ci;", false},            //
	{"MTB-L2-002", "alter table t1 convert to character set utf8 collate utf8mb4_general_ci;", true},          //
	{"MTB-L2-003", "alter table t1 COLLATE 'utf8mb4_bin';", true},                                             //
	{"MTB-L2-003", "alter table t1 CHARSET 'latin1';", false},                                                 // 原表的排序规则utf8mb4
	{"MTB-L2-003", "alter table t1 CHARSET 'utf8mb4' COLLATE 'utf8mb4_bin';", true},                           //
	{"MTB-L2-003", "alter table t1 CHARSET 'latin1' COLLATE 'latin1_general_ci';", true},                      //
	{"MTB-L2-003", "alter table t1 CHARSET 'utf8mb4' COLLATE 'latin1_general_ci';", false},                    //
	{"MTB-L2-004", "alter table t1 ENGINE InnoDB;", true},                                                     //
	{"MTB-L2-004", "alter table t1 ENGINE MyISAM;", false},                                                    //
	{"MTB-L2-005", "alter table t1 ADD c1 INT;", true},
	{"MTB-L2-005", "alter table t1 ADD `$% c1` INT, change c2 `MHzzz#ssz` int ;", false},
	{"MTB-L2-005", "alter table t1 ADD c2 INT, change c2  __xxx int ;", false},
	{"MTB-L2-005", "alter table t1 modify c2 INT, modify  __xxx int ;", true},
	{"MTB-L2-005", "alter table t1 ADD C1 INT COMMENT '';", true},                      //
	{"MTB-L2-005", "alter table t1 ADD `c 1` INT COMMENT 'c1';", false},                //
	{"MTB-L2-006", "alter table t1 ADD c1 INT;", true},                                 //
	{"MTB-L2-006", "alter table t1 ADD C1 INT COMMENT '';", false},                     //
	{"MTB-L2-007", "alter table t1 ADD c1012345678901234567 INT;", true},               // 长度 = 20
	{"MTB-L2-007", "alter table t1 ADD c10123456789012345678 INT COMMENT '';", false},  // 长度 > 20
	{"MTB-L2-007", "alter table t1 modify c10123456789012345678 VARCHAR(10);", true},   // 原有字段，不关心
	{"MTB-L2-007", "alter table t1 modify c1 VARCHAR(10) CHARACTER SET latin1;", true}, // 原有字段，不关心
	{"MTB-L2-007", "alter table t1 change c2 c2012345678901234567 char(40);", true},    // 长度 = 20
	{"MTB-L2-007", "alter table t1 change c2 c20123456789012345678 char(40);", false},  // 长度 > 20
	{"MTB-L2-008", "alter table t1 add Id INT;", true},                                 // 允许类型
	{"MTB-L2-008", "alter table t1 add Id TINYINT;", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 add Id SMALLINT;", true},                            // 允许类型
	{"MTB-L2-008", "alter table t1 add Id MEDIUMINT;", true},                           // 允许类型
	{"MTB-L2-008", "alter table t1 add Id BIGINT;", true},                              // 允许类型
	{"MTB-L2-008", "alter table t1 add Id VARCHAR(1);", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 add Id CHAR(1);", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 add Id DATE;", true},                                // 允许类型
	{"MTB-L2-008", "alter table t1 add Id TIME;", true},                                // 允许类型
	{"MTB-L2-008", "alter table t1 add Id DATETIME;", true},                            // 允许类型
	{"MTB-L2-008", "alter table t1 add Id YEAR;", true},                                // 允许类型
	{"MTB-L2-008", "alter table t1 add Id DECIMAL(10,1);", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 add Id FLOAT;", false},                              // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id FLOAT(3,1);", false},                         // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id DOUBLE;", false},                             // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id DOUBLE(10,1);", false},                       // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id REAL;", false},                               // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id BINARY(10);", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 add Id VARBINARY(10);", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 add Id TINYBLOB;", true},                            // 允许类型
	{"MTB-L2-008", "alter table t1 add Id BLOB;", true},                                // 允许类型
	{"MTB-L2-008", "alter table t1 add Id MEDIUMBLOB;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 add Id LONGBLOB;", true},                            // 允许类型
	{"MTB-L2-008", "alter table t1 add Id TINYTEXT;", true},                            // 允许类型
	{"MTB-L2-008", "alter table t1 add Id TEXT;", true},                                // 允许类型
	{"MTB-L2-008", "alter table t1 add Id MEDIUMTEXT;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 add Id LONGTEXT;", true},                            // 允许类型
	{"MTB-L2-008", "alter table t1 add Id TIMESTAMP;", true},                           // 允许类型
	{"MTB-L2-008", "alter table t1 add Id BIT;", false},                                // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id SET('a','b');", false},                       // 不允许类型
	{"MTB-L2-008", "alter table t1 add Id ENUM('a','b');", false},                      // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id INT;", true},                              // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id TINYINT;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id SMALLINT;", true},                         // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id MEDIUMINT;", true},                        // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id BIGINT;", true},                           // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id VARCHAR(1);", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id CHAR(1);", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id DATE;", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id TIME;", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id DATETIME;", true},                         // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id YEAR;", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id DECIMAL(10,1);", true},                    // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id FLOAT;", false},                           // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id FLOAT(3,1);", false},                      // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id DOUBLE;", false},                          // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id DOUBLE(10,1);", false},                    // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id REAL;", false},                            // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id BINARY(10);", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id VARBINARY(10);", true},                    // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id TINYBLOB;", true},                         // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id BLOB;", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id MEDIUMBLOB;", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id LONGBLOB;", true},                         // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id TINYTEXT;", true},                         // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id TEXT;", true},                             // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id MEDIUMTEXT;", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id LONGTEXT;", true},                         // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id TIMESTAMP;", true},                        // 允许类型
	{"MTB-L2-008", "alter table t1 modify Id BIT;", false},                             // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id SET('a','b');", false},                    // 不允许类型
	{"MTB-L2-008", "alter table t1 modify Id ENUM('a','b');", false},                   // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 INT;", true},                           // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 TINYINT;", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 SMALLINT;", true},                      // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 MEDIUMINT;", true},                     // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 BIGINT;", true},                        // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 VARCHAR(1);", true},                    // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 CHAR(1);", true},                       // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 DATE;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 TIME;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 DATETIME;", true},                      // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 YEAR;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 DECIMAL(10,1);", true},                 // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 FLOAT;", false},                        // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 FLOAT(3,1);", false},                   // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 DOUBLE;", false},                       // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 DOUBLE(10,1);", false},                 // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 REAL;", false},                         // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 BINARY(10);", true},                    // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 VARBINARY(10);", true},                 // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 TINYBLOB;", true},                      // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 BLOB;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 MEDIUMBLOB;", true},                    // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 LONGBLOB;", true},                      // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 TINYTEXT;", true},                      // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 TEXT;", true},                          // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 MEDIUMTEXT;", true},                    // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 LONGTEXT;", true},                      // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 TIMESTAMP;", true},                     // 允许类型
	{"MTB-L2-008", "alter table t1 change id c1 BIT;", false},                          // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 SET('a','b');", false},                 // 不允许类型
	{"MTB-L2-008", "alter table t1 change id c1 ENUM('a','b');", false},                // 不允许类型

	{"MTB-L2-009", "alter table t1 add c22 INT comment 'id',change t2 t2_2 float(3,1) comment 'float';", true},                                 // TableAlterColumnCommentRequired
	{"MTB-L2-009", "alter table t1 add Id INT comment 'id',change t2 t2_2 tinyint comment 't2_2',modify t3 double comment 'float';", true},     // TableAlterColumnCommentRequired
	{"MTB-L2-009", "alter table t1 add c2 INT comment 'id',modify c tinyint comment 'modify',change tc tc2 float(3,1) comment 'fload';", true}, // TableAlterColumnCommentRequired
	{"MTB-L2-009", "alter table t1 add Id INT comment 'id',change t2 t2_2 float(3,1);", false},                                                 // TableAlterColumnCommentRequired
	{"MTB-L2-009", "alter table t1 add Id INT comment 'id',change t2 t2_2 tinyint comment 't2_2',modify t3 double ;", false},                   // TableAlterColumnCommentRequired
	{"MTB-L2-009", "alter table t1 add c2 INT comment 'id',modify c tinyint comment 'modify',change tc tc2 float(3,1);", false},                // TableAlterColumnCommentRequired

	{"MTB-L2-010", "alter table t1 add c2 INT comment 'id';", true},                                                            //TableAlterColumnAvailableCharsets
	{"MTB-L2-010", "alter table t1 add c2 INT comment 'id',modify c tinyint comment 'modify' ;", true},                         //TableAlterColumnAvailableCharsets
	{"MTB-L2-010", "alter table t1 add c2 INT comment 'id',modify c tinyint comment 'modify',change tc tc2 float(3,1);", true}, //TableAlterColumnAvailableCharsets

	{"MTB-L2-010", "alter table t1 add c2  varchar(32) CHARACTER SET utf8mb4 comment 'id';", true},                                                      //TableAlterColumnAvailableCharsets
	{"MTB-L2-010", "alter table t1 add c2 INT comment 'id',modify c char(32)  CHARACTER SET utf8mb4  comment 'modify'  ;", true},                        //TableAlterColumnAvailableCharsets
	{"MTB-L2-010", "alter table t1 add c2 INT comment 'id',modify c tinyint comment 'modify',change tc tc2 varchar(32)  CHARACTER SET utf8mb4 ;", true}, //TableAlterColumnAvailableCharsets

	{"MTB-L2-010", "alter table t1 add c2 char(32) CHARACTER SET latin1  ,modify c tinyint ,change tc tc2 float(3,1);", false}, //TableAlterColumnAvailableCharsets
	{"MTB-L2-010", "alter table t1 add c2 INT,modify c varchar(32)  CHARACTER SET latin1 ,change tc tc2 float(3,1);", false},   //TableAlterColumnAvailableCharsets
	{"MTB-L2-010", "alter table t1 add c2 INT,modify c tinyint,change tc tc2 varchar(32) CHARACTER SET latin1 ;", false},       //TableAlterColumnAvailableCharsets

	{"MTB-L2-011", "alter table t1 add c1 VARCHAR(1) collate utf8_unicode_ci;", false},                                                  //TableAlterColumnAvailableCollates
	{"MTB-L2-011", "alter table t1 add c1 varchar(1), modify c3 VARCHAR(1) collate utf8_unicode_ci;", false},                            //TableAlterColumnAvailableCollates
	{"MTB-L2-011", "alter table t1 add c1 varchar(1), modify c3 VARCHAR(1),change t2 t_2 varchar(1) collate utf8_unicode_ci;", false},   //TableAlterColumnAvailableCollates
	{"MTB-L2-011", "alter table t1 add c1 VARCHAR(1) collate utf8mb4_unicode_ci;", true},                                                //TableAlterColumnAvailableCollates
	{"MTB-L2-011", "alter table t1 add c1 varchar(1), modify c3 VARCHAR(1) collate utf8mb4_unicode_ci;", true},                          //TableAlterColumnAvailableCollates
	{"MTB-L2-011", "alter table t1 add c1 varchar(1), modify c3 VARCHAR(1),change t2 t_2 varchar(1) collate utf8mb4_unicode_ci;", true}, //TableAlterColumnAvailableCollates

	{"MTB-L2-014", "alter table t1 add index index_2(name,password);", true},                                                     //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add index index_2(name,password),add key index_2(n,p);", true},                                //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add unique index unique_2(name,password);", true},                                             //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add unique index unique_2(name,password);", true},                                             //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add unique index unique_2(name,password);", true},                                             //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 drop index idx_c, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true},   //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add index ind2ex_2(name,password);", true},                                                    //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add index index_2(name,password),add key in2dex_2(n,p);", true},                               //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add unique index uniqu2e_2(name,password);", true},                                            //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 add unique index unique_2(name,password),add index idx_2(name);", true},                       //TableAlterIndexNameExplicit
	{"MTB-L2-014", "alter table t1 drop index idx_c, add constraint 22FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterIndexNameExplicit

	{"MTB-L2-015", "alter table t1 add index index_2(name,password);", true},                               //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add index index_2(name,password),add key index_2(n,p);", true},          //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add unique index unique_2(name,password);", true},                       //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add unique index unique_2(name,password);", true},                       //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add unique index unique_2(name,password);", true},                       //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add index ind2ex_2(name,password);", true},                              //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add index index_2(name,password),add key in2dex_2(n,p);", true},         //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add unique index uniqu2e_2(name,password);", true},                      //TableAlterIndexNameExplicit
	{"MTB-L2-015", "alter table t1 add unique index unique_2(name,password),add index idx_2(name);", true}, //TableAlterIndexNameExplicit

	{"MTB-L2-016", "alter table t1 add index index_2(name,password);", true},                        //TableAlterIndexNameExplicit
	{"MTB-L2-016", "alter table t1 add index index_2(name,password),add key index_2(n,p);", true},   //TableAlterIndexNameExplicit
	{"MTB-L2-016", "alter table t1 add index ind2eX_2(name,password);", false},                      //TableAlterIndexNameExplicit
	{"MTB-L2-016", "alter table t1 add index inUex_2(name,password),add key in2dex_2(n,p);", false}, //TableAlterIndexNameExplicit

	{"MTB-L2-017", "alter table t1 add index index_7890(name,password);", true},                            //TableAlterIndexNameMaxLength
	{"MTB-L2-017", "alter table t1 add index index_78901(name,password),add key index_78902(n,p);", false}, //TableAlterIndexNameMaxLength
	{"MTB-L2-017", "alter table t1 add index ind2eX_2(name,password);", true},                              //TableAlterIndexNameMaxLength

	{"MTB-L2-018", "alter table t1 add index index_2(name,password);", true},                                                                                        //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add index index_222222222222222222222222222222(name,password),add key index_2(n,p);", true},                                      //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add unique index unique_2(name,password);", true},                                                                                //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                    //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add index ind2eX_2(name,password);", false},                                                                                      //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", false},                                   //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add unique index uUiqu2e_2(name,password);", true},                                                                               //TableAlterIndexNamePrefixRequired
	{"MTB-L2-018", "alter table t1 add unique index Unique_2(name,password),add index idx_2(name);", false},                                                         //TableAlterIndexNamePrefixRequired
	//^index_[1-9][1-9]*$
	{"MTB-L2-018", "alter table t1 drop index idx_c, add constraint 22FK_ID22232342423424234242 foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterIndexNamePrefixRequired

	//TableAlterUniqueNameExplicit
	{"MTB-L2-019", "alter table t1 add unique index (name,password);", false},            //TableAlterUniqueNameExplicit
	{"MTB-L2-019", "alter table t1 add unique index unique_2(name,password);", true},     //TableAlterUniqueNameExplicit
	{"MTB-L2-019", "alter table t1 drop index idx_c, add unique index `  `(id);", false}, //TableAlterUniqueNameExplicit
	{"MTB-L2-019", "alter table t1 drop index idx_c, add unique index ``(id);", false},   //

	{"MTB-L2-020", "alter table t1 add index index_2(name,password);", true},                                                                                        //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add index index_222222222222222222222222222222(name,password),add key index_2(n,p);", true},                                      //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add unique index unique_2(name,password);", true},                                                                                //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                    //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                       //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", true},                                    //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add unique index `uUiqu2e_2`(name,password);", true},                                                                             //TableAlterUniqueNameQualified
	{"MTB-L2-020", "alter table t1 add unique index Unique_2(name,password),add index idx_2(name);", true},                                                          //TableAlterUniqueNameQualified 跳过空串检测
	{"MTB-L2-020", "alter table t1 drop index idx_c, add unique index `  `(id);", true},                                                                             //TableAlterUniqueNameQualified 跳过空串检测
	{"MTB-L2-020", "alter table t1 drop index idx_c, add unique index ``(id);", true},

	//TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add index index_2(name,password);", true},                                                                                        //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add index index_222222222222222222222222222222(name,password),add key index_2(n,p);", true},                                      //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add unique index unique_2(name,password);", true},                                                                                //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                    //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                       //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", true},                                    //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add unique index `uUiqu2e_2`(name,password);", false},                                                                            //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 add unique index Unique_2(name,password),add index idx_2(name);", false},                                                         //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 drop index idx_c, add unique index `  `(id);", true},                                                                             //TableAlterUniqueNameLowerCaseRequired
	{"MTB-L2-021", "alter table t1 drop index idx_c, add unique index ``(id);", true},

	{"MTB-L2-022", "alter table t1 add index index_2(name,password);", true},                                                                                        //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add index index_222222222222222222222222222222(name,password),add key index_2(n,p);", true},                                      //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add unique index unique_2(name,password);", true},                                                                                //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", false},                                                   //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                       //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", true},                                    //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add unique index `uUiqu2e_2`(name,password);", true},                                                                             //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 add unique index Unique_2(name,password),add index idx_2(name);", true},                                                          //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 drop index idx_c, add unique index `  `(id);", true},                                                                             //TableAlterUniqueNameMaxLength
	{"MTB-L2-022", "alter table t1 drop index idx_c, add unique index ``(id);", true},

	{"MTB-L2-023", "alter table t1 add index index_2(name,password);", true},                                                                                        //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add index index_222222222222222222222222222222(name,password),add key index_2(n,p);", true},                                      //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add unique index unique_2(name,password);", true},                                                                                //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                    //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                       //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", true},                                    //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add unique index `uUiqu2e_2`(name,password);", false},                                                                            //TableAlterUniqueNamePrefixRequired
	{"MTB-L2-023", "alter table t1 add unique index Unique_2(name,password),add index idx_2(name);", false},                                                         //TableAlterUniqueNamePrefixRequired 跳过空串检测
	{"MTB-L2-023", "alter table t1 drop index idx_c, add unique index `  `(id);", true},                                                                             //TableAlterUniqueNamePrefixRequired 跳过空串检测
	{"MTB-L2-023", "alter table t1 drop index idx_c, add unique index ``(id);", true},

	{"MTB-L2-024", "alter table t1 add index index_2(name,password);", true},                                                                                         //TableAlterForeignKeyEnabled
	{"MTB-L2-024", "alter table t1 add index index_222222222222222222222222222222(name,password),add key index_2(n,p);", true},                                       //TableAlterForeignKeyEnabled
	{"MTB-L2-024", "alter table t1 add unique index unique_2(name,password);", true},                                                                                 //TableAlterForeignKeyEnabled
	{"MTB-L2-024", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                     //TableAlterForeignKeyEnabled
	{"MTB-L2-024", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", false}, //TableAlterForeignKeyEnabled

	{"MTB-L2-025", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                    //TableAlterForeignKeyNameExplicit
	{"MTB-L2-025", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterForeignKeyNameExplicit
	{"MTB-L2-025", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                       //TableAlterForeignKeyNameExplicit
	{"MTB-L2-025", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", true},                                    //TableAlterForeignKeyNameExplicit

	{"MTB-L2-026", "alter table t1 add unique index unique_2(name,password);", true},                                                    //TableAlterForeignKeyNameQualified
	{"MTB-L2-026", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                        //TableAlterForeignKeyNameQualified
	{"MTB-L2-026", "alter table t1 drop index idx_c234242424, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterForeignKeyNameQualified
	{"MTB-L2-026", "alter table t1 add index ind2eX_2(name,password);", true},                                                           //TableAlterForeignKeyNameQualified
	{"MTB-L2-026", "alter table t1 add unique index `uUiqu2e_2`(name,password);", true},                                                 //TableAlterForeignKeyNameQualified

	{"MTB-L2-027", "alter table t1 add unique index unique_2(name,password);", true},                                                                                 //TableAlterForeignKeyNameLowerCaseRequired
	{"MTB-L2-027", "alter table t1 add unique index unique_22333333333234222222222222222(name,password);", true},                                                     //TableAlterForeignKeyNameLowerCaseRequired
	{"MTB-L2-027", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint FK_ID foreign key(user_id) REFERENCES tb_user(id);", false}, //TableAlterForeignKeyNameLowerCaseRequired
	{"MTB-L2-027", "alter table t1 drop index idx_c234242424zzzzzzzzzzzzzzzzzzzzzzzzzzzz, add constraint fk_id foreign key(user_id) REFERENCES tb_user(id);", true},  //TableAlterForeignKeyNameLowerCaseRequired
	{"MTB-L2-027", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                        //TableAlterForeignKeyNameLowerCaseRequired
	{"MTB-L2-027", "alter table t1 add index inUex_2(name,password),add key index_22342424242424324242342424242423(n,p);", true},                                     //TableAlterForeignKeyNameLowerCaseRequired

	//TableAlterForeignKeyNameMaxLength

	{"MTB-L2-028", "alter table t1 add unique index unique_2(name,password);", true},                                                                          //TableAlterForeignKeyNameMaxLength
	{"MTB-L2-028", "alter table t1 drop index idx_c234242424, add constraint FK_ID22222222222222212345 foreign key(user_id) REFERENCES tb_user(id);", true},   //TableAlterForeignKeyNameMaxLength
	{"MTB-L2-028", "alter table t1 drop index idx_c234242424, add constraint FK_ID222222222222222123456 foreign key(user_id) REFERENCES tb_user(id);", false}, //TableAlterForeignKeyNameMaxLength
	{"MTB-L2-028", "alter table t1 add index ind2eX_2(name,password);", true},                                                                                 //TableAlterForeignKeyNameMaxLength

	//TableAlterForeignKeyNamePrefixRequired
	{"MTB-L2-029", "alter table t1 add unique index unique_2(name,password);", true},                                                             //TableAlterForeignKeyNamePrefixRequired
	{"MTB-L2-029", "alter table t1 drop index idx_c234242424 , add constraint FK_ID_222222 foreign key(user_id) REFERENCES tb_user(id);", false}, //TableAlterForeignKeyNamePrefixRequired
	{"MTB-L2-029", "alter table t1 drop index idx_c234242424, add constraint fk_id foreign key(user_id) REFERENCES tb_user(id);", true},          //TableAlterForeignKeyNamePrefixRequired
	{"MTB-L2-029", "alter table t1 add index ind2eX_2(name,password);", true},                                                                    //TableAlterForeignKeyNamePrefixRequired

	// TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 add unique index unique_2(name,password);", true},                          //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to zzzz;", true},                                                    //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to UNNNzzzPP;", true},                                               //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to UNN________________NzzzPP;", true},                               //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to `UNNNzz_____$@#$___zPP`;", false},                                //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to `111111`;", false},                                               //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1  add constraint fk_id foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to ` `;", false},                                                    //TableAlterNewTableNameQualified
	{"MTB-L2-030", "alter table t1 rename to ``;", false},                                                     //TableAlterNewTableNameQualified

	//TableAlterNewTableNameLowerCaseRequired
	{"MTB-L2-031", "alter table t1 add unique index unique_2(name,password);", true},                          //TableAlterNewTableNameLowerCaseRequired
	{"MTB-L2-031", "alter table t1 rename to zzzz;", true},                                                    //TableAlterNewTableNameLowerCaseRequired
	{"MTB-L2-031", "alter table t1 rename to UNNNzzzPP;", false},                                              //TableAlterNewTableNameLowerCaseRequired
	{"MTB-L2-031", "alter table t1  add constraint fk_id foreign key(user_id) REFERENCES tb_user(id);", true}, //TableAlterNewTableNameLowerCaseRequired
	{"MTB-L2-031", "alter table t1 rename to ` `;", true},                                                     //TableAlterNewTableNameLowerCaseRequired
	{"MTB-L2-031", "alter table t1 rename to ``;", true},                                                      //TableAlterNewTableNameLowerCaseRequired

	//TableAlterNewTableNameMaxLength
	{"MTB-L2-032", "alter table t1 add unique index unique_2(name,password);", true}, //TableAlterNewTableNameMaxLength
	{"MTB-L2-032", "alter table t1 rename to zzzz1111111111111111;", true},           //TableAlterNewTableNameMaxLength
	{"MTB-L2-032", "alter table t1 rename to zzzz11111111111111123;", false},         //TableAlterNewTableNameMaxLength
	{"MTB-L2-032", "alter table t1 rename to ` `;", true},                            //TableAlterNewTableNameMaxLength
	{"MTB-L2-032", "alter table t1 rename to ``;", true},                             //TableAlterNewTableNameMaxLength
	{"MTB-L2-039", "ALTER TABLE t4 ADD KEY index_1 (c1, c2, c3, c4);", false},        //
	{"MTB-L2-039", "ALTER TABLE t4 ADD KEY index_1 (c1, c2, c3);", true},             //

	{"MTB-L2-014", "alter table t1 add key index_1 (c1);", true},                    //
	{"MTB-L2-010", "alter table t1 add c1 VARCHAR(1) CHARACTER SET utf8mb4;", true}, //
	{"MTB-L2-010", "alter table t1 add c1 VARCHAR(1) CHARACTER SET latin1;", false},
}

// TestTableCreateTableCharsetCollateMatch 建表是校验规则与字符集必须匹配
// RULE:CTB-L2-003
var updateCasess = []RuleTestCase{
	{"UPD-L2-001", "update t1 set id =1;", false},                                          // UpdateWithoutWhereEnabled
	{"UPD-L2-001", "update t1 set id =2 where 1=1;", true},                                 // UpdateWithoutWhereEnabled
	{"UPD-L2-001", "update t1 set id =2 where 'z'>'a';", true},                             // UpdateWithoutWhereEnabled
	{"UPD-L2-001", "update t1 ,c set c.id =t.id where c.id >1;", true},                     // UpdateWithoutWhereEnabled
	{"UPD-L2-001", "update t1 c set c.id=15 limit 10;", false},                             // UpdateWithoutWhereEnabled
	{"UPD-L2-001", "update t1 ,c set c.id =t.id where c.id in (select id from t2);", true}, // UpdateWithoutWhereEnabled
	{"UPD-L3-001", "update starwars.users set name='x' where id=1;", true},                 // 目标库存在
	{"UPD-L3-001", "update nodb.users set name='x' where id=1;", false},                    // 目标库不存在
	{"UPD-L3-002", "update users set name='x' where id=1;", true},                          // 目标表存在
	{"UPD-L3-002", "update starwars.nouser set name='x' where id=1;", false},               // 目标表不存在
	{"UPD-L3-003", "update users set name='x' where id=1;", true},                          // 目标列存在
	{"UPD-L3-003", "update users set bogus='x' where id=1;", false},                        // 目标列不存在
	{"UPD-L3-005", "update users set name='x' limit 100;", true},                           // 未超过最大更新行数
	{"UPD-L3-005", "update users set name='x';", false},                                    // 无 LIMIT，视为不满足
}

// TestUpdateTargetDatabaseExists 目标库是否存在
var viewCases = []RuleTestCase{
	{"CVW-L2-001", "create view tHLLLzz_21324___ as select * from c;", true},            // ViewCreateViewNameQualified
	{"CVW-L2-001", "create view `tHLLzz 2124__` as select * from c;", false},            // ViewCreateViewNameQualified
	{"CVW-L2-002", "create view thllzzz_21324___ as select * from c;", true},            // ViewCreateViewNameLowerCaseRequired
	{"CVW-L2-002", "create view tHLLzz_21324___ as select * from c;", false},            // ViewCreateViewNameLowerCaseRequired
	{"CVW-L2-003", "create view cczz as select * from t ;", true},                       // ViewCreateViewNameMaxLength
	{"CVW-L2-003", "create view tzzzzzzzzzzzzzzzzzzzzzzzz as select * from t;", true},   // ViewCreateViewNameMaxLength
	{"CVW-L2-003", "create view tzzzzzzzzzzzzzzzzzzzzzzzz1 as select * from t;", false}, // ViewCreateViewNameMaxLength
	{"CVW-L2-004", "create view vw_concrete_view as select * from t;", true},            // ViewCreateViewNamePrefixRequired
	{"CVW-L2-004", "create view v_concrete_view as select * from t;", false},            // ViewCreateViewNamePrefixRequired
	{"CVW-L3-001", "create view starwars.v_ok as select 1;", true},                      // 目标库存在
	{"CVW-L3-001", "create view nodb.v_bad as select 1;", false},                        // 目标库不存在
	{"CVW-L3-002", "create view starwars.v_new as select 1;", true},                     // 视图不存在
	{"CVW-L3-002", "create view starwars.halo_v as select 1;", false},                   // 视图已存在
	{"MVW-L3-001", "alter view starwars.halo_v as select 1;", true},                     // 目标库存在
	{"MVW-L3-001", "alter view nodb.v_bad as select 1;", false},                         // 目标库不存在
	{"MVW-L3-002", "alter view starwars.halo_v as select 1;", true},                     // 视图存在
	{"MVW-L3-002", "alter view starwars.no_v as select 1;", false},                      // 视图不存在
	{"DVW-L3-001", "drop view starwars.v_ok;", true},                                    // 目标库存在
	{"DVW-L3-001", "drop view nodb.v_bad;", false},                                      // 目标库不存在
	{"DVW-L3-002", "drop view starwars.halo_v;", true},                                  // 视图存在
	{"DVW-L3-002", "drop view starwars.no_v;", false},                                   // 视图不存在
}

// TestViewCreateTargetDatabaseExists 目标库是否存在
var insertCases = []RuleTestCase{
	{"INS-L2-001", "insert into t1 values(1,2,3,'c');", false},                                           // InsertExplicitColumnRequired
	{"INS-L2-001", "insert into t1 select 1,a,b,z from t2 ;", false},                                     // InsertExplicitColumnRequired
	{"INS-L2-002", "insert into t1(a,b,c,d) select a,b,c,d from t2  where t2.id > 100  limit 10", false}, // InsertUsingSelectEnabled
	{"INS-L2-002", "insert into t1 select * from c;", false},                                             // InsertUsingSelectEnabled
	{"INS-L2-002", "insert into t1(a,b,c,d) values(1,1,1,1);", true},                                     // InsertUsingSelectEnabled
	{"INS-L2-002", "insert into t1(a,b,c,d) select a,b,c,d from t2 limit 10", false},                     // InsertUsingSelectEnabled
	{"INS-L2-005", "insert into users (name,email) values ('a','b');", true},                             // 列与值数量一致
	{"INS-L2-005", "insert into users (name,email) values ('a');", false},                                // 列与值数量不一致
	{"INS-L3-001", "insert into starwars.users (name,email) values ('a','b');", true},                    // 目标库存在
	{"INS-L3-001", "insert into nodb.users (name,email) values ('a','b');", false},                       // 目标库不存在
	{"INS-L3-002", "insert into users (name,email) values ('a','b');", true},                             // 目标表存在
	{"INS-L3-002", "insert into starwars.notable (id) values (1);", false},                               // 目标表不存在
	{"INS-L3-003", "insert into users (name,email) values ('a','b');", true},                             // 目标列存在
	{"INS-L3-003", "insert into users (name,bogus) values ('a','b');", false},                            // 目标列不存在
	{"INS-L3-004", "insert into users (name,email) values ('a','b');", true},                             // 非空列均有值
	{"INS-L3-004", "insert into users (name) values ('a');", false},                                      // 缺少非空 email 列
	{"INS-L3-006", "insert into t1(id) values(1),(2);", true},                                            // 未超过行数上限
	{"INS-L3-006", insertManyValues(1001), false},                                                        // 超过行数上限
}

// TestInsertExplicitColumnRequired 是否要求显式列申明
var indexCases = []RuleTestCase{
	{"CIX-L2-001", "CREATE INDEX index_1 ON t1(name,a,t,c,d,e,f);", false}, // IndexCreateIndexMaxColumnLimit
	{"CIX-L2-001", "CREATE INDEX index_1 ON t1 (c1, c2, c3, c4);", false},  // IndexCreateIndexMaxColumnLimit
	{"CIX-L2-001", "CREATE INDEX index_1 ON t1 (c1, c2, c3);", true},       // IndexCreateIndexMaxColumnLimit
	{"CIX-L2-001", "CREATE UNIQUE INDEX unique_1 ON t1(name);", true},      // IndexCreateIndexMaxColumnLimit
	{"CIX-L2-002", "CREATE INDEX indexa ON t1(name);", true},               // IndexCreateIndexNameQualified 全小写
	{"CIX-L2-002", "CREATE INDEX INDEXA ON t1(name);", true},               // IndexCreateIndexNameQualified 全大写
	{"CIX-L2-002", "CREATE INDEX index1 ON t1(name);", true},               // IndexCreateIndexNameQualified 小写+数字
	{"CIX-L2-002", "CREATE INDEX index_2 ON t1(name);", true},              // IndexCreateIndexNameQualified 小写+数字+下划线
	{"CIX-L2-002", "CREATE INDEX index2__a_dd_3 ON t1(name);", true},       // IndexCreateIndexNameQualified 小写+数字+下划线
	{"CIX-L2-002", "CREATE INDEX index________ ON t1(name);", true},        // IndexCreateIndexNameQualified 小写+下划线
	{"CIX-L2-002", "CREATE INDEX INDEX1 ON t1(name);", true},               // IndexCreateIndexNameQualified 大写+数字
	{"CIX-L2-002", "CREATE INDEX INDEX_2 ON t1(name);", true},              // IndexCreateIndexNameQualified 大写+数字+下划线
	{"CIX-L2-002", "CREATE INDEX INDEX2__A_DD_3 ON t1(name);", true},       // IndexCreateIndexNameQualified 大写+数字+下划线
	{"CIX-L2-002", "CREATE INDEX INDEX________ ON t1(name);", true},        // IndexCreateIndexNameQualified 大写+下划线
	{"CIX-L2-002", "CREATE INDEX Idx_2 ON t1(name);", true},                // IndexCreateIndexNameQualified 混合
	{"CIX-L2-002", "CREATE INDEX `id x2` ON t1(name);", false},             // IndexCreateIndexNameQualified 有空格
	{"CIX-L2-002", "CREATE INDEX _indexa ON t1(name);", false},             // IndexCreateIndexNameQualified 非字符开头
	{"CIX-L2-002", "CREATE INDEX 9indexa ON t1(name);", false},             // IndexCreateIndexNameQualified 非字符开头
	{"CIX-L2-003", "CREATE INDEX indexabc ON t1(name,name);", true},        // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX index2 ON t1(name,name);", true},          // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX index_ ON t1(name,name);", true},          // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX index_2abc ON t1(name,name);", true},      // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX Indexabc ON t1(name,name);", false},       // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX Index2 ON t1(name,name);", false},         // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX Index_ ON t1(name,name);", false},         // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-003", "CREATE INDEX Index_2abc ON t1(name,name);", false},     // IndexCreateIndexNameLowerCaseRequired
	{"CIX-L2-004", "CREATE INDEX index_2adf on t1(name);", true},           // IndexCreateIndexNameMaxLength
	{"CIX-L2-004", "CREATE INDEX index_2adfd on t1(name);", false},         // IndexCreateIndexNameMaxLength
	{"CIX-L2-005", "CREATE INDEX index_2 ON t1(name);", true},              // IndexCreateIndexNamePrefixRequired
	{"CIX-L2-005", "CREATE INDEX idx_2 ON t1(name);", false},               // IndexCreateIndexNamePrefixRequired
	{"CIX-L2-006", "CREATE INDEX index_2HGHK ON t1(name);", true},          // IndexCreateDuplicateIndexColumn
	{"CIX-L2-006", "CREATE INDEX index_2 ON t1(name,name);", false},        // IndexCreateDuplicateIndexColumn
	{"CIX-L3-001", "CREATE INDEX index_2 ON t1(name);", true},              //
	{"CIX-L3-001", "CREATE INDEX index_2 ON starwars.t1(name);", true},     //
	{"CIX-L3-001", "CREATE INDEX index_2 ON mock.t1(name);", false},        //
	{"CIX-L3-002", "CREATE INDEX index_2 ON t1(name);", true},              //
	{"CIX-L3-002", "CREATE INDEX index_2 ON starwars.t1(name);", true},     //
	{"CIX-L3-002", "CREATE INDEX index_2 ON mock.t1(name);", false},        //
	{"CIX-L3-002", "CREATE INDEX index_2 ON mock(name);", false},           //
	{"CIX-L3-003", "CREATE INDEX index_2 ON t1(id);", true},                //
	{"CIX-L3-003", "CREATE INDEX index_2 ON starwars.t1(id);", true},       //
	{"CIX-L3-003", "CREATE INDEX index_2 ON t1(name);", false},             //
	{"CIX-L3-003", "CREATE INDEX index_2 ON starwars.t1(name);", false},    //
	{"CIX-L3-003", "CREATE INDEX index_2 ON mock.t1(name);", true},         // 表不存在，不处理
	{"CIX-L3-004", "CREATE INDEX index_1 ON users(status);", true},         // 与已有索引不重叠
	{"CIX-L3-004", "CREATE INDEX index_1 ON users(email);", false},         // 与 unique_1(email) 重叠
	{"CIX-L3-005", "CREATE INDEX index_1 ON users(email);", true},          // 索引名不重复
	{"CIX-L3-005", "CREATE INDEX unique_1 ON users(status);", false},       // 索引名重复
	{"CIX-L3-006", "CREATE INDEX index_6 ON users(status);", true},         // 未超过索引数量上限
	{"CIX-L3-006", "CREATE INDEX index_6 ON idx_many(id);", false},         // 超过索引数量上限
	{"CIX-L3-007", "CREATE INDEX index_1 ON users(status);", true},         // 非文本/二进制列
	{"CIX-L3-007", "CREATE INDEX index_1 ON large_t(payload);", false},     // TEXT 列不允许建索引
	{"RIX-L3-001", "drop index index_1 on idx_many;", true},                // 目标库存在
	{"RIX-L3-001", "drop index index_1 on mock.idx_many;", false},          // 目标库不存在
	{"RIX-L3-002", "drop index index_1 on idx_many;", true},                // 目标表存在
	{"RIX-L3-002", "drop index index_1 on notable;", false},                // 目标表不存在
	{"RIX-L3-003", "drop index index_1 on idx_many;", true},                // 目标索引存在
	{"RIX-L3-003", "drop index index_missing on idx_many;", false},         // 目标索引不存在
}

// TestIndexCreateTargetDatabaseExists 添加索引的表所属库是否存在
var deleteCases = []RuleTestCase{
	{"DEL-L2-001", "delete from t1 where t.id >10;", true},        //  DeleteWithoutWhereEnabled
	{"DEL-L2-001", "delete from t1 where 1=1 and t.id>1; ", true}, //  DeleteWithoutWhereEnabled
	{"DEL-L2-001", "delete from t1 where 1>=1;", true},            //  DeleteWithoutWhereEnabled
	{"DEL-L2-001", "delete from t1 limit 10", false},              //  DeleteWithoutWhereEnabled
	{"DEL-L2-001", "delete from t1 ;", false},                     //  DeleteWithoutWhereEnabled
	{"DEL-L3-002", "delete from t1;", true},                       //
	{"DEL-L3-002", "delete from starwars.t1;", true},              //
	{"DEL-L3-002", "delete from mock.t1;", false},                 // 目标库不存在
	{"DEL-L3-003", "delete from t1;", true},                       //
	{"DEL-L3-003", "delete from starwars.t1;", true},              //
	{"DEL-L3-003", "delete from t2;", false},                      // 目标表不存在
	{"DEL-L3-003", "delete from starwars.t2;", false},             // 目标表不存在
	{"DEL-L3-003", "delete from mock.t1;", false},                 // 目标表不存在
	{"DEL-L3-001", "delete from users limit 10;", true},           // 未超过最大删除行数
	{"DEL-L3-001", "delete from users;", false},                   // 无 LIMIT，视为不满足
	{"DEL-L3-004", "delete from users where id=1;", true},         // 过滤列存在
	{"DEL-L3-004", "delete from users where bogus=1;", false},     // 过滤列不存在
}

// TestDeleteRowsLimit 单次删除的最大行数
var databaseCases = []RuleTestCase{
	{"CDB-L2-001", "create database tickets charset=utf8mb4; ", true},                         // DatabaseCreateAvailableCharsets
	{"CDB-L2-001", "create database tickets; ", false},                                        // DatabaseCreateAvailableCharsets
	{"CDB-L2-001", "create database tickets charset=utf8; ", false},                           // DatabaseCreateAvailableCharsets
	{"CDB-L2-002", "create database tickets charset=utf8 collate utf8_general_ci;", false},    // DatabaseCreateAvailableCollates
	{"CDB-L2-002", "create database tickets;", false},                                         // DatabaseCreateAvailableCollates
	{"CDB-L2-002", "create database tickets charset=utf8 collate utf8mb4_general_ci;", true},  // DatabaseCreateAvailableCollates
	{"CDB-L2-003", "create database tickets charset=utf8 collate utf8mb4_general_ci;", false}, //
	{"CDB-L2-003", "create database tickets charset=utf8 collate utf8_general_ci;", true},     //
	{"CDB-L2-003", "create database tickets charset=utf8mb4 collate utf8mb4_bin;", true},      //
	{"CDB-L2-003", "create database tickets charset=utf8mb4;", true},                          // 不指定排序规则
	{"CDB-L2-003", "create database tickets collate utf8mb4_bin;", true},                      // 不指定字符集
	{"CDB-L2-003", "create database tickets;", true},                                          // 全部不指定
	{"CDB-L2-004", "create database t__icketsAAAACC charset=utf8;", true},                     // DatabaseCreateDatabaseNameQualified
	{"CDB-L2-004", "create database t_0012_icketsAAAACC charset=utf8;", true},                 // DatabaseCreateDatabaseNameQualified
	{"CDB-L2-004", "create database tickets charset=utf8mb4;", true},                          // DatabaseCreateDatabaseNameQualified
	{"CDB-L2-004", "create database ` ` charset=utf8mb4;", false},                             // DatabaseCreateDatabaseNameQualified
	{"CDB-L2-004", "create database `abc$$` charset=utf8mb4;", false},                         // DatabaseCreateDatabaseNameQualified
	{"CDB-L2-005", "create database ticketsAAAACC charset=utf8;", false},                      // DatabaseCreateDatabaseNameLowerCaseRequired
	{"CDB-L2-005", "create database __tickets__ charset=utf8;", true},                         // DatabaseCreateDatabaseNameLowerCaseRequired
	{"CDB-L2-006", "create database tick_0123456789 charset=utf8;", true},                     // DatabaseCreateDatabaseNameMaxLength
	{"CDB-L2-006", "create database tick_01234567890 charset=utf8;", false},                   // DatabaseCreateDatabaseNameMaxLength
	{"CDB-L2-007", "create database starwars charset=utf8;", false},                           //
	{"CDB-L2-007", "create database mock charset=utf8;", true},                                //
	{"MDB-L2-001", "ALTER DATABASE db1 DEFAULT CHARACTER SET utf8mb4;", true},
	{"MDB-L2-001", "ALTER DATABASE db1 CHARACTER SET utf8mb4;", true},
	{"MDB-L2-001", "ALTER DATABASE db1 DEFAULT CHARACTER SET = utf8mb4;", true},
	{"MDB-L2-001", "ALTER DATABASE db1 CHARACTER SET = utf8mb4;", true},
	{"MDB-L2-001", "ALTER DATABASE db1 DEFAULT CHARACTER SET utf8;", false},
	{"MDB-L2-001", "ALTER DATABASE db1 CHARACTER SET utf8;", false},
	{"MDB-L2-001", "ALTER DATABASE db1 DEFAULT CHARACTER SET = utf8;", false},
	{"MDB-L2-001", "ALTER DATABASE db1 CHARACTER SET = utf8;", false},
	{"MDB-L2-002", "ALTER DATABASE db1 DEFAULT COLLATE utf8mb4_unicode_ci;", true},
	{"MDB-L2-002", "ALTER DATABASE db1 COLLATE utf8mb4_unicode_ci;", true},
	{"MDB-L2-002", "ALTER DATABASE db1 DEFAULT COLLATE = utf8mb4_unicode_ci;", true},
	{"MDB-L2-002", "ALTER DATABASE db1 COLLATE = utf8mb4_unicode_ci;", true},
	{"MDB-L2-002", "ALTER DATABASE db1 DEFAULT COLLATE latin1_general_ci;", false},
	{"MDB-L2-002", "ALTER DATABASE db1 COLLATE latin1_general_ci;", false},
	{"MDB-L2-002", "ALTER DATABASE db1 DEFAULT COLLATE = latin1_general_ci;", false},
	{"MDB-L2-002", "ALTER DATABASE db1 COLLATE = latin1_general_ci;", false},
	{"MDB-L2-003", "ALTER DATABASE mock CHARSET = latin1;", true},                                 // 数据库不存在，不判断                           // 原数据库是utf8mb4_general_ci\
	{"MDB-L2-003", "ALTER DATABASE starwars COLLATE = latin1_general_ci;", false},                 // 原数据库是utf8mb4
	{"MDB-L2-003", "ALTER DATABASE starwars CHARSET = latin1;", false},                            // 原数据库是utf8mb4_general_ci\
	{"MDB-L2-003", "ALTER DATABASE starwars CHARSET = latin1 COLLATE = latin1_general_ci;", true}, // 正确
	{"MDB-L2-004", "ALTER DATABASE starwars CHARSET = utf8 COLLATE = latin1_general_ci;", true},   // 目标库存在
	{"MDB-L2-004", "ALTER DATABASE mock COLLATE = latin1_general_ci;", false},                     // 目标库不存在
	{"DDB-L2-001", "drop database starwars;", true},
	{"DDB-L2-001", "drop database mock;", false},
}

var replaceCases = []RuleTestCase{
	{"RPL-L2-002", "replace into users (name,email) values ('a','b');", true},              // 未使用 SELECT
	{"RPL-L2-002", "replace into users (name,email) select name,email from users;", false}, // 不允许 REPLACE...SELECT
	{"RPL-L2-004", "replace into t1(id) values(1);", true},                                 // RPL-L2-004 行数上限
	{"RPL-L2-004", replaceManyValues(1001), false},                                         // 超过行数上限
	{"RPL-L2-005", "replace into users (name,email) values ('a','b');", true},              // 列与值数量一致
	{"RPL-L2-005", "replace into users (name,email) values ('a');", false},                 // 列与值数量不一致
	{"RPL-L3-001", "replace into starwars.users (name,email) values ('a','b');", true},     // 目标库存在
	{"RPL-L3-001", "replace into nodb.users (name,email) values ('a','b');", false},        // 目标库不存在
	{"RPL-L3-002", "replace into users (name,email) values ('a','b');", true},              // 目标表存在
	{"RPL-L3-002", "replace into starwars.nouser (name,email) values ('a','b');", false},   // 目标表不存在
	{"RPL-L3-003", "replace into users (name,email) values ('a','b');", true},              // 目标列存在
	{"RPL-L3-003", "replace into users (name,bogus) values ('a','b');", false},             // 目标列不存在
	{"RPL-L3-004", "replace into users (name,email) values ('a','b');", true},              // 非空列均有值
	{"RPL-L3-004", "replace into users (name) values ('a');", false},                       // 非空列缺失
}

// TestReplaceUsingSelectEnabled 是否允许REPLACE...SELECT
var miscCases = []RuleTestCase{
	{"MSC-L1-002", "FLUSH NO_WRITE_TO_BINLOG TABLES t1 WITH READ LOCK;", false},
	{"MSC-L1-002", "FLUSH TABLES;", false},
	{"MSC-L1-002", "FLUSH TABLES t1;", false},
	{"MSC-L1-002", "FLUSH NO_WRITE_TO_BINLOG TABLES t1;", false},
	{"MSC-L1-002", "FLUSH TABLES WITH READ LOCK;", false},
	{"MSC-L1-002", "FLUSH TABLES t1, t2, t3;", false},
	{"MSC-L1-002", "FLUSH TABLES t1, t2, t3 WITH READ LOCK;", false},
	{"MSC-L1-002", "FLUSH PRIVILEGES;", false},
	{"MSC-L1-002", "FLUSH STATUS;", false},
	{"MSC-L1-002", "SELECT 1;", true},
	{"MSC-L1-003", "truncate table t;", false}, // TruncateTableEnabled
	{"MSC-L1-003", "SELECT 1;", true},
	{"MSC-L1-007", "kill 2;", false}, // KillEnabled
	{"MSC-L1-007", "SELECT 1;", true},
	{"MSC-L1-001", "LOCK TABLES a READ;", false}, // LockTableEnabled
	{"MSC-L1-001", "SELECT 1;", true},
	{"MSC-L1-006", "UNLOCK TABLES;", false}, // UnlockTableEnabled
	{"MSC-L1-006", "SELECT 1;", true},
	{"MSC-L1-005", "PURGE BINARY LOGS BEFORE '2008-04-02 22:46:26';", false},
	{"MSC-L1-005", "PURGE MASTER LOGS TO 'mysql-bin.010';", false},
	{"MSC-L1-005", "PURGE MASTER LOGS BEFORE '2008-04-02 22:46:26';", false},
	{"MSC-L1-005", "SELECT 1;", true},
}

func TestGeneralRules(t *testing.T) {
	a := buildApp(t)
	ensureFixture(t, a)
	ctx := buildContext(t)
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
		valid, ok := runCase(t, ctx, p, c)
		if !ok {
			continue
		}
		assert.Equal(c.Valid, valid, c.RuleName, c.Text)
	}
}

// TestTableCreateColumnCountLimit 表允许的最大列数
// RULE:CTB-L2-014
func TestTableCreateColumnCountLimit(t *testing.T) {
	buildApp(t)
	ctx := buildContext(t)
	p := parser.New()
	assert := assert.New(t)

	cases := []RuleTestCase{
		{"CTB-L2-014", `
		CREATE TABLE t1 (
			a INT,
			b INT,
			c INT,
			d INT,
			e INT,
			f INT,
			g INT,
			h INT,
			i INT,
			j INT,
			k INT,
			l INT,
			m INT,
			n INT,
			o INT,
			p INT,
			q INT,
			r INT,
			s INT,
			t INT,
			u INT,
			v INT,
			w INT,
			x INT,
			y INT
		)
		`, true},
		{"CTB-L2-014", `
		CREATE TABLE t1 (
			a INT,
			b INT,
			c INT,
			d INT,
			e INT,
			f INT,
			g INT,
			h INT,
			i INT,
			j INT,
			k INT,
			l INT,
			m INT,
			n INT,
			o INT,
			p INT,
			q INT,
			r INT,
			s INT,
			t INT,
			u INT,
			v INT,
			w INT,
			x INT,
			y INT,
			z INT
		)
		`, false},
	}

	for _, c := range cases {
		valid, ok := runCase(t, ctx, p, c)
		if !ok {
			continue
		}
		assert.Equal(c.Valid, valid, c.RuleName, c.Text)
	}
}

// mergeCases are the multi-statement cases. They need the multi-statement
// parser, so they cannot go through the single-statement runCase.
var mergeCases = []RuleTestCase{
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

// runMergeCases drives the multi-statement cases for one rule ("" = all).
func runMergeCases(t *testing.T, name string) {
	t.Helper()
	ensureFixture(t, buildApp(t))
	ctx := buildContext(t)
	p := parser.New()

	for _, c := range mergeCases {
		if name != "" && c.RuleName != name {
			continue
		}
		mr := cache.RulesMap.Any(func(elem *models.Rule) bool {
			return elem.Name == c.RuleName
		})
		if mr == nil {
			t.Errorf("rule %s not found", c.RuleName)
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
		// The seed ships some implemented rules disabled (even Bitwise); force
		// the rule on so a case can exercise its handler regardless of seed state.
		rule := bridge.RuleToDomain(mr)
		rule.Bitwise |= 1
		ctx.Rules = []*validate.Rule{{Rule: *rule}}
		validate.RunRegistry(ctx)
	}
}

func TestMergeStatement(t *testing.T) {
	runMergeCases(t, "")
}
