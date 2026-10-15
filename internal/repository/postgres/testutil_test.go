package postgresrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// testDB is opened once per package run. Every test runs inside its own
// transaction that is rolled back afterwards, so tests never commit and never
// see each other's rows.
var testDB *sql.DB

func TestMain(m *testing.M) {
	if dsn := testDSN(); dsn != "" {
		if db, err := sql.Open("pgx", dsn); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := db.PingContext(ctx); err == nil {
				testDB = db
			} else {
				_ = db.Close()
			}
			cancel()
		}
	}
	code := m.Run()
	if testDB != nil {
		_ = testDB.Close()
	}
	os.Exit(code)
}

// testDSN prefers HALO_TEST_PG_DSN and falls back to the committed PG service
// config, so a developer with a running metadata database gets the tests free.
func testDSN() string {
	if dsn := os.Getenv("HALO_TEST_PG_DSN"); dsn != "" {
		return dsn
	}
	data, err := os.ReadFile("../../../bin/config.postgresql.json")
	if err != nil {
		return ""
	}
	var cfg struct {
		Database struct {
			Addr string `json:"addr"`
		} `json:"database"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ""
	}
	return cfg.Database.Addr
}

// setupTestDB returns a *Queries bound to a rolled-back transaction. Tests skip
// when no metadata database is reachable.
func setupTestDB(t *testing.T) *Queries {
	t.Helper()
	if testDB == nil {
		t.Skip("no metadata database: set HALO_TEST_PG_DSN or provide bin/config.postgresql.json")
	}
	tx, err := testDB.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback()
	})
	return New(tx)
}

func newUUID() string { return uuid.NewString() }

// newSuffix gives short unique strings for VARCHAR columns that carry a unique
// key (names, aliases, emails) so tests never collide with seed rows.
func newSuffix() string { return uuid.NewString()[:8] }

func newEmail() string { return newSuffix() + "-" + newSuffix() + "@example.test" }

func newIP() string {
	return fmt.Sprintf("10.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func newPort() uint { return uint(1024 + rand.Intn(60000)) }

// newID gives a positive key that is unlikely to collide with seed rows in
// tables whose primary key is a composite the test cannot randomise otherwise.
func newID() uint { return uint(1 + rand.Intn(1<<30)) }

func nowUnix() uint { return uint(time.Now().Unix()) }
