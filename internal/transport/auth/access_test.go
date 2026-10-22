package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

func TestClusterAccessible(t *testing.T) {
	user := &models.User{UserID: 1}
	role := &models.Role{RoleID: 2, Name: "Developer"}
	cluster := &models.Cluster{ClusterID: 10, GroupID: 5}
	ungrouped := &models.Cluster{ClusterID: 11, GroupID: 0}

	direct := &models.Edge{
		Type:         uint(domain.EdgeTypeUserToCluster),
		AncestorID:   1,
		DescendantID: 10,
	}
	groupGrant := &models.Edge{
		Type:         uint(domain.EdgeTypeRoleToGroup),
		AncestorID:   2,
		DescendantID: 5,
	}

	require.True(t, ClusterAccessible(user, []*models.Role{role}, cluster, []*models.Edge{direct}))
	require.True(t, ClusterAccessible(user, []*models.Role{role}, cluster, []*models.Edge{groupGrant}))
	require.False(t, ClusterAccessible(user, nil, cluster, []*models.Edge{groupGrant}))
	require.False(t, ClusterAccessible(user, []*models.Role{role}, cluster, nil))
	require.False(t, ClusterAccessible(user, []*models.Role{role}, ungrouped, []*models.Edge{groupGrant}))
	require.False(t, ClusterAccessible(nil, []*models.Role{role}, cluster, []*models.Edge{direct}))
}
