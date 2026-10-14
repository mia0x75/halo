package agent

import (
	"context"
	"fmt"

	// Register the target drivers; the agent dials them directly.
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/mia0x75/halo/internal/agent/dialect"
	"github.com/mia0x75/halo/internal/domain"
)

// Target 一个已建立连接的目标库，绑定其方言实现。
type Target struct {
	db      *sqlx.DB
	dialect dialect.Dialect
}

// Close 释放连接池。
func (t *Target) Close() error {
	return t.db.Close()
}

// openTarget 按群集类型建立到目标库的连接。密码由调用方解密后传入，
// 不落盘、不缓存。
func openTarget(ctx context.Context, c *domain.Cluster, password string) (*Target, error) {
	// DSN format differs by driver.
	switch c.Type {
	case domain.ClusterTypeMySQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			c.User, password, c.IP, c.Port, c.Database)
		db, err := sqlx.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("open mysql target %s: %w", c.Alias, err)
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			return nil, fmt.Errorf("ping mysql target %s: %w", c.Alias, err)
		}
		return &Target{db: db, dialect: dialect.MySQL{}}, nil
	case domain.ClusterTypePostgres:
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			c.User, password, c.IP, c.Port, c.Database)
		db, err := sqlx.Open("pgx", dsn)
		if err != nil {
			return nil, fmt.Errorf("open postgres target %s: %w", c.Alias, err)
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			return nil, fmt.Errorf("ping postgres target %s: %w", c.Alias, err)
		}
		return &Target{db: db, dialect: dialect.Postgres{}}, nil
	}
	return nil, fmt.Errorf("unsupported cluster type: %s", c.Type)
}
