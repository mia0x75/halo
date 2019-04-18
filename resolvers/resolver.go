package resolvers

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
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
		credential := tools.Credential{}
		ok := false

		credential, ok = ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		if !ok {
			payload := handler.GetInitPayload(ctx)
			if payload == nil {
				rc = gqlapi.ReturnCodeUnauthorized
				err = fmt.Errorf("错误代码: %s, 错误信息: 没有提供访问令牌、访问令牌被收回或已过期，访问被拒绝。", rc)
				break
			}
			token := payload[g.Config().Secret.Jwt.TokenName].(string)
			credential, err = tools.New().ParseToken(token)
			if err != nil {
				rc = gqlapi.ReturnCodeUnauthorized
				err = fmt.Errorf("错误代码: %s, 错误信息: 没有提供访问令牌、访问令牌被收回或已过期，访问被拒绝。", rc)
				break
			}
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
type Resolver struct{}
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

// RequestCols TODO: 有问题
func (r *Resolver) RequestCols(fields []graphql.CollectedField, model interface{}) (columns []string) {
	// 从model的成员和底层映射表字段的对照关系
	table := g.Engine.TableInfo(model)
	cols := table.Columns()
	for _, col := range cols {
		for _, field := range fields {
			if field.Name == col.FieldName {
				columns = append(columns, col.Name)
			}
		}
	}

	return
}

func updateEdges(edgeType uint, userID uint, edges []*models.Edge) (err error) {
	session := g.Engine.NewSession()
	defer session.Close()
	session.Begin()

	edge := models.Edge{
		Type:       uint(edgeType),
		AncestorID: userID,
	}

	for {
		if _, err = session.Delete(&edge); err != nil {
			session.Rollback()
			break
		}
		if _, err = session.Insert(&edges); err != nil {
			session.Rollback()
			break
		}
		if err = session.Commit(); err != nil {
			break
		}
		break
	}
	return
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

func isClusterAvalaible() bool {
	return true
}

func isUserAvalaible() bool {
	return true
}
