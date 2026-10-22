// Package seed loads the base data (roles, rules, templates, options) into a
// metadata database that has just been migrated.
package seed

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/platform/db"
	"github.com/mia0x75/halo/internal/platform/migrate"
	"github.com/mia0x75/halo/scripts"
)

// Apply runs the base data script of the given engine.
func Apply(ctx context.Context, pool *db.DB, driver string) error {
	body, err := scripts.Seed(driver)
	if err != nil {
		return err
	}
	statements, err := migrate.SplitStatements(string(body))
	if err != nil {
		return fmt.Errorf("解析基础数据脚本失败: %w", err)
	}
	for _, statement := range statements {
		trimmed := strings.TrimSpace(statement)
		if skip(trimmed) {
			continue
		}
		if _, err := pool.ExecContext(ctx, trimmed); err != nil {
			if duplicate(err) {
				// 迁移已经写过的行不用再写一遍。
				log.Warnf("[W] 基础数据已存在，跳过: %v", err)
				continue
			}
			return fmt.Errorf("执行基础数据脚本失败: %w", err)
		}
	}
	return nil
}

// duplicate reports whether the error is a unique constraint violation.
func duplicate(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

// skip drops the statements that only matter for an interactive import.
// LOCK TABLES pins the table to one connection, so the following inserts would
// block on a different connection of the pool.
func skip(statement string) bool {
	upper := strings.ToUpper(statement)
	return strings.HasPrefix(upper, "LOCK TABLES") ||
		strings.HasPrefix(upper, "UNLOCK TABLES") ||
		strings.HasPrefix(upper, "SET ")
}
