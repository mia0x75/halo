package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateOptionValue(t *testing.T) {
	require.NoError(t, ValidateOptionValue("smtp.enabled", "true"))
	require.ErrorIs(t, ValidateOptionValue("smtp.enabled", "yes"), ErrInvalidInput)

	require.NoError(t, ValidateOptionValue("retention.logs-days", "180"))
	require.ErrorIs(t, ValidateOptionValue("retention.logs-days", "abc"), ErrInvalidInput)
	require.ErrorIs(t, ValidateOptionValue("retention.logs-days", "0"), ErrInvalidInput)
	require.ErrorIs(t, ValidateOptionValue("retention.logs-days", "9999"), ErrInvalidInput)
}

func TestMaskOptionValue(t *testing.T) {
	require.Equal(t, "********", MaskOptionValue("smtp.password", "secret"))
	require.Equal(t, "", MaskOptionValue("smtp.password", ""))
	require.Equal(t, "90", MaskOptionValue("retention.logs-days", "90"))
}
