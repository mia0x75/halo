package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/models"
)

func TestHasRoleName(t *testing.T) {
	root := &models.Role{RoleID: 1, Name: "Root"}
	reviewer := &models.Role{RoleID: 2, Name: "Reviewer"}
	developer := &models.Role{RoleID: 3, Name: "Developer"}
	viewer := &models.Role{RoleID: 4, Name: "Viewer"}
	user := &models.Role{RoleID: 5, Name: "User"}
	guest := &models.Role{RoleID: 6, Name: "Guest"}

	require.True(t, HasRoleName([]*models.Role{root}, RoleRoot))
	require.True(t, HasRoleName([]*models.Role{reviewer}, RoleReviewer))
	require.True(t, HasRoleName([]*models.Role{developer}, RoleDeveloper))
	require.True(t, HasRoleName([]*models.Role{viewer}, RoleViewer))
	require.True(t, HasRoleName([]*models.Role{user}, RoleUser))
	require.True(t, HasRoleName([]*models.Role{guest}, RoleGuest))
	// Names are compared case-insensitively and trimmed.
	require.True(t, HasRoleName([]*models.Role{{Name: "  root "}}, RoleRoot))

	require.False(t, HasRoleName([]*models.Role{guest}, RoleUser))
	require.False(t, HasRoleName(nil, RoleRoot))
}

func TestHasAnyRoleName(t *testing.T) {
	viewer := &models.Role{RoleID: 4, Name: "Viewer"}

	require.True(t, HasAnyRoleName([]*models.Role{viewer}, []string{RoleRoot, RoleViewer}))
	require.False(t, HasAnyRoleName([]*models.Role{viewer}, []string{RoleRoot, RoleReviewer}))
}

func TestIsAdmin(t *testing.T) {
	require.True(t, IsAdmin([]*models.Role{{RoleID: 1, Name: "Root"}}))
	require.False(t, IsAdmin([]*models.Role{{RoleID: 2, Name: "Reviewer"}}))
}

func TestKnownRoleNamesMatchSeed(t *testing.T) {
	require.ElementsMatch(t, []string{
		"Root", "Reviewer", "Developer", "Viewer", "User", "Guest",
	}, KnownRoleNames())
}
