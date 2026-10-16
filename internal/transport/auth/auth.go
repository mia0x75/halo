package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/mia0x75/halo/internal/auth"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Credential is the resolved session principal attached to a request context.
type Credential struct {
	User  *models.User
	Roles []*models.Role
}

type authenticator struct {
	jwt   *auth.JWT
	store sync.Map
}

var authn = &authenticator{jwt: auth.NewDefault()}

// Issue signs a token for userUUID and records its hash.
func Issue(userUUID string) (string, error) {
	token, err := authn.jwt.CreateToken(userUUID)
	if err != nil {
		return "", err
	}
	authn.store.Store(userUUID, sha1hex(token))
	return token, nil
}

// Authenticate verifies a token and resolves the bound user and roles.
func Authenticate(input string) (*Credential, error) {
	uuid, err := authn.jwt.ParseToken(input)
	if err != nil {
		return nil, err
	}
	v, ok := authn.store.Load(uuid)
	if !ok {
		return nil, fmt.Errorf("cannot find the token from token storage")
	}
	if !strings.EqualFold(v.(string), sha1hex(input)) {
		return nil, fmt.Errorf("token is obsolete (a newer one was issued)")
	}

	user := cache.UsersMap.Any(func(elem *models.User) bool { return elem.UUID == uuid })
	if user == nil {
		return nil, fmt.Errorf("the token owner does not exist in system")
	}

	edges := cache.EdgesMap.Filter(func(elem *models.Edge) bool {
		return elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole] &&
			elem.AncestorID == user.UserID
	})
	if edges == nil {
		return nil, fmt.Errorf("there's no role binding the user")
	}

	roles := cache.RolesMap.Filter(func(elem *models.Role) bool {
		for _, r := range edges {
			if elem.RoleID == r.DescendantID {
				return true
			}
		}
		return false
	})
	if roles == nil {
		return nil, fmt.Errorf("there's no role binding the user")
	}

	return &Credential{User: user, Roles: roles}, nil
}

// Revoke drops the stored token for the given user UUID.
func Revoke(userUUID string) {
	authn.store.Delete(userUUID)
}

// FromContext reads the credential stored by the HTTP middleware.
func FromContext(ctx context.Context) *Credential {
	value := ctx.Value(config.CREDENTIAL_KEY)
	if value == nil {
		return nil
	}
	if credential, ok := value.(Credential); ok {
		return &credential
	}
	return nil
}

func sha1hex(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
