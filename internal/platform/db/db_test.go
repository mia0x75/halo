package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectDriver(t *testing.T) {
	cases := []struct {
		addr string
		want string
	}{
		{"postgres://halo:pw@10.0.0.7:5432/halo?sslmode=disable", "postgres"},
		{"postgresql://halo:pw@10.0.0.7:5432/halo", "postgres"},
		{"  POSTGRES://u@h/db  ", "postgres"},
		{"mysql://halo:pw@127.0.0.1:3306/halo", "mysql"},
		{"halo:pw@tcp(127.0.0.1:3306)/halo?charset=utf8mb4&parseTime=true", "mysql"},
		{"halo:pw@/halo", "mysql"},
		{"halo:pw@unix(/run/mysqld/mysqld.sock)/halo", "mysql"},
		{"", "mysql"},
	}
	for _, c := range cases {
		got, err := DetectDriver(c.addr)
		require.NoError(t, err, c.addr)
		require.Equal(t, c.want, got, c.addr)
	}

	_, err := DetectDriver("oracle://u@h/db")
	require.Error(t, err)
}

// The conflict and unknown-driver errors fire before any connection is
// opened, so they need no live database.
func TestNewRejectsDriverMismatch(t *testing.T) {
	ctx := context.Background()
	pgAddr := "postgres://halo:pw@10.0.0.7:5432/halo?sslmode=disable"
	myAddr := "halo:pw@tcp(127.0.0.1:3306)/halo"

	_, err := New(ctx, Config{Addr: pgAddr, Driver: "mysql"})
	require.ErrorContains(t, err, "conflicts")

	_, err = New(ctx, Config{Addr: myAddr, Driver: "postgres"})
	require.ErrorContains(t, err, "conflicts")

	_, err = New(ctx, Config{Addr: myAddr, Driver: "oracle"})
	require.ErrorContains(t, err, "unknown metadata driver")
}
