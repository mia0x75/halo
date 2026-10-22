package validate

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/driver"

	"github.com/mia0x75/halo/internal/platform/config"
)

// 保留，不可以删除
var _ = driver.ValueExpr{}

// TestMain parses the config flag outside init(), so the testing package's own
// flags are already registered when Parse runs. Paths are tried in order and
// the first existing one wins; the search tolerates the test binary's working
// directory being the package dir rather than the repo root.
func TestMain(m *testing.M) {
	cfg := flag.String("c", "", "configuration file")
	flag.Parse()

	candidates := []string{
		*cfg,
		os.Getenv("HALO_TEST_CONFIG"),
		"bin/config.test.json",
		"../../bin/config.test.json",
		"/etc/halo/config.json",
	}
	for _, path := range candidates {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			config.ParseConfig(path)
			break
		}
	}
	os.Exit(m.Run())
}

func TestSingleSQL(t *testing.T) {
	p := parser.New()
	stmt, _ := p.ParseOneStmt("alter table t1 drop index idx_c, add unique index `  `(id);", "", "")
	for _, spec := range stmt.(*ast.AlterTableStmt).Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq &&
			c.Tp != ast.ConstraintUniqIndex &&
			c.Tp != ast.ConstraintUniqKey {
			continue
		}
		fmt.Printf("Name:<%d>", len(strings.TrimSpace(c.Name)))
	}
}
