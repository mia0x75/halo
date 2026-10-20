package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Groups 群集分组列表
func (r *queryRootResolver) Groups(ctx context.Context) (L []*models.Group, err error) {
	L = cache.GroupsMap.All()
	return
}

// Group 群集所属分组的名称
func (r *clusterResolver) Group(ctx context.Context, obj *models.Cluster) (*string, error) {
	if obj.GroupID == 0 {
		return nil, nil
	}
	group := cache.GroupsMap.Any(func(elem *models.Group) bool {
		return elem.GroupID == obj.GroupID
	})
	if group == nil {
		return nil, nil
	}
	name := group.Name
	return &name, nil
}

// Groups 角色可管理的群集分组
func (r *roleResolver) Groups(ctx context.Context, obj *models.Role) (L []*models.Group, err error) {
	edgeType := gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup]
	for _, edge := range cache.EdgesMap.All() {
		if edge.Type != edgeType || edge.AncestorID != obj.RoleID {
			continue
		}
		group := cache.GroupsMap.Any(func(elem *models.Group) bool {
			return elem.GroupID == edge.DescendantID
		})
		if group == nil {
			continue
		}
		L = append(L, group)
	}
	return
}

// GrantGroups 关联角色到群集分组
func (r *mutationRootResolver) GrantGroups(ctx context.Context, input gqlapi.GrantGroupsInput) (ok bool, err error) {
	role := cache.RolesMap.Any(func(elem *models.Role) bool {
		return elem.UUID == input.RoleUUID
	})
	if role == nil {
		return false, fmt.Errorf("错误代码: %s, 错误信息: 角色(uuid=%s)不存在。", gqlapi.ReturnCodeNotFound, input.RoleUUID)
	}

	edges := make([]*models.Edge, 0, len(input.GroupUUIDs))
	for _, groupUUID := range input.GroupUUIDs {
		group := cache.GroupsMap.Any(func(elem *models.Group) bool {
			return elem.UUID == groupUUID
		})
		if group == nil {
			return false, fmt.Errorf("错误代码: %s, 错误信息: 分组(uuid=%s)不存在。", gqlapi.ReturnCodeNotFound, groupUUID)
		}
		edges = append(edges, &models.Edge{
			Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup],
			AncestorID:   role.RoleID,
			DescendantID: group.GroupID,
		})
	}

	domainEdges := make([]*domain.Edge, 0, len(edges))
	for _, edge := range edges {
		domainEdges = append(domainEdges, bridge.EdgeToDomain(edge))
	}

	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	edgeType := gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup]
	if err = r.app.Services.Edge.Replace(ctx, edgeType, role.RoleID, domainEdges, credential.User.UserID); err != nil {
		return false, fmt.Errorf("错误代码: %s, 错误信息: %s", gqlapi.ReturnCodeUnknowError, err.Error())
	}

	cache.EdgesMap.Reload()
	return true, nil
}

// RevokeGroups 收回角色到群集分组的关联
func (r *mutationRootResolver) RevokeGroups(ctx context.Context, input gqlapi.GrantGroupsInput) (ok bool, err error) {
	role := cache.RolesMap.Any(func(elem *models.Role) bool {
		return elem.UUID == input.RoleUUID
	})
	if role == nil {
		return false, fmt.Errorf("错误代码: %s, 错误信息: 角色(uuid=%s)不存在。", gqlapi.ReturnCodeNotFound, input.RoleUUID)
	}

	edgeType := gqlapi.EdgeEnumMap[gqlapi.EdgeEnumRoleToGroup]
	credential := ctx.Value(config.CREDENTIAL_KEY).(transportauth.Credential)
	for _, edge := range cache.EdgesMap.All() {
		if edge.Type != edgeType || edge.AncestorID != role.RoleID {
			continue
		}
		if err = r.app.Services.Edge.Delete(ctx, edge.EdgeID, credential.User.UserID); err != nil {
			return false, fmt.Errorf("错误代码: %s, 错误信息: %s", gqlapi.ReturnCodeUnknowError, err.Error())
		}
	}

	cache.EdgesMap.Reload()
	return true, nil
}
