package auth

import (
	"fmt"
	"sort"
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

// KnownRoleNames returns every role name the permission mapping understands.
func KnownRoleNames() []string {
	names := make([]string, 0, 8)
	for _, group := range roleNames {
		names = append(names, group...)
	}
	sort.Strings(names)
	return names
}

// ValidateRoles rejects role names the permission mapping does not know.
//
// Roles are immutable (there is no create or delete entry point), and a role's
// name carries its permissions. Renaming one in the database would silently
// strip every permission it grants, so the server refuses to start instead.
func ValidateRoles(roles []*models.Role) error {
	known := make(map[string]bool, 8)
	for _, name := range KnownRoleNames() {
		known[strings.ToLower(name)] = true
	}
	var unknown []string
	for _, role := range roles {
		name := strings.TrimSpace(role.Name)
		if name == "" || !known[strings.ToLower(name)] {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	return fmt.Errorf("未知角色名 %s：可用角色为 %s",
		strings.Join(unknown, ", "), strings.Join(KnownRoleNames(), ", "))
}
