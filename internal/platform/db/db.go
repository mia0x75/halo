package db

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB wraps the metadata connection pool. The metadata store is MySQL-only for
// v0.1 (see TECH_STACK), so the driver is fixed rather than selected here.
type DB struct {
	*sqlx.DB
}

// New opens and verifies the metadata pool. It does not migrate or seed; callers
// run migrations separately so a schema drift never starts the service silently.
func New(ctx context.Context, cfg Config) (*DB, error) {
	conn, err := sqlx.Open("mysql", cfg.Addr)
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

	return &DB{DB: conn}, nil
}
