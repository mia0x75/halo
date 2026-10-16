package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestRuleSourceLoad(t *testing.T) {
	repo := NewRuleRepo(setupTestDB(t))
	ctx := context.Background()

	rule := &domain.Rule{
		Name:        "r" + newSuffix()[:8],
		UUID:        newUUID(),
		Group:       1,
		Description: "a rule",
		Level:       domain.RuleLevelError,
		VldrGroup:   110,
		Operator:    "in",
		Parameters:  "a,b",
		Bitwise:     1,
		Func:        "func",
		Message:     "not allowed",
		Element:     "input",
		Version:     1,
		CreateAt:    nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, rule))

	loaded, err := NewRuleSource(repo).Load(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, loaded)

	var found *domain.Rule
	for _, r := range loaded {
		if r.UUID == rule.UUID {
			found = r
			break
		}
	}
	require.NotNil(t, found, "the first page must include every rule")
	require.Equal(t, rule.Name, found.Name)
	require.Equal(t, rule.VldrGroup, found.VldrGroup)
	require.Equal(t, rule.Level, found.Level)
}
