package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

func TestClusterGrants(t *testing.T) {
	user := &models.User{UserID: 1}
	role := &models.Role{RoleID: 2, Name: "Developer"}
	clusters := []*models.Cluster{
		{ClusterID: 10, GroupID: 5},
		{ClusterID: 11, GroupID: 0},
	}
	edges := []*models.Edge{
		{Type: gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster], AncestorID: 1, DescendantID: 10},
		{Type: gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup], AncestorID: 2, DescendantID: 5},
	}
	groupName := func(groupID uint) string {
		if groupID == 5 {
			return "核心交易"
		}
		return ""
	}

	grants := ClusterGrants(user, []*models.Role{role}, clusters, edges, groupName)
	require.Len(t, grants, 2)

	require.True(t, grants[0].Direct)
	require.True(t, grants[0].Inherited)
	require.Equal(t, []string{"Developer / 核心交易"}, grants[0].Sources)

	require.False(t, grants[1].Direct)
	require.False(t, grants[1].Inherited)
	require.Empty(t, grants[1].Sources)
}

func TestClusterGrantsIgnoresOtherUsersEdges(t *testing.T) {
	user := &models.User{UserID: 1}
	clusters := []*models.Cluster{{ClusterID: 10, GroupID: 0}}
	edges := []*models.Edge{
		{Type: gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster], AncestorID: 99, DescendantID: 10},
	}

	grants := ClusterGrants(user, nil, clusters, edges, func(uint) string { return "" })
	require.Len(t, grants, 1)
	require.False(t, grants[0].Direct)
	require.False(t, grants[0].Inherited)
}
