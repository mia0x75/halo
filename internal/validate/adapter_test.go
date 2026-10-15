package validate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestAdapterNotImplemented(t *testing.T) {
	a := NewAdapter(nil, nil)
	_, err := a.Validate(context.Background(), domain.ValidateRequest{})
	require.ErrorIs(t, err, ErrNotImplemented)
}
