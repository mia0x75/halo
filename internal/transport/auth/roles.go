package auth

import (
	"strings"

	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// roleNames maps the wire enum to seeded role names. The enum integers are
// historical and do not line up with mm_roles.role_id, so permissions are
// resolved by name instead of by id.
var roleNames = map[gqlapi.RoleEnum][]string{
	gqlapi.RoleEnumAdmin:     {"Root", "Admin"},
	gqlapi.RoleEnumReviewer:  {"Reviewer"},
	gqlapi.RoleEnumDeveloper: {"Developer"},
	gqlapi.RoleEnumUser:      {"User", "Viewer"},
	gqlapi.RoleEnumGuest:     {"Guest"},
}

// HasRole reports whether one of the roles satisfies the required enum.
func HasRole(roles []*models.Role, required gqlapi.RoleEnum) bool {
	names := roleNames[required]
	for _, role := range roles {
		for _, name := range names {
			if strings.EqualFold(strings.TrimSpace(role.Name), name) {
				return true
			}
		}
	}
	return false
}

// HasAnyRole reports whether one of the roles satisfies any requirement.
func HasAnyRole(roles []*models.Role, required []gqlapi.RoleEnum) bool {
	for _, item := range required {
		if HasRole(roles, item) {
			return true
		}
	}
	return false
}

// IsAdmin reports whether the roles include the platform administrator.
func IsAdmin(roles []*models.Role) bool {
	return HasRole(roles, gqlapi.RoleEnumAdmin)
}

// FindRoleByName looks up a seeded role by its display name.
func FindRoleByName(name string) *models.Role {
	return cache.RolesMap.Any(func(elem *models.Role) bool {
		return strings.EqualFold(strings.TrimSpace(elem.Name), name)
	})
}
