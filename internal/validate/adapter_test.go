package validate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

// TestAdapterWithoutCluster checks that a request without a target cluster
// yields no hits instead of an error: the existence and charset rules cannot
// be checked without the target metadata.
func TestAdapterWithoutCluster(t *testing.T) {
	a := NewAdapter(nil, nil)
	violations, err := a.Validate(context.Background(), domain.ValidateRequest{SQL: "SELECT 1"})
	require.NoError(t, err)
	require.Empty(t, violations)
}
