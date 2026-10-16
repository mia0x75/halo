package validate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// TestLookup covers known and unknown (group, func) pairs.
func TestLookup(t *testing.T) {
	fn, ok := lookup(100, "AvailableCharsets")
	require.True(t, ok, "group 100 AvailableCharsets must be registered")
	require.NotNil(t, fn)

	// Same func name in another group resolves to a different handler.
	other, ok := lookup(110, "AvailableCharsets")
	require.True(t, ok)
	require.NotNil(t, other)

	// Unknown func within a known group.
	_, ok = lookup(100, "NoSuchRule")
	require.False(t, ok)

	// Unknown group.
	_, ok = lookup(9999, "AvailableCharsets")
	require.False(t, ok)
}

// TestCall_Dispatch verifies Call routes (group, func) to the handler and that
// an unknown handler is a logged no-op rather than a panic.
func TestCall_Dispatch(t *testing.T) {
	s := &Subject{Violations: &models.Violations{}}
	r := &Rule{Rule: domain.Rule{
		Name:      "MSC-L1-001",
		VldrGroup: 230,
		Func:      "LockTableProhibited",
		Level:     1,
		Message:   "lock not allowed",
	}}

	v := &MiscVldr{}
	v.Ctx = &Context{Subjects: []*Subject{s}}

	// Unknown func: must not panic.
	require.NotPanics(t, func() { Call(v, 230, "NoSuchRule", nil, r) })

	// Known func runs without panic (subject has no LockTable node, so no hit).
	require.NotPanics(t, func() { Call(v, 230, "LockTableProhibited", nil, r) })
}

// TestValidateRegistry_Missing logs (does not fail) when a rule has no handler.
func TestValidateRegistry_Missing(t *testing.T) {
	require.NotPanics(t, func() {
		ValidateRegistry([]*domain.Rule{
			{VldrGroup: 100, Func: "AvailableCharsets"},
			{VldrGroup: 230, Func: "GhostRule"},
		})
	})
}
