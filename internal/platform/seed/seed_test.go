package seed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipDropsImportOnlyStatements(t *testing.T) {
	require.True(t, skip("LOCK TABLES `mm_roles` WRITE"))
	require.True(t, skip("UNLOCK TABLES;"))
	require.True(t, skip("SET NAMES utf8mb4"))
	require.True(t, skip("set time_zone = '+00:00'"))
	require.False(t, skip("INSERT INTO mm_roles VALUES (1, 'x')"))
}

func TestApplyRejectsUnknownDriver(t *testing.T) {
	err := Apply(t.Context(), nil, "sqlite")
	require.ErrorContains(t, err, "元数据库类型")
}
