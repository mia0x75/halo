package resolvers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentialErrorsHideSecrets(t *testing.T) {
	require.Contains(t, loginMismatchError("user@example.com").Error(), "password=********")
	require.Contains(t, loginMismatchError("user@example.com").Error(), "user@example.com")
	require.NotContains(t, samePasswordError().Error(), "password=")
	require.NotContains(t, wrongOldPasswordError().Error(), "password=")
}
