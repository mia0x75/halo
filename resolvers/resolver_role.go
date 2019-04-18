package resolvers

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

// Role 查看某一个角色信息
func (r *queryRootResolver) Role(ctx context.Context, id string) (role *models.Role, err error) {
	rc := gqlapi.ReturnCodeOK
	role = caches.RolesMap.Any(func(elem *models.Role) bool {
		if elem.UUID == id {
			return true
		}
		return false
	})
	if role == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 角色(uuid=%s)不存在。", rc, id)
	}
	return
}

// Roles 因为角色不会太多，所以不需要考虑数据分页问题
func (r *queryRootResolver) Roles(ctx context.Context) (L []*models.Role, err error) {
	L = caches.RolesMap.All()
	return
}

type roleResolver struct{ *Resolver }

// Users 查看角色关联的用户信息
func (r *roleResolver) Users(ctx context.Context, obj *models.Role, after *string, before *string, first *int, last *int) (*gqlapi.UserConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	// TODO: 从光标解出来用户ID，然后根据ID查询，前提缓存数据需要排序
	edges := caches.EdgesMap.Filter(func(elem *models.Edge) bool {
		if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole] &&
			elem.DescendantID == obj.RoleID {
			return true
		}
		return false
	})

	if edges == nil {
		return &gqlapi.UserConnection{}, nil
	}

	users := caches.UsersMap.Filter(func(elem *models.User) bool {
		for _, r := range edges {
			if elem.UserID == r.AncestorID {
				return true
			}
		}
		return false
	})

	if users == nil {
		return &gqlapi.UserConnection{}, nil
	}

	userEdges := []gqlapi.UserEdge{}
	for _, user := range users {
		userEdges = append(userEdges, gqlapi.UserEdge{
			Node:   user,
			Cursor: EncodeCursor(fmt.Sprintf("%d", user.UserID)),
		})
	}
	if len(userEdges) == 0 {
		return &gqlapi.UserConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     "",
		EndCursor:       "",
	}

	return &gqlapi.UserConnection{
		PageInfo:   pageInfo,
		Edges:      userEdges,
		TotalCount: len(userEdges),
	}, nil
}
