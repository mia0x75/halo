package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func numberOption(value string) *Option {
	lower, upper := 1, 3650
	return &Option{Name: "retention.logs-days", Element: OptionElementNumber, Parameter: value, MinValue: &lower, MaxValue: &upper}
}

func TestValidateOptionValue(t *testing.T) {
	switchOption := &Option{Name: "mail.enabled", Element: OptionElementSwitch, Parameter: "true"}
	require.NoError(t, ValidateOptionValue(switchOption))
	switchOption.Parameter = "yes"
	require.ErrorIs(t, ValidateOptionValue(switchOption), ErrInvalidInput)

	require.NoError(t, ValidateOptionValue(numberOption("180")))
	require.ErrorIs(t, ValidateOptionValue(numberOption("abc")), ErrInvalidInput)
	require.ErrorIs(t, ValidateOptionValue(numberOption("0")), ErrInvalidInput)
	require.ErrorIs(t, ValidateOptionValue(numberOption("9999")), ErrInvalidInput)

	// 没有范围时不做上下界检查。
	require.NoError(t, ValidateOptionValue(&Option{Name: "x", Element: OptionElementNumber, Parameter: "9999"}))
}

func TestMaskOptionValue(t *testing.T) {
	require.Equal(t, "********", MaskOptionValue(&Option{Name: "mail.password", Element: OptionElementPassword, Parameter: "secret"}))
	require.Equal(t, "", MaskOptionValue(&Option{Name: "mail.password", Element: OptionElementPassword}))
	require.Equal(t, "90", MaskOptionValue(&Option{Name: "retention.logs-days", Element: OptionElementNumber, Parameter: "90"}))
}

func TestOptionSecretAndWritable(t *testing.T) {
	require.True(t, (&Option{Element: OptionElementPassword}).Secret())
	require.False(t, (&Option{Element: OptionElementInput}).Secret())
	require.True(t, (&Option{Writable: 1}).IsWritable())
	require.False(t, (&Option{Writable: 0}).IsWritable())
}
