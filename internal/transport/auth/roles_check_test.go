package auth

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/models"
)

func TestValidateRolesAcceptsSeededNames(t *testing.T) {
	roles := []*models.Role{
		{RoleID: 1, Name: "Root"},
		{RoleID: 2, Name: "Reviewer"},
		{RoleID: 3, Name: "Developer"},
		{RoleID: 4, Name: "Viewer"},
		{RoleID: 5, Name: "User"},
		{RoleID: 6, Name: "Guest"},
	}
	require.NoError(t, ValidateRoles(roles))
}

func TestValidateRolesRejectsUnknownName(t *testing.T) {
	err := ValidateRoles([]*models.Role{{RoleID: 1, Name: "Superuser"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "Superuser")
}

func TestValidateRolesRejectsEmptyName(t *testing.T) {
	err := ValidateRoles([]*models.Role{{RoleID: 1, Name: "  "}})
	require.Error(t, err)
}
