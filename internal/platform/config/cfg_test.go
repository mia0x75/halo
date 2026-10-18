package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// writeConfig drops a JSON config in a temp dir so ParseConfig can load it
// without touching /etc or the executable directory.
func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cfg.json")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	return path
}

func TestParseConfigLoadsFields(t *testing.T) {
	path := writeConfig(t, `{
		"listen": ":8080",
		"log": {"level": "debug"},
		"database": {"addr": "postgres://localhost:5432/halo", "max_idle": 3, "max_connections": 9, "driver": "postgres"},
		"mail": {"enabled": true, "addr": "smtp.test:465", "user": "halo"},
		"secret": {"crypto": "0123456789abcdef", "jwt": {"issuer": "halo"}}
	}`)

	ParseConfig(path)

	cfg := Config()
	require.NotNil(t, cfg)
	require.Equal(t, ":8080", cfg.Listen)
	require.Equal(t, "debug", cfg.Log.Level)
	require.Equal(t, "postgres", cfg.Database.Driver)
	require.Equal(t, 9, cfg.Database.MaxConnections)
	require.True(t, cfg.Mail.Enabled)
	require.Equal(t, "halo", cfg.Secret.Jwt.Issuer)
	require.Equal(t, path, ConfigFile)
}

// TestParseConfigDefaultsCerts covers the empty-cert/key fallback applied after
// parsing, which the TLS server relies on.
func TestParseConfigDefaultsCerts(t *testing.T) {
	path := writeConfig(t, `{}`)
	ParseConfig(path)

	cfg := Config()
	require.NotNil(t, cfg)
	require.Equal(t, "cert.pem", config.Cert)
	require.Equal(t, "key.pem", config.Key)
}

func TestTemplateConstantsAreDistinct(t *testing.T) {
	tpls := []Template{
		TplTicketCreated, TplTicketUpdated, TplTicketRemoved, TplTicketExecuted,
		TplTicketFailed, TplTicketScheduled, TplTicketClosed, TplTicketMrvFailure,
		TplTicketLgtm, TplUserRegistered, TplPasswordUpdated, TplEmailUpdated,
		TplProfileUpdated, TplUserCreated, TplCommentCreated, TplCronCancelled,
	}
	seen := make(map[Template]struct{}, len(tpls))
	for _, tpl := range tpls {
		require.NotEmpty(t, tpl)
		if _, dup := seen[tpl]; dup {
			t.Fatalf("duplicate template id: %s", tpl)
		}
		seen[tpl] = struct{}{}
	}
}
