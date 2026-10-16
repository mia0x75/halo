package dialect

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// postgresTestDSN returns a DSN to a live PostgreSQL target, or skips when none
// is configured. It mirrors the repository suite: HALO_TEST_PG_DSN first, then
// the committed dev config.
func postgresTestDSN(t *testing.T) string {
	t.Helper()
	if dsn := os.Getenv("HALO_TEST_PG_DSN"); dsn != "" {
		return dsn
	}
	data, err := os.ReadFile("../../../bin/config.postgresql.json")
	if err != nil {
		t.Skip("no postgres target: set HALO_TEST_PG_DSN")
	}
	// The dev config points at the metadata instance; reuse it as the target.
	const key = `"addr"`
	i := strings.Index(string(data), key)
	if i < 0 {
		t.Skip("no postgres target: config.postgresql.json has no addr")
	}
	rest := string(data)[i+len(key):]
	rest = rest[strings.Index(rest, `"`)+1:]
	return rest[:strings.Index(rest, `"`)]
}

// TestPostgresDatabasesAreDatabases guards the PG semantics fix: Databases must
// return database names (from pg_database), not schemas. Returning schemas made
// statDatabase and the validate target-database rule reject valid targets.
func TestPostgresDatabasesAreDatabases(t *testing.T) {
	db, err := sqlx.Open("pgx", postgresTestDSN(t))
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	if err := db.PingContext(context.Background()); err != nil {
		t.Skipf("postgres target unreachable: %v", err)
	}

	dbs, err := (Postgres{}).Databases(context.Background(), db)
	require.NoError(t, err)
	require.NotEmpty(t, dbs)

	names := make(map[string]bool, len(dbs))
	for _, d := range dbs {
		names[d.Name] = true
	}
	// A schema name must not appear; the connected database must.
	require.False(t, names["public"], "Databases returned a schema, not a database")
	require.False(t, names["information_schema"], "Databases returned a system schema")
}
