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
	driver := cfg.Driver
	if driver == "" {
		// No explicit engine: sniff it from the DSN so a postgres:// URL
		// just works without touching the config.
		detected, err := DetectDriver(cfg.Addr)
		if err != nil {
			return nil, err
		}
		driver = detected
	}
	sqlDriver := "mysql"
	if driver == "postgres" || driver == "pgx" {
		driver, sqlDriver = "postgres", "pgx"
	} else if driver != "mysql" {
		return nil, fmt.Errorf("unknown metadata driver %q: want mysql or postgres", cfg.Driver)
	}
	if detected, err := DetectDriver(cfg.Addr); err == nil && detected != driver {
		return nil, fmt.Errorf("metadata driver %q conflicts with DSN (looks like %q)",
			cfg.Driver, detected)
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
// database.driver unset: postgres(ql):// URLs mean postgres, mysql:// means
// mysql, and anything without a URL scheme is the go-sql-driver DSN form,
// hence mysql. A URL with any other scheme is an error.
func DetectDriver(addr string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(addr))
	if strings.HasPrefix(s, "postgres://") || strings.HasPrefix(s, "postgresql://") {
		return "postgres", nil
	}
	if strings.HasPrefix(s, "mysql://") {
		return "mysql", nil
	}
	if strings.Contains(s, "://") {
		return "", fmt.Errorf("unknown DSN scheme: set database.driver explicitly")
	}
	return "mysql", nil
}
