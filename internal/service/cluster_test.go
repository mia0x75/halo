package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/domain"
)

func TestClusterServiceTestConnectionSucceeds(t *testing.T) {
	ctx := context.Background()
	exec := &fakeExecutor{}
	svc := NewClusterService(&fakeTx{}, &fakeClusterRepo{}, exec)

	result, err := svc.TestConnection(ctx, &domain.Cluster{IP: "127.0.0.1", Port: 3306}, "secret")
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, "8.0.36", result.Version)
	require.Equal(t, []string{"app"}, result.Databases)
}

func TestClusterServiceTestConnectionReportsFailure(t *testing.T) {
	ctx := context.Background()
	exec := &fakeExecutor{connErr: agent.ErrNotImplemented}
	svc := NewClusterService(&fakeTx{}, &fakeClusterRepo{}, exec)

	result, err := svc.TestConnection(ctx, &domain.Cluster{IP: "127.0.0.1", Port: 3306}, "secret")
	require.NoError(t, err)
	require.False(t, result.Success)
	require.NotEmpty(t, result.Message)
}
