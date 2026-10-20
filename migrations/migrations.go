// Package migrations embeds the per-engine schema migration files so a single
// binary can migrate the metadata database without an external tool.
package migrations

import (
	"embed"
	"io/fs"
)

//go:embed mysql/*.sql
var mysqlFS embed.FS

//go:embed postgres/*.sql
var postgresFS embed.FS

// MySQL returns the migration set for the MySQL metadata engine.
func MySQL() (fs.FS, error) {
	return fs.Sub(mysqlFS, "mysql")
}

// Postgres returns the migration set for the PostgreSQL metadata engine.
func Postgres() (fs.FS, error) {
	return fs.Sub(postgresFS, "postgres")
}
