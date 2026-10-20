package resolvers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscriptionTokenAcceptsConfiguredName(t *testing.T) {
	require.Equal(t, "abc", subscriptionToken(map[string]any{"Authentication": "abc"}))
}

func TestSubscriptionTokenAcceptsAuthorizationHeader(t *testing.T) {
	require.Equal(t, "abc", subscriptionToken(map[string]any{"Authorization": "Bearer abc"}))
	require.Equal(t, "abc", subscriptionToken(map[string]any{"Authorization": "abc"}))
}

func TestSubscriptionTokenIgnoresNonString(t *testing.T) {
	require.Equal(t, "", subscriptionToken(map[string]any{"Authentication": 42}))
}

func TestSubscriptionTokenMissing(t *testing.T) {
	require.Equal(t, "", subscriptionToken(map[string]any{}))
	require.Equal(t, "", subscriptionToken(map[string]any{"Authentication": "  "}))
}
