// Package authz maps the GraphQL @auth enums onto seeded role names.
//
// The mapping lives here instead of internal/transport/auth because it is part
// of the GraphQL contract: the generic auth helpers only know role names.
package authz

import (
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// WireRoleNames maps each @auth requirement to the roles that satisfy it.
var WireRoleNames = map[gqlapi.RoleEnum][]string{
	gqlapi.RoleEnumAdmin:     {auth.RoleRoot},
	gqlapi.RoleEnumReviewer:  {auth.RoleReviewer},
	gqlapi.RoleEnumDeveloper: {auth.RoleDeveloper},
	gqlapi.RoleEnumUser:      {auth.RoleUser, auth.RoleViewer},
	gqlapi.RoleEnumGuest:     {auth.RoleGuest},
}

// HasRole reports whether one of the roles satisfies the requirement.
func HasRole(roles []*models.Role, required gqlapi.RoleEnum) bool {
	return auth.HasAnyRoleName(roles, WireRoleNames[required])
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
