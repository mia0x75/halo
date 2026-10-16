package mysqlrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestRuleRepo(t *testing.T) {
	repo := NewRuleRepo(setupTestDB(t))
	ctx := context.Background()

	rule := &domain.Rule{
		Name:        "r" + newSuffix()[:8],
		UUID:        newUUID(),
		Group:       1,
		Description: "a rule",
		Level:       domain.RuleLevelError,
		VldrGroup:   1,
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

	got, err := repo.GetByUUID(ctx, rule.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.RuleLevelError, got.Level)
	require.Equal(t, rule.Parameters, got.Parameters)

	byName, err := repo.GetByName(ctx, rule.Name)
	require.NoError(t, err)
	require.Equal(t, rule.UUID, byName.UUID)

	updated := nowUnix()
	rule.Level = domain.RuleLevelWarning
	rule.Message = "be careful"
	rule.Version = 2
	rule.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, rule))

	got, err = repo.GetByUUID(ctx, rule.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.RuleLevelWarning, got.Level)
	require.Equal(t, "be careful", got.Message)

	list, err := repo.List(ctx, rule.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasRule(list, rule.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, rule.UUID))
	_, err = repo.GetByUUID(ctx, rule.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasRule(list []*domain.Rule, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
