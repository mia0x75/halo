package migrate

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/migrations"
)

func TestSplitStatements(t *testing.T) {
	script := `
-- 注释里的分号不算结束
CREATE TABLE t1 (id INT COMMENT 'a;b');
/* 块注释 ; 也不算 */
INSERT INTO t1 VALUES (1);
`
	stmts, err := SplitStatements(script)
	require.NoError(t, err)
	require.Len(t, stmts, 2)
	require.Equal(t, "CREATE TABLE t1 (id INT COMMENT 'a;b')", stmts[0])
	require.Equal(t, "INSERT INTO t1 VALUES (1)", stmts[1])
}

func TestSplitStatementsKeepsBacktickIdentifiers(t *testing.T) {
	stmts, err := SplitStatements("ALTER TABLE `mm_a` ADD COLUMN `x` INT;")
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	require.Contains(t, stmts[0], "`mm_a`")
}

func TestSplitStatementsIgnoresEmptyInput(t *testing.T) {
	stmts, err := SplitStatements("-- only a comment\n\n")
	require.NoError(t, err)
	require.Empty(t, stmts)
}

func TestLoadOrdersByVersion(t *testing.T) {
	fsys := fstest.MapFS{
		"0002_second.sql": {Data: []byte("SELECT 2;")},
		"0001_first.sql":  {Data: []byte("SELECT 1;")},
		"README.md":       {Data: []byte("ignored")},
	}

	migrations, err := Load(fsys)
	require.NoError(t, err)
	require.Len(t, migrations, 2)
	require.Equal(t, 1, migrations[0].Version)
	require.Equal(t, "first", migrations[0].Name)
	require.Equal(t, 2, migrations[1].Version)
	require.Equal(t, "second", migrations[1].Name)
}

func TestLoadRejectsDuplicatedVersion(t *testing.T) {
	fsys := fstest.MapFS{
		"0001_first.sql": {Data: []byte("SELECT 1;")},
		"0001_other.sql": {Data: []byte("SELECT 2;")},
	}

	_, err := Load(fsys)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicated")
}

func TestLoadRejectsBadFilename(t *testing.T) {
	fsys := fstest.MapFS{"init.sql": {Data: []byte("SELECT 1;")}}

	_, err := Load(fsys)
	require.Error(t, err)
	require.Contains(t, err.Error(), "NNNN_name.sql")
}

func TestLoadComputesChecksum(t *testing.T) {
	fsys := fstest.MapFS{"0001_first.sql": {Data: []byte("SELECT 1;")}}

	migrations, err := Load(fsys)
	require.NoError(t, err)
	require.Len(t, migrations, 1)
	require.Len(t, migrations[0].Checksum, 64)
}

func TestParseFilename(t *testing.T) {
	version, name, err := parseFilename("0007_add_metric.sql")
	require.NoError(t, err)
	require.Equal(t, 7, version)
	require.Equal(t, "add_metric", name)

	_, _, err = parseFilename("0007.sql")
	require.Error(t, err)

	_, _, err = parseFilename("abc_name.sql")
	require.Error(t, err)
}

func TestVerifyDetectsChecksumDrift(t *testing.T) {
	m := &Migrator{
		migrations: []Migration{
			{Version: 1, Name: "first", Checksum: "aaa"},
		},
	}

	err := m.verify(map[int]Record{1: {Version: 1, Checksum: "bbb", Success: true}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "checksum")
}

func TestVerifyDetectsMissingFile(t *testing.T) {
	m := &Migrator{
		migrations: []Migration{
			{Version: 1, Name: "first", Checksum: "aaa"},
		},
	}

	err := m.verify(map[int]Record{2: {Version: 2, Checksum: "bbb", Success: true}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "缺失")
}

func TestVerifyDetectsVersionGap(t *testing.T) {
	m := &Migrator{
		migrations: []Migration{
			{Version: 1, Name: "first", Checksum: "aaa"},
			{Version: 3, Name: "third", Checksum: "ccc"},
		},
	}

	err := m.verify(map[int]Record{3: {Version: 3, Checksum: "ccc", Success: true}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "不连续")
}

// TestEmbeddedSetsAreLoadable keeps the shipped migration files parseable.
func TestEmbeddedSetsAreLoadable(t *testing.T) {
	mysqlFS, err := migrations.MySQL()
	require.NoError(t, err)

	mysqlMigrations, err := Load(mysqlFS)
	require.NoError(t, err)
	require.NotEmpty(t, mysqlMigrations)
	for _, mig := range mysqlMigrations {
		if mig.Folded {
			continue
		}
		require.NotEmpty(t, mig.Statements, "migration %s has no statement", mig.Filename)
	}

	pgFS, err := migrations.Postgres()
	require.NoError(t, err)

	pgMigrations, err := Load(pgFS)
	require.NoError(t, err)
	require.NotEmpty(t, pgMigrations)
	for _, mig := range pgMigrations {
		if mig.Folded {
			continue
		}
		require.NotEmpty(t, mig.Statements, "migration %s has no statement", mig.Filename)
	}
}

func TestSplitStatementsKeepsDollarQuotedBody(t *testing.T) {
	script := `DO $halo$
BEGIN
    IF EXISTS(SELECT 1 FROM t) THEN
        INSERT INTO a VALUES (1);
    END IF;
END
$halo$;`

	stmts, err := SplitStatements(script)
	require.NoError(t, err)
	require.Len(t, stmts, 1)
	require.Contains(t, stmts[0], "INSERT INTO a VALUES (1);")
}
