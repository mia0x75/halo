package migrate

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/migrations"
)

// TestMigratorUpMySQL applies the shipped MySQL migration set to the test
// database. The test database already holds the schema, so every baseline
// statement is a no-op and the assertions focus on the ledger.
func TestMigratorUpMySQL(t *testing.T) {
	dsn := os.Getenv("HALO_TEST_DSN")
	if dsn == "" {
		t.Skip("no metadata database: set HALO_TEST_DSN")
	}

	db, err := sqlx.Connect("mysql", dsn)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	fsys, err := migrations.MySQL()
	require.NoError(t, err)

	migrator, err := New(db, EngineMySQL, fsys)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err = migrator.Up(ctx)
	require.NoError(t, err)

	status, err := migrator.Status(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, status)
	for _, item := range status {
		require.Equal(t, StateApplied, item.State, "migration %s", item.Migration.Filename)
	}

	// A second run must be a no-op: the ledger keeps every migration recorded.
	applied, err := migrator.Up(ctx)
	require.NoError(t, err)
	require.Empty(t, applied)
}

// TestMigratorUpPostgres applies the shipped PostgreSQL migration set to the
// test database.
func TestMigratorUpPostgres(t *testing.T) {
	dsn := os.Getenv("HALO_TEST_PG_DSN")
	if dsn == "" {
		t.Skip("no metadata database: set HALO_TEST_PG_DSN")
	}

	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	fsys, err := migrations.Postgres()
	require.NoError(t, err)

	migrator, err := New(db, EnginePostgres, fsys)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	_, err = migrator.Up(ctx)
	require.NoError(t, err)

	status, err := migrator.Status(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, status)
	for _, item := range status {
		require.Equal(t, StateApplied, item.State, "migration %s", item.Migration.Filename)
	}

	applied, err := migrator.Up(ctx)
	require.NoError(t, err)
	require.Empty(t, applied)
}
