package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

// Migration is one schema migration file.
type Migration struct {
	// Version is the leading number of the file name and orders the set.
	Version int
	// Name is the file name without version and extension.
	Name string
	// Filename is the original file name, kept for error messages.
	Filename string
	// Checksum is the sha256 of the raw file content.
	Checksum string
	// Statements are the executable statements of the file.
	Statements []string
	// Folded marks a migration whose changes are already part of the baseline.
	// The migrator records it in the ledger without executing its statements.
	Folded bool
}

// foldedDirective marks a migration as folded into the baseline.
const foldedDirective = "-- migrate:folded-into-baseline"

// Load reads and parses every .sql file of the given set. Files are ordered by
// version; duplicated versions are rejected because the ledger is keyed by
// version.
func Load(fsys fs.FS) ([]Migration, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read migration dir: %w", err)
	}

	migrations := make([]Migration, 0, len(entries))
	seen := make(map[int]string, len(entries))

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, name, err := parseFilename(entry.Name())
		if err != nil {
			return nil, err
		}
		if other, dup := seen[version]; dup {
			return nil, fmt.Errorf("migration version %d duplicated by %s and %s", version, other, entry.Name())
		}
		seen[version] = entry.Name()

		content, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}

		statements, err := SplitStatements(string(content))
		if err != nil {
			return nil, fmt.Errorf("parse migration %s: %w", entry.Name(), err)
		}

		sum := sha256.Sum256(content)
		migrations = append(migrations, Migration{
			Version:    version,
			Name:       name,
			Filename:   entry.Name(),
			Checksum:   hex.EncodeToString(sum[:]),
			Statements: statements,
			Folded:     foldedDirectivePresent(string(content)),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// foldedDirectivePresent reports whether the file asks to be recorded only.
func foldedDirectivePresent(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == foldedDirective {
			return true
		}
	}
	return false
}

// parseFilename splits "0002_add_cluster_type.sql" into version and name.
func parseFilename(filename string) (int, string, error) {
	base := strings.TrimSuffix(filename, ".sql")
	idx := strings.Index(base, "_")
	if idx <= 0 {
		return 0, "", fmt.Errorf("migration file %s must be named NNNN_name.sql", filename)
	}

	version, err := strconv.Atoi(base[:idx])
	if err != nil {
		return 0, "", fmt.Errorf("migration file %s has an invalid version: %w", filename, err)
	}
	if version <= 0 {
		return 0, "", fmt.Errorf("migration file %s must have a positive version", filename)
	}

	return version, base[idx+1:], nil
}

// SplitStatements splits a SQL script into executable statements. It tracks
// quoting and comments so that a semicolon inside a string or a comment does
// not end a statement.
func SplitStatements(content string) ([]string, error) {
	var (
		statements []string
		current    strings.Builder
		runes      = []rune(content)
		inSingle   bool
		inDouble   bool
		inBacktick bool
		inLine     bool
		inBlock    bool
	)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		var next rune
		if i+1 < len(runes) {
			next = runes[i+1]
		}

		switch {
		case inLine:
			if ch == '\n' {
				inLine = false
				current.WriteRune(ch)
			}
			continue
		case inBlock:
			if ch == '*' && next == '/' {
				inBlock = false
				i++
			}
			continue
		case inSingle:
			current.WriteRune(ch)
			if ch == '\\' && next != 0 {
				current.WriteRune(next)
				i++
				continue
			}
			if ch == '\'' {
				inSingle = false
			}
			continue
		case inDouble:
			current.WriteRune(ch)
			if ch == '"' {
				inDouble = false
			}
			continue
		case inBacktick:
			current.WriteRune(ch)
			if ch == '`' {
				inBacktick = false
			}
			continue
		}

		switch ch {
		case '\'':
			inSingle = true
			current.WriteRune(ch)
		case '"':
			inDouble = true
			current.WriteRune(ch)
		case '`':
			inBacktick = true
			current.WriteRune(ch)
		case '$':
			// PostgreSQL dollar quoting: $tag$ ... $tag$ keeps the body verbatim.
			tag, ok := dollarTag(runes, i)
			if !ok {
				current.WriteRune(ch)
				continue
			}
			end := indexRunes(runes, i+len(tag), tag)
			if end < 0 {
				return nil, fmt.Errorf("unterminated dollar quote %s", tag)
			}
			current.WriteString(string(runes[i : end+len(tag)]))
			i = end + len(tag) - 1
		case '-':
			if next == '-' {
				inLine = true
				i++
				continue
			}
			current.WriteRune(ch)
		case '/':
			if next == '*' {
				inBlock = true
				i++
				continue
			}
			current.WriteRune(ch)
		case ';':
			if stmt := strings.TrimSpace(current.String()); stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		default:
			current.WriteRune(ch)
		}
	}

	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements, nil
}

// dollarTag returns the dollar quote tag starting at the given position, for
// example "$$" or "$halo$". It fails when the position does not start one.
func dollarTag(runes []rune, pos int) (string, bool) {
	if pos >= len(runes) || runes[pos] != '$' {
		return "", false
	}
	for i := pos + 1; i < len(runes); i++ {
		ch := runes[i]
		if ch == '$' {
			return string(runes[pos : i+1]), true
		}
		if !isTagChar(ch) {
			return "", false
		}
	}
	return "", false
}

// isTagChar reports whether the rune may appear inside a dollar quote tag.
func isTagChar(ch rune) bool {
	return ch == '_' ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9')
}

// indexRunes finds the next occurrence of the needle at or after the position.
func indexRunes(runes []rune, pos int, needle string) int {
	target := []rune(needle)
	for i := pos; i+len(target) <= len(runes); i++ {
		match := true
		for j := range target {
			if runes[i+j] != target[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
