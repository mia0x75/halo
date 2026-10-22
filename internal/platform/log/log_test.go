package log

import (
	"os"
	"path/filepath"
	"testing"

	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/platform/config"
)

func loadConfig(t *testing.T, body string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	config.ParseConfig(path)
}

func TestInitLogSetsLevel(t *testing.T) {
	loadConfig(t, `{"log":{"level":"debug"}}`)
	InitLog()
	require.Equal(t, logrus.DebugLevel, logrus.GetLevel())

	loadConfig(t, `{"log":{"level":"warn"}}`)
	InitLog()
	require.Equal(t, logrus.WarnLevel, logrus.GetLevel())
}

// TestInitLogIgnoresInvalidLevel keeps the previous level when the config holds
// an unparsable value instead of failing startup.
func TestInitLogIgnoresInvalidLevel(t *testing.T) {
	loadConfig(t, `{"log":{"level":"warn"}}`)
	InitLog()
	require.Equal(t, logrus.WarnLevel, logrus.GetLevel())

	loadConfig(t, `{"log":{"level":"not-a-level"}}`)
	InitLog()
	require.Equal(t, logrus.WarnLevel, logrus.GetLevel())
}
