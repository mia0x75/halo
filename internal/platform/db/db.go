package db

import (
	"context"
	"fmt"
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
		driver = "mysql"
	}
	sqlDriver := "mysql"
	if driver == "postgres" || driver == "pgx" {
		driver, sqlDriver = "postgres", "pgx"
	} else if driver != "mysql" {
		return nil, fmt.Errorf("unknown metadata driver %q: want mysql or postgres", cfg.Driver)
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
