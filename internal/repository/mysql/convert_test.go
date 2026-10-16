package mysqlrepo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClampUint32 guards against int64->uint32 truncation on cursor/id inputs.
func TestClampUint32(t *testing.T) {
	require.Equal(t, uint32(0), clampUint32(-1))
	require.Equal(t, uint32(0), clampUint32(0))
	require.Equal(t, uint32(42), clampUint32(42))
	require.Equal(t, uint32(math.MaxUint32), clampUint32(math.MaxUint32))
	require.Equal(t, uint32(math.MaxUint32), clampUint32(int64(math.MaxUint32)+1))
	require.Equal(t, uint32(math.MaxUint32), clampUint32(math.MaxInt64))
}

// TestClampUint guards against uint->uint32 truncation on 64-bit inputs.
func TestClampUint(t *testing.T) {
	require.Equal(t, uint32(7), clampUint(7))
	require.Equal(t, uint32(math.MaxUint32), clampUint(uint(math.MaxUint32)+1))
}

// TestClampInt32 guards against int->int32 truncation on page sizes.
func TestClampInt32(t *testing.T) {
	require.Equal(t, int32(0), clampInt32(-1))
	require.Equal(t, int32(0), clampInt32(0))
	require.Equal(t, int32(50), clampInt32(50))
	require.Equal(t, int32(math.MaxInt32), clampInt32(math.MaxInt32))
	require.Equal(t, int32(math.MaxInt32), clampInt32(math.MaxInt32+1))
}
