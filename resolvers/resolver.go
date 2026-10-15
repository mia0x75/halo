package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/tools"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/models"
)

// TicketSub 为每个订阅的用户分配一个通道，所以用UserUUID做键就可以了
var TicketSub struct {
	sync.RWMutex
	Subscribers map[string]chan *gqlapi.TicketStatusChangePayload
}

// SubscriptionAuth WebSocket的鉴权
func SubscriptionAuth(ctx context.Context, requires []gqlapi.RoleEnum) (user *models.User, err error) {
	for {
		rc := gqlapi.ReturnCodeOK
		credential := transportauth.Credential{}
		ok := false

		credential, ok = ctx.Value(g.CREDENTIAL_KEY).(transportauth.Credential)
		if !ok {
			payload := transport.GetInitPayload(ctx)
			if payload == nil {
				rc = gqlapi.ReturnCodeUnauthorized
				err = fmt.Errorf("错误代码: %s, 错误信息: 没有提供访问令牌、访问令牌被收回或已过期，访问被拒绝。", rc)
				break
			}
			token := payload[g.Config().Secret.Jwt.TokenName].(string)
			cred, err := transportauth.Authenticate(token)
			if err != nil {
				rc = gqlapi.ReturnCodeUnauthorized
				err = fmt.Errorf("错误代码: %s, 错误信息: 没有提供访问令牌、访问令牌被收回或已过期，访问被拒绝。", rc)
				break
			}
			credential = *cred
		}
		user = credential.User
		// TODO: 此处由于gqlgen支持不好，代码跟Directive.Auth代码雷同
		switch user.Status {
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal]:
			// 有令牌，但是找不到匹配的角色
			rc = gqlapi.ReturnCodeForbidden
			err = fmt.Errorf("错误代码: %s, 错误信息: 权限不足，访问被拒绝。", rc)
			requireRoleIDs := []uint{} // requireRoleIds
			for _, require := range requires {
				requireRoleIDs = append(requireRoleIDs, gqlapi.RoleEnumMap[require])
			}
			for _, role := range credential.Roles {
				if tools.Contains(requireRoleIDs, role.RoleID) {
					rc = gqlapi.ReturnCodeOK
					err = nil
					break
				}
			}

			if err != nil {
				user = nil
			}
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending]:
			rc = gqlapi.ReturnCodeUserStatusPending
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)当前状态是等待验证。", rc, user.UUID)
			user = nil
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked]:
			rc = gqlapi.ReturnCodeUserStatusBlocked
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)已经被禁用。", rc, user.UUID)
			user = nil
		default:
			rc = gqlapi.ReturnCodeUserStatusUnknown
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)当前状态是未知。", rc, user.UUID)
			user = nil
		}

		break
	}

	return
}

func init() {
	TicketSub.Lock()
	defer TicketSub.Unlock()
	TicketSub.Subscribers = make(map[string]chan *gqlapi.TicketStatusChangePayload, 0)
}

// Resolver resolver
type Resolver struct {
	app *app.App
}

// NewResolver returns a resolver wired to the application container. During the
// 2D migration the container runs alongside the legacy engine.

func NewResolver(a *app.App) *Resolver {
	return &Resolver{app: a}
}

type mutationRootResolver struct{ *Resolver }
type queryRootResolver struct{ *Resolver }
type subscriptionRootResolver struct{ *Resolver }

// Comment TODO: 添加描述
func (r *Resolver) Comment() gqlapi.CommentResolver {
	return &commentResolver{r}
}

// Log TODO: 添加描述
func (r *Resolver) Log() gqlapi.LogResolver {
	return &logResolver{r}
}

// MutationRoot TODO: 添加描述
func (r *Resolver) MutationRoot() gqlapi.MutationRootResolver {
	return &mutationRootResolver{r}
}

// Query TODO: 添加描述
func (r *Resolver) Query() gqlapi.QueryResolver {
	return &queryResolver{r}
}

// QueryRoot TODO: 添加描述
func (r *Resolver) QueryRoot() gqlapi.QueryRootResolver {
	return &queryRootResolver{r}
}

// Role TODO: 添加描述
func (r *Resolver) Role() gqlapi.RoleResolver {
	return &roleResolver{r}
}

// Statement TODO: 添加描述
func (r *Resolver) Statement() gqlapi.StatementResolver {
	return &statementResolver{r}
}

// SubscriptionRoot TODO: 添加描述
func (r *Resolver) SubscriptionRoot() gqlapi.SubscriptionRootResolver {
	return &subscriptionRootResolver{r}
}

// Ticket TODO: 添加描述
func (r *Resolver) Ticket() gqlapi.TicketResolver {
	return &ticketResolver{r}
}

// User TODO: 添加描述
func (r *Resolver) User() gqlapi.UserResolver {
	return &userResolver{r}
}

// EncodeCursor 对分页的光标进行编码
func EncodeCursor(s string) string {
	i, _ := strconv.Atoi(s)
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i)))
}

// updateEdges replaces every edge of one kind leaving userID, in one tx.
func (r *Resolver) updateEdges(ctx context.Context, edgeType uint, userID uint, edges []*models.Edge) error {
	domainEdges := make([]*domain.Edge, 0, len(edges))
	for _, e := range edges {
		domainEdges = append(domainEdges, bridge.EdgeToDomain(e))
	}
	return r.app.Services.Edge.Replace(ctx, edgeType, userID, domainEdges, userID)
}

func notifyTicketStatusChange(userUUID string, ticketUUID string, message string) {
	TicketSub.Lock()
	if sub, ok := TicketSub.Subscribers[userUUID]; ok {
		payload := &gqlapi.TicketStatusChangePayload{
			TicketUUID: ticketUUID,
			Message:    message,
		}
		sub <- payload
	}
	TicketSub.Unlock()
}

// statDatabase verifies the database exists on the target cluster. It mirrors
// the removed models.Cluster.Stat existence check.
func statDatabase(ctx context.Context, exec agent.Executor, target *domain.Cluster, password, database string) error {
	databases, err := exec.Databases(ctx, target, password)
	if err != nil {
		return err
	}
	for _, db := range databases {
		if db.Name == database {
			return nil
		}
	}
	return fmt.Errorf("unknown database '%s'", database)
}
