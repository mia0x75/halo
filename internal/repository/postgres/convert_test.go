package postgresrepo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClampInt32 guards against int->int32 truncation on page sizes.
func TestClampInt32(t *testing.T) {
	require.Equal(t, int32(0), clampInt32(-1))
	require.Equal(t, int32(0), clampInt32(0))
	require.Equal(t, int32(50), clampInt32(50))
	require.Equal(t, int32(math.MaxInt32), clampInt32(math.MaxInt32))
	require.Equal(t, int32(math.MaxInt32), clampInt32(math.MaxInt32+1))
}
