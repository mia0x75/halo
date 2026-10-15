package validate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// TestRunRegistryEmptyContext exercises the pure engine with no statements: it
// must be a no-op and must not touch the target database.
func TestRunRegistryEmptyContext(t *testing.T) {
	require.Empty(t, RunRegistry(&Context{}))
}

// TestCollectViolations flattens per-statement hits.
func TestCollectViolations(t *testing.T) {
	s := &models.Statement{Violations: &models.Violations{}}
	s.Violations.Append(&models.Clause{Level: 1, Description: "boom"})

	got := collectViolations([]*models.Statement{s})
	require.Len(t, got, 1)
	require.Same(t, s, got[0].Statement)
	require.Equal(t, "boom", got[0].Clause.Description)

	require.Empty(t, collectViolations([]*models.Statement{{}}))
}
