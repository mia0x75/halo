package mysqlrepo

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestClusterRepo(t *testing.T) {
	repo := NewClusterRepo(setupTestDB(t))
	ctx := context.Background()

	c := &domain.Cluster{
		UUID:        newUUID(),
		Host:        "host-" + newSuffix(),
		Alias:       "alias-" + newSuffix(),
		IP:          newIP(),
		Port:        newPort(),
		User:        "root",
		Type:        domain.ClusterTypeMySQL,
		Password:    []byte("secret"),
		Fingerprint: []byte("fp"),
		Status:      domain.ClusterStatusNormal,
		Version:     1,
		CreateAt:    nowUnix(),
	}
	groupID := newID()
	c.GroupID = groupID
	id, err := repo.Create(ctx, c)
	require.NoError(t, err)
	require.NotZero(t, id)
	c.ClusterID = uint(id)

	count, err := repo.CountByGroup(ctx, groupID)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	got, err := repo.GetByUUID(ctx, c.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.ClusterTypeMySQL, got.Type)
	require.Equal(t, []byte("secret"), got.Password)
	require.Equal(t, c.Port, got.Port)

	byID, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, c.UUID, byID.UUID)

	updated := nowUnix()
	c.Alias = "alias-updated-" + newSuffix()
	c.Status = domain.ClusterStatusDisabled
	c.Version = 2
	c.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, c))

	got, err = repo.GetByUUID(ctx, c.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.ClusterStatusDisabled, got.Status)
	require.Equal(t, c.Alias, got.Alias)

	list, err := repo.List(ctx, int64(math.MaxUint32), 100)
	require.NoError(t, err)
	require.True(t, hasCluster(list, c.UUID))

	require.NoError(t, repo.DeleteByUUID(ctx, c.UUID))
	_, err = repo.GetByUUID(ctx, c.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasCluster(list []*domain.Cluster, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
