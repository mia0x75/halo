package validate_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/validate"
)

// This file provides a database-free harness for the rule cases. Rules are
// parsed from the authoritative seed (scripts/halodb_seed_mysql.sql) instead of
// read through the cache, and the target metadata is synthesised to match what
// the integration fixture provides (starwars with t1(id), plus a "db" schema).

// seedRules parses mm_rules out of the MySQL seed once per test binary. The
// seed is the source of truth, so the unit cases stay in lockstep with the
// rules the server loads.
func seedRules(t *testing.T) map[string]*domain.Rule {
	t.Helper()
	path := filepath.Join("..", "..", "scripts", "halodb_seed_mysql.sql")
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	block := regexp.MustCompile(`(?s)INSERT INTO ` + "`mm_rules`" + ` VALUES\n(.*?)\nUNLOCK TABLES;`).
		FindStringSubmatch(string(data))
	require.NotNil(t, block, "mm_rules block not found in seed")

	body := strings.TrimSuffix(strings.TrimSpace(block[1]), ";")
	rows := regexp.MustCompile(`\),\s*\n\(`).Split(body, -1)
	rows[0] = strings.TrimPrefix(rows[0], "(")
	rows[len(rows)-1] = strings.TrimSuffix(rows[len(rows)-1], ")")

	out := make(map[string]*domain.Rule, len(rows))
	for _, row := range rows {
		f := splitSQLFields(row)
		require.GreaterOrEqual(t, len(f), 11, "short rule row: %s", row)
		level, _ := strconv.Atoi(sqlUnquote(f[4]))
		vldrGroup, _ := strconv.Atoi(sqlUnquote(f[5]))
		bitwise, _ := strconv.Atoi(sqlUnquote(f[8]))
		r := &domain.Rule{
			Name:        sqlUnquote(f[0]),
			Group:       uint8(level), // not read by rules; keep non-zero
			Description: sqlUnquote(f[3]),
			Level:       domain.RuleLevel(level),
			VldrGroup:   uint8(vldrGroup),
			Operator:    sqlUnquote(f[6]),
			Parameters:  sqlUnquote(f[7]),
			Bitwise:     uint8(bitwise),
			Func:        sqlUnquote(f[9]),
			Message:     sqlUnquote(f[10]),
		}
		out[r.Name] = r
	}
	return out
}

// splitSQLFields splits one VALUES tuple on top-level commas, respecting single
// quotes and backslash escapes.
func splitSQLFields(row string) []string {
	var out []string
	var cur strings.Builder
	inQuote, esc := false, false
	for i := 0; i < len(row); i++ {
		c := row[i]
		switch {
		case inQuote:
			cur.WriteByte(c)
			if esc {
				esc = false
			} else if c == '\\' {
				esc = true
			} else if c == '\'' {
				inQuote = false
			}
		case c == '\'':
			inQuote = true
			cur.WriteByte(c)
		case c == ',':
			out = append(out, cur.String())
			cur.Reset()
		default:
			cur.WriteByte(c)
		}
	}
	out = append(out, cur.String())
	return out
}

// sqlUnquote trims a SQL literal and unescapes the sequences the seed uses.
func sqlUnquote(s string) string {
	s = strings.TrimSpace(s)
	if s == "NULL" {
		return ""
	}
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
		s = strings.ReplaceAll(s, `\"`, `"`)
		s = strings.ReplaceAll(s, `\'`, `'`)
		s = strings.ReplaceAll(s, `\\`, `\`)
	}
	return s
}

// unitContext synthesises the metadata the rule cases were authored against:
// the starwars schema (utf8mb4/utf8mb4_unicode_ci) holding t1(id INT, body
// TEXT) with index idx_id(id), plus the users/large_t/idx_many fixtures the DML
// cases use, and a second schema "db". Everything else (t2, mock) is
// deliberately absent so the "does not exist" rules have something to hit.
func unitContext() *validate.Context {
	col := func(name string, t meta.ColumnType) *meta.Column {
		return &meta.Column{Name: name, Type: t}
	}
	return &validate.Context{
		Ticket: &domain.Ticket{Database: "starwars"},
		Databases: []models.Database{
			{Name: "starwars", Charset: "utf8mb4", Collate: "utf8mb4_unicode_ci"},
			{Name: "db", Charset: "utf8mb4", Collate: "utf8mb4_unicode_ci"},
		},
		Tables: map[string][]*meta.Table{
			"starwars": {
				{
					Name:    "t1",
					Engine:  "InnoDB",
					Charset: "utf8mb4",
					Collate: "utf8mb4_unicode_ci",
					Columns: map[string]*meta.Column{
						"id":   col("id", "INT"),
						"body": col("body", "TEXT"),
					},
					Indexes: map[string]*meta.Index{
						"idx_id": {Name: "idx_id", Cols: []string{"id"}},
					},
				},
				{
					Name:    "users",
					Engine:  "InnoDB",
					Charset: "utf8mb4",
					Collate: "utf8mb4_unicode_ci",
					Columns: map[string]*meta.Column{
						"id":     {Name: "id", Type: "INT", AutoIncrement: true},
						"name":   col("name", "VARCHAR"),
						"email":  col("email", "VARCHAR"),
						"status": {Name: "status", Type: "TINYINT", Default: "1"},
					},
					Indexes: map[string]*meta.Index{
						"PRIMARY":  {Name: "PRIMARY", Unique: true, Cols: []string{"id"}},
						"unique_1": {Name: "unique_1", Unique: true, Cols: []string{"email"}},
					},
				},
				{
					// TEXT payload exercises the blob/text index guard (CIX-L3-007).
					Name:    "large_t",
					Engine:  "InnoDB",
					Charset: "utf8mb4",
					Collate: "utf8mb4_unicode_ci",
					Columns: map[string]*meta.Column{
						"id":      col("id", "INT"),
						"payload": col("payload", "TEXT"),
					},
					Indexes: map[string]*meta.Index{},
				},
				{
					// Five indexes trip the per-table index cap (CIX-L3-006, limit 5).
					Name:    "idx_many",
					Engine:  "InnoDB",
					Charset: "utf8mb4",
					Collate: "utf8mb4_unicode_ci",
					Columns: map[string]*meta.Column{
						"id": col("id", "INT"),
					},
					Indexes: map[string]*meta.Index{
						"index_1": {Name: "index_1", Cols: []string{"id"}},
						"index_2": {Name: "index_2", Cols: []string{"id"}},
						"index_3": {Name: "index_3", Cols: []string{"id"}},
						"index_4": {Name: "index_4", Cols: []string{"id"}},
						"index_5": {Name: "index_5", Cols: []string{"id"}},
					},
				},
			},
		},
	}
}
