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

// ClusterGrant explains why a user may use a cluster.
//
// Permissions are the union of a direct grant and the groups granted to the
// user's roles. The management page has to show both layers, otherwise
// "unchecking a cluster" looks like it revokes access when it only drops one
// of the two sources.
type ClusterGrant struct {
	Cluster *models.Cluster
	// Direct reports a user-to-cluster grant.
	Direct bool
	// Inherited reports access through a group granted to one of the roles.
	Inherited bool
	// Sources names the roles and groups that grant inherited access.
	Sources []string
}

// ClusterGrants computes the effective access of one user over every cluster.
func ClusterGrants(
	user *models.User,
	roles []*models.Role,
	clusters []*models.Cluster,
	edges []*models.Edge,
	groupName func(groupID uint) string,
) []ClusterGrant {
	grants := make([]ClusterGrant, 0, len(clusters))
	if user == nil {
		return grants
	}

	roleIDs := make(map[uint]bool, len(roles))
	roleNames := make(map[uint]string, len(roles))
	for _, role := range roles {
		roleIDs[role.RoleID] = true
		roleNames[role.RoleID] = role.Name
	}

	direct := make(map[uint]bool)
	for _, edge := range edges {
		if edge.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] && edge.AncestorID == user.UserID {
			direct[edge.DescendantID] = true
		}
	}

	// groupID -> 角色名列表，先按分组聚合继承来源。
	inherited := make(map[uint][]string)
	for _, edge := range edges {
		if edge.Type != gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup] {
			continue
		}
		name, ok := roleNames[edge.AncestorID]
		if !ok {
			continue
		}
		source := name
		if group := groupName(edge.DescendantID); group != "" {
			source = name + " / " + group
		}
		inherited[edge.DescendantID] = append(inherited[edge.DescendantID], source)
	}

	for _, cluster := range clusters {
		sources := inherited[cluster.GroupID]
		grants = append(grants, ClusterGrant{
			Cluster:   cluster,
			Direct:    direct[cluster.ClusterID],
			Inherited: len(sources) > 0,
			Sources:   sources,
		})
	}
	return grants
}
