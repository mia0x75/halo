package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

var token *JWT     // TODO: 唯一实例，名称不太好
var once sync.Once // 单例模式

// TODO: 结构体名称不太好
type JWT struct {
	issuer   string
	secret   string
	duration time.Duration
	cache    sync.Map
}

type Credential struct {
	User  *models.User
	Roles []*models.Role
}

type Claims struct {
	jwt.StandardClaims
	UserUUID string `json:"uuid"`
}

// TODO: New这个函数名称不好，容易造成误解
func New() *JWT {
	once.Do(func() {
		token = &JWT{}
	})
	return token
}

func (this *JWT) Load(key string) (string, bool) {
	if value, ok := this.cache.Load(key); ok {
		return value.(string), ok
	} else {
		return "", ok
	}
}

func (this *JWT) Store(key string, value string) {
	this.cache.Store(key, Sha1(value))
}

func (this *JWT) Delete(key string) {
	this.cache.Delete(key)
}

func (this *JWT) CreateToken(user *models.User) (token string, err error) {
	claims := &Claims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 24).Unix()),
			Issuer:    this.issuer,
		},
		user.UUID,
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(this.secret))
	this.Store(user.UUID, token)
	return
}

func (this *JWT) FromContext(ctx context.Context) *Credential {
	value := ctx.Value(g.CREDENTIAL_KEY)
	if value == nil {
		return nil
	}
	credential := ctx.Value(g.CREDENTIAL_KEY).(Credential)
	return &credential
}

// ParseToken 解析令牌，生成凭证
func (this *JWT) ParseToken(input string) (cred Credential, err error) {
	for {
		token := &jwt.Token{}
		token, err = jwt.ParseWithClaims(input, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(this.secret), nil
		})
		if err != nil {
			// 无效的令牌
			err = fmt.Errorf("Invalid token!")
			break
		}
		claims := &Claims{}
		ok := false
		if claims, ok = token.Claims.(*Claims); !ok || !token.Valid {
			// 无效的令牌
			err = fmt.Errorf("Invalid token")
			break
		}

		// 不存在的令牌不处理
		if value, ok := this.Load(claims.UserUUID); !ok {
			// 找不到
			err = fmt.Errorf("Cannot find the token from token storage")
			break
		} else {
			// 找到，比较令牌的签名
			// 令牌的键存在，但值不同，说明重新登录过，现在提供的是废弃的令牌
			if !strings.EqualFold(value, Sha1(input)) {
				err = fmt.Errorf("The token is obsolete, a new token is generated on server side, typically it is because of user login from another device")
				break
			}
		}

		// 验证通过
		UUID := claims.UserUUID
		// 必须要找到用户
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == UUID {
				return true
			}
			return false
		})
		if user == nil {
			// 令牌中记录的用户ID在系统中不存在
			err = fmt.Errorf("The token owner does not exists in system")
			break
		}

		// 用户必须要有角色关联
		edges := caches.EdgesMap.Filter(func(elem *models.Edge) bool {
			if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole] &&
				elem.AncestorID == user.UserID {
				return true
			}
			return false
		})
		if edges == nil {
			// 用户没有角色关联
			err = fmt.Errorf("There's no role binding the user")
			break
		}
		roles := caches.RolesMap.Filter(func(elem *models.Role) bool {
			for _, r := range edges {
				if elem.RoleID == r.DescendantID {
					return true
				}
			}
			return false
		})
		if roles == nil {
			// 用户没有角色关联
			err = fmt.Errorf("There's no role binding the user")
			break
		}
		cred = Credential{
			User:  user,
			Roles: roles,
		}

		break
	}

	return
}
