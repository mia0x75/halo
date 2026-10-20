package worker

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveCommandKeepsAbsolutePath(t *testing.T) {
	require.Equal(t, "/usr/local/bin/halocli", resolveCommand("/usr/local/bin/halocli"))
}

func TestResolveCommandKeepsOtherCommands(t *testing.T) {
	require.Equal(t, "/bin/sh", resolveCommand("/bin/sh"))
	require.Equal(t, "other-tool", resolveCommand("other-tool"))
}

func TestResolveCommandFallsBackToPath(t *testing.T) {
	// The test binary has no sibling "halocli", so the bare name is kept and
	// the PATH lookup applies at execution time.
	resolved := resolveCommand(cmd)
	require.True(t, resolved == cmd || filepath.IsAbs(resolved))
}
