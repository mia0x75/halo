// Package scripts embeds the schema and seed SQL files so a single binary can
// bootstrap a fresh metadata database without external tooling.
package scripts

import (
	"embed"
	"fmt"
)

//go:embed halodb_seed_mysql.sql halodb_seed_postgres.sql
var seedFS embed.FS

// Seed returns the base data script of the given metadata engine.
func Seed(driver string) ([]byte, error) {
	switch driver {
	case "mysql":
		return seedFS.ReadFile("halodb_seed_mysql.sql")
	case "postgres":
		return seedFS.ReadFile("halodb_seed_postgres.sql")
	default:
		return nil, fmt.Errorf("未知的元数据库类型: %s", driver)
	}
}
