package validate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestRuleInt(t *testing.T) {
	r := ruleFromDomain(&domain.Rule{Name: "n", Func: "f", Values: "42"})
	got, err := r.Int()
	require.NoError(t, err)
	require.Equal(t, 42, got)

	// memoised: a second read returns the cached value
	got, err = r.Int()
	require.NoError(t, err)
	require.Equal(t, 42, got)

	_, err = ruleFromDomain(&domain.Rule{Name: "n", Values: "not-an-int"}).Int()
	require.Error(t, err)

	got, err = ruleFromDomain(&domain.Rule{Name: "n"}).Int()
	require.NoError(t, err)
	require.Zero(t, got)

	got, err = ruleFromDomain(&domain.Rule{Name: "n", Values: "nil"}).Int()
	require.NoError(t, err)
	require.Zero(t, got)
}

func TestRuleStringSlice(t *testing.T) {
	r := ruleFromDomain(&domain.Rule{Name: "n", Values: `["a","b"]`})
	got, err := r.StringSlice()
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b"}, got)

	_, err = ruleFromDomain(&domain.Rule{Name: "n", Values: `not json`}).StringSlice()
	require.Error(t, err)

	got, err = ruleFromDomain(&domain.Rule{Name: "n"}).StringSlice()
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestRuleRegexp(t *testing.T) {
	r := ruleFromDomain(&domain.Rule{Name: "n", Values: `^[a-z]+$`})
	re, err := r.Regexp()
	require.NoError(t, err)
	require.True(t, re.MatchString("abc"))
	require.False(t, re.MatchString("ABC"))

	_, err = ruleFromDomain(&domain.Rule{Name: "n", Values: `(`}).Regexp()
	require.Error(t, err)

	re, err = ruleFromDomain(&domain.Rule{Name: "n"}).Regexp()
	require.NoError(t, err)
	require.Nil(t, re)
}

func TestValidateRules(t *testing.T) {
	ok := []*domain.Rule{
		{Name: "a", Func: "f", Operator: "in", Values: `["x"]`},
		{Name: "b", Func: "g", Operator: "regexp", Values: `^x$`},
	}
	require.NoError(t, ValidateRules(ok))

	badJSON := []*domain.Rule{{Name: "c", Func: "f", Operator: "in", Values: `oops`}}
	require.Error(t, ValidateRules(badJSON))

	badRegexp := []*domain.Rule{{Name: "d", Func: "f", Operator: "regexp", Values: `(`}}
	require.Error(t, ValidateRules(badRegexp))

	emptyFunc := []*domain.Rule{{Name: "e", Operator: "in", Values: `["x"]`}}
	require.Error(t, ValidateRules(emptyFunc))

	dup := []*domain.Rule{
		{Name: "f", Func: "g", Operator: "in", Values: `["x"]`},
		{Name: "f", Func: "g", Operator: "in", Values: `["x"]`},
	}
	require.Error(t, ValidateRules(dup))

	require.NoError(t, ValidateRules(nil))
}
