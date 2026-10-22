package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguredNeedsDatabaseAndSecrets(t *testing.T) {
	path := writeConfig(t, `{"listen": ":8080"}`)
	ParseConfig(path)
	require.False(t, Configured(), "缺元数据库与密钥时不算已配置")

	ParseConfig(writeConfig(t, `{"database": {"addr": "root@tcp(127.0.0.1:3306)/halo"}}`))
	require.False(t, Configured(), "缺密钥时不算已配置")

	ParseConfig(writeConfig(t, `{"database": {"addr": "root@tcp(127.0.0.1:3306)/halo"},
		"secret": {"crypto": "0123456789abcdef", "jwt": {"secret_key": "secret"}}}`))
	require.True(t, Configured())
}

func TestSaveWritesFileAndKeepsOtherSections(t *testing.T) {
	path := writeConfig(t, `{"listen": ":8443", "log": {"level": "debug"}}`)
	ParseConfig(path)

	err := Save(func(c *GlobalConfig) {
		c.Database = &DatabaseConfig{Addr: "root@tcp(127.0.0.1:3306)/halo", Driver: "mysql"}
		c.Secret = &SecretConfig{
			Jwt:    &JwtConfig{TokenName: "Authentication", Issuer: "halo", SecretKey: "secret"},
			Crypto: "0123456789abcdef",
		}
	})
	require.NoError(t, err)

	body, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(body), `"listen": ":8443"`)
	require.Contains(t, string(body), "root@tcp(127.0.0.1:3306)/halo")

	// 内存里的配置也要跟着更新，服务重启前就能用。
	require.True(t, Configured())
	require.Equal(t, ":8443", Config().Listen)
}

func TestSaveKeepsFilePermissions(t *testing.T) {
	path := writeConfig(t, `{"listen": ":8443"}`)
	ParseConfig(path)

	require.NoError(t, Save(func(c *GlobalConfig) {
		c.Database = &DatabaseConfig{Addr: "root@tcp(127.0.0.1:3306)/halo"}
	}))

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestSaveReplacesFileAtomically(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"listen": ":8443"}`), 0o600))
	ParseConfig(path)

	require.NoError(t, Save(func(c *GlobalConfig) {
		c.Listen = ":9443"
	}))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1, "临时文件不应残留")
	require.Equal(t, ":9443", Config().Listen)
}
