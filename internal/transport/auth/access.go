package auth

import (
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// ClusterAccessible reports whether the user may operate on the cluster.
//
// Access comes from two sources: a direct user-to-cluster grant, or a role of
// the user being granted the group the cluster belongs to. The second one is
// what makes the group model a real isolation boundary instead of a label.
func ClusterAccessible(user *models.User, roles []*models.Role, cluster *models.Cluster, edges []*models.Edge) bool {
	if user == nil || cluster == nil {
		return false
	}
	for _, edge := range edges {
		if edge.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
			edge.AncestorID == user.UserID &&
			edge.DescendantID == cluster.ClusterID {
			return true
		}
	}
	if cluster.GroupID == 0 || len(roles) == 0 {
		return false
	}
	roleIDs := make(map[uint]bool, len(roles))
	for _, role := range roles {
		roleIDs[role.RoleID] = true
	}
	for _, edge := range edges {
		if edge.Type != gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup] {
			continue
		}
		if edge.DescendantID != cluster.GroupID {
			continue
		}
		if roleIDs[edge.AncestorID] {
			return true
		}
	}
	return false
}

// CanAccessCluster is ClusterAccessible against the current edge cache.
func CanAccessCluster(user *models.User, roles []*models.Role, cluster *models.Cluster) bool {
	return ClusterAccessible(user, roles, cluster, cache.EdgesMap.All())
}
