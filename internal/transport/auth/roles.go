package auth

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/models"
)

// Seeded role names. Permissions are resolved by name because the wire enum
// integers do not line up with mm_roles.role_id.
const (
	RoleRoot      = "Root"
	RoleReviewer  = "Reviewer"
	RoleDeveloper = "Developer"
	RoleUser      = "User"
	RoleViewer    = "Viewer"
	RoleGuest     = "Guest"
)

// AdminRoleNames carry the ADMIN requirement. The seed data has a single
// administrator role, Root; there is no "Admin" role in mm_roles.
var AdminRoleNames = []string{RoleRoot}

// allRoleNames is every role name the permission mapping understands.
var allRoleNames = []string{
	RoleRoot, RoleReviewer, RoleDeveloper, RoleUser, RoleViewer, RoleGuest,
}

// HasRoleName reports whether the roles include one with the given name.
func HasRoleName(roles []*models.Role, name string) bool {
	for _, role := range roles {
		if strings.EqualFold(strings.TrimSpace(role.Name), name) {
			return true
		}
	}
	return false
}

// HasAnyRoleName reports whether the roles include any of the given names.
func HasAnyRoleName(roles []*models.Role, names []string) bool {
	for _, name := range names {
		if HasRoleName(roles, name) {
			return true
		}
	}
	return false
}

// IsAdmin reports whether the roles include the platform administrator.
func IsAdmin(roles []*models.Role) bool {
	return HasAnyRoleName(roles, AdminRoleNames)
}

// FindRoleByName looks up a seeded role by its display name.
func FindRoleByName(name string) *models.Role {
	return cache.RolesMap.Any(func(elem *models.Role) bool {
		return strings.EqualFold(strings.TrimSpace(elem.Name), strings.TrimSpace(name))
	})
}

// KnownRoleNames returns every role name the permission mapping understands.
func KnownRoleNames() []string {
	names := append([]string(nil), allRoleNames...)
	sort.Strings(names)
	return names
}

// ValidateRoles rejects role names the permission mapping does not know.
//
// Roles are immutable (there is no create or delete entry point), and a role's
// name carries its permissions. Renaming one in the database would silently
// strip every permission it grants, so the server refuses to start instead.
func ValidateRoles(roles []*models.Role) error {
	known := make(map[string]bool, len(allRoleNames))
	for _, name := range allRoleNames {
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
