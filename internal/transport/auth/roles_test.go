package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

func TestHasRole(t *testing.T) {
	root := &models.Role{RoleID: 1, Name: "Root"}
	reviewer := &models.Role{RoleID: 2, Name: "Reviewer"}
	developer := &models.Role{RoleID: 3, Name: "Developer"}
	viewer := &models.Role{RoleID: 4, Name: "Viewer"}
	user := &models.Role{RoleID: 5, Name: "User"}
	guest := &models.Role{RoleID: 6, Name: "Guest"}

	require.True(t, HasRole([]*models.Role{root}, gqlapi.RoleEnumAdmin))
	require.True(t, HasRole([]*models.Role{reviewer}, gqlapi.RoleEnumReviewer))
	require.True(t, HasRole([]*models.Role{developer}, gqlapi.RoleEnumDeveloper))
	require.True(t, HasRole([]*models.Role{viewer}, gqlapi.RoleEnumUser))
	require.True(t, HasRole([]*models.Role{user}, gqlapi.RoleEnumUser))
	require.True(t, HasRole([]*models.Role{guest}, gqlapi.RoleEnumGuest))

	require.False(t, HasRole([]*models.Role{guest}, gqlapi.RoleEnumUser))
	require.False(t, HasRole([]*models.Role{viewer}, gqlapi.RoleEnumDeveloper))
	require.False(t, HasRole(nil, gqlapi.RoleEnumAdmin))
}

func TestHasAnyRole(t *testing.T) {
	viewer := &models.Role{RoleID: 4, Name: "Viewer"}

	require.True(t, HasAnyRole([]*models.Role{viewer}, []gqlapi.RoleEnum{gqlapi.RoleEnumAdmin, gqlapi.RoleEnumUser}))
	require.False(t, HasAnyRole([]*models.Role{viewer}, []gqlapi.RoleEnum{gqlapi.RoleEnumAdmin, gqlapi.RoleEnumReviewer}))
}

func TestIsAdmin(t *testing.T) {
	require.True(t, IsAdmin([]*models.Role{{RoleID: 1, Name: "Root"}}))
	require.False(t, IsAdmin([]*models.Role{{RoleID: 2, Name: "Reviewer"}}))
}
