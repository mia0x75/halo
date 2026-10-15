package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB wraps the metadata connection pool. The metadata store runs on MySQL or
// PostgreSQL (see TECH_STACK); Driver selects the database/sql driver.
type DB struct {
	*sqlx.DB
	// Driver is the canonical metadata engine name ("mysql" or "postgres").
	Driver string
}

// New opens and verifies the metadata pool. It does not migrate or seed; callers
// run migrations separately so a schema drift never starts the service silently.
func New(ctx context.Context, cfg Config) (*DB, error) {
	driver, sqlDriver, err := resolveDriver(cfg.Driver, cfg.Addr)
	if err != nil {
		return nil, err
	}
	conn, err := sqlx.Open(sqlDriver, cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("open metadata db: %w", err)
	}

	conn.SetMaxIdleConns(cfg.MaxIdle)
	conn.SetMaxOpenConns(cfg.MaxConnections)
	// WaitTimeout is used as connection lifetime to preserve the behaviour of the
	// legacy database initialization in g/database.go.

	conn.SetConnMaxLifetime(time.Duration(cfg.WaitTimeout) * time.Second)

	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping metadata db: %w", err)
	}

	return &DB{DB: conn, Driver: driver}, nil
}

// DetectDriver sniffs the metadata engine from the DSN so callers can leave
// database.driver unset: postgres(ql):// URLs mean postgres, and anything
// without a URL scheme is the go-sql-driver DSN form, hence mysql.
//
// A mysql:// URL is an error, not mysql: go-sql-driver cannot parse URLs and
// would misread the scheme as the username (surfacing later as a confusing
// access-denied). Use user:password@tcp(host:port)/dbname instead.
func DetectDriver(addr string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(addr))
	if strings.HasPrefix(s, "postgres://") || strings.HasPrefix(s, "postgresql://") {
		return "postgres", nil
	}
	if strings.HasPrefix(s, "mysql://") {
		return "", fmt.Errorf("mysql DSN must use user:password@tcp(host:port)/dbname form, not a URL")
	}
	if strings.Contains(s, "://") {
		return "", fmt.Errorf("unknown DSN scheme: set database.driver explicitly")
	}
	return "mysql", nil
}

// resolveDriver normalizes the configured engine (or the sniffed one when
// unset) and maps it to its database/sql driver name. A configured engine
// that contradicts a positively identified DSN fails fast instead of
// connecting to the wrong database.
func resolveDriver(configured, addr string) (driver, sqlDriver string, err error) {
	detected, err := DetectDriver(addr)
	if err != nil {
		return "", "", err
	}
	if configured == "" {
		configured = detected
	}
	if configured == "pgx" {
		configured = "postgres"
	}
	if configured != detected {
		return "", "", fmt.Errorf("metadata driver %q conflicts with DSN (looks like %q)",
			configured, detected)
	}
	if configured == "postgres" {
		return "postgres", "pgx", nil
	}
	if configured == "mysql" {
		return "mysql", "mysql", nil
	}
	return "", "", fmt.Errorf("unknown metadata driver %q: want mysql or postgres", configured)
}
