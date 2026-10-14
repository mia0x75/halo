package mysqlrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
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
		Values:      "a,b",
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

	var found *models.Rule
	for _, r := range loaded {
		if r.UUID == rule.UUID {
			found = r
			break
		}
	}
	require.NotNil(t, found, "the first page must include every rule")
	require.Equal(t, rule.Name, found.Name)
	require.Equal(t, uint16(rule.VldrGroup), found.VldrGroup)
	require.Equal(t, uint8(rule.Level), found.Level)
}

func TestLegacyRuleFromDomain(t *testing.T) {
	updateAt := uint(123)
	in := &domain.Rule{
		Name:        "n",
		UUID:        "u",
		Group:       1,
		Description: "d",
		Level:       domain.RuleLevelWarning,
		VldrGroup:   110,
		Operator:    "in",
		Values:      "v",
		Bitwise:     1,
		Func:        "f",
		Message:     "m",
		Element:     "e",
		Version:     2,
		UpdateAt:    &updateAt,
		CreateAt:    9,
	}

	got := legacyRuleFromDomain(in)
	require.Equal(t, in.Name, got.Name)
	require.Equal(t, in.UUID, got.UUID)
	require.Equal(t, in.Group, got.Group)
	require.Equal(t, uint8(in.Level), got.Level)
	require.Equal(t, uint16(in.VldrGroup), got.VldrGroup)
	require.Equal(t, in.Values, got.Values)
	require.Equal(t, updateAt, got.UpdateAt)
	require.Equal(t, in.CreateAt, got.CreateAt)

	in.UpdateAt = nil
	require.Zero(t, legacyRuleFromDomain(in).UpdateAt)
}
