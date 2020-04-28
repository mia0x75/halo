package resolvers

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mia0x75/yql"
	"golang.org/x/crypto/bcrypt"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/events"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// Me 当前登陆用户信息
func (r queryRootResolver) Me(ctx context.Context) (user *models.User, err error) {
	// Directive.Auth确保了一定可以获得正确的凭证
	credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
	user = credential.User
	return
}

// User 根据ID查看用户详细信息
func (r queryRootResolver) User(ctx context.Context, id string) (user *models.User, err error) {
	rc := gqlapi.ReturnCodeOK
	user = caches.UsersMap.Any(func(elem *models.User) bool {
		if elem.UUID == id {
			return true
		}
		return false
	})
	if user == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)不存在。", rc, id)
	}
	return
}

// Users 分页查看用户列表
func (r queryRootResolver) Users(ctx context.Context, after *string, before *string, first *int, last *int) (*gqlapi.UserConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	from := math.MaxInt64
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return nil, err
		}
		from = i
	}
	hasPreviousPage := true
	hasNextPage := true

	if from == math.MaxInt64 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []*gqlapi.UserEdge{}

	// TODO: 一次从数据库查出足够的数据
	// TODO: 对于分页中下一页的处理：
	//       如果不需要记录总数或者说总页数，假设页大小为n，那么总是查n+1是否可行
	//       如果需要查记录总数或者说总页数，假设页大小为n，那么用offset+n和记录总数判断是否可行
	//       在Edge.Cursor中包含offset是否可行
	users := []*models.User{}
	var err error
	if err = g.Engine.Desc("user_id").Where("user_id < ?", from).Limit(*first).Find(&users); err != nil {
		return nil, err
	}

	for _, user := range users {
		edges = append(edges, &gqlapi.UserEdge{
			Node:   user,
			Cursor: EncodeCursor(fmt.Sprintf("%d", user.UserID)),
		})
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return nil, nil
	}

	// 获取pageInfo
	startCursor := EncodeCursor(fmt.Sprintf("%d", users[0].UserID))
	endCursor := EncodeCursor(fmt.Sprintf("%d", users[len(users)-1].UserID))
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     startCursor,
		EndCursor:       endCursor,
	}

	return &gqlapi.UserConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// UserSearch 用户搜索，TODO: 分页处理
func (r queryRootResolver) UserSearch(ctx context.Context, search string, after *string, before *string, first *int, last *int) (rs *gqlapi.UserConnection, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		// 参数判断，只允许 first/before first/after last/before last/after 模式
		if first != nil && last != nil {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
			break L
		}
		if after != nil && before != nil {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
			break L
		}

		var rule yql.Ruler
		if rule, err = yql.Rule(search); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 语法错误，请参考帮助文档。", rc)
			break L
		}
		users := caches.UsersMap.Filter(func(elem *models.User) bool {
			// yql最好可以使用反射，直接操作struct对象，而不是map[string]interface{}
			data := structs.Map(elem)
			if matched, _ := rule.Match(data); matched {
				return true
			}
			return false
		})
		count := len(users)
		if count == 0 {
			break L
		}
		// 获取edges
		edges := []*gqlapi.UserEdge{}
		for _, user := range users {
			edges = append(edges, &gqlapi.UserEdge{
				Node:   user,
				Cursor: EncodeCursor(strconv.Itoa(0)),
			})
		}
		pageInfo := &gqlapi.PageInfo{
			StartCursor:     edges[0].Cursor,
			EndCursor:       edges[count-1].Cursor,
			HasPreviousPage: false,
			HasNextPage:     false,
		}

		rs = &gqlapi.UserConnection{
			PageInfo:   pageInfo,
			Edges:      edges,
			TotalCount: count,
		}

		break L
	}

	return
}

// Register 用户注册
func (r mutationRootResolver) Register(ctx context.Context, input models.UserRegisterInput) (user *models.User, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		user = caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.Email == input.Email {
				return true
			}
			return false
		})

		if user != nil {
			switch user.Status {
			case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending]:
				rc = gqlapi.ReturnCodeUserStatusPending
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)当前状态是等待验证。", rc, user.Email)
			case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal]:
				rc = gqlapi.ReturnCodeUserEmailTaken
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)已经被注册。", rc, input.Email)
			case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked]:
				rc = gqlapi.ReturnCodeUserStatusBlocked
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)已经被禁用。", rc, user.Email)
			default:
				rc = gqlapi.ReturnCodeUserStatusUnknown
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)当前状态是未知。", rc, user.Email)
			}
			break L
		}

		bs, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user = &models.User{
			Email:    input.Email,
			Password: string(bs),
			AvatarID: 1,
			Status:   gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending],
		}
		session := g.Engine.NewSession()
		defer session.Close()
		session.Begin()
		if _, err = session.Insert(user); err != nil {
			session.Rollback()
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		edge := &models.Edge{
			Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole], // TODO: 消除硬编码
			AncestorID:   user.UserID,                                   //
			DescendantID: gqlapi.RoleEnumMap["USER"],                    // 默认给USER权限
		}
		if _, err = session.Insert(edge); err != nil {
			session.Rollback()
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		if err = session.Commit(); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 主动同步缓存
		caches.UsersMap.Append(user)
		caches.EdgesMap.Reload()

		events.Fire(events.EventUserRegistered, &events.UserRegisteredArgs{
			User: *user,
		})

		// 退出for循环
		break
	}

	return
}

// Activate 账号激活
func (r mutationRootResolver) Activate(ctx context.Context, input models.ActivateInput) (payload *gqlapi.ActivatePayload, err error) {
L:
	// 1, 解密
	// 2, 数据有效性判断
	// 3, 是否过期
	// 4, 用户是否存在，状态是否正确
	for {
		var bs []byte
		bs, err = hex.DecodeString(input.Code)
		if err != nil {
			break L
		}
		data := map[string]string{}

		bs, err = tools.DecryptAES(bs, g.Config().Secret.Crypto)
		if err != nil {
			break L
		}

		err = json.Unmarshal(bs, &data)
		if err != nil {
			err = fmt.Errorf("")
			break L
		}

		var email, expire string
		if value, ok := data["Email"]; ok {
			email = value
			payload = &gqlapi.ActivatePayload{
				Email: email,
			}

		} else {
			// TODO:
			err = fmt.Errorf("")
			break L
		}
		if value, ok := data["Expire"]; ok {
			expire = value
		} else {
			// TODO:
			err = fmt.Errorf("")
			break L
		}

		var dt time.Time
		dt, err = time.Parse("2006-01-02 15:04:05", expire)
		if err != nil {
			err = fmt.Errorf("")
			break L
		}

		diff := time.Now().Sub(dt)
		if int(diff.Seconds()) > 24*60*60 {
			err = fmt.Errorf("")
			break L
		}

		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.Email == email {
				return true
			}
			return false
		})

		if user == nil {
			// TODO:
			err = fmt.Errorf("")
			break L
		}

		if user.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending] {
			err = fmt.Errorf("")
			break L
		}

		user.Status = gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal]

		if _, err = g.Engine.ID(user.UserID).Update(user); err != nil {
			err = fmt.Errorf("")
			break L
		}

		events.Fire(events.EventUserRegistered, &events.UserRegisteredArgs{
			User: *user,
		})

		break L
	}

	return
}

// LostPasswd 忘记密码
func (r mutationRootResolver) LostPasswd(ctx context.Context, input models.LostPasswdInput) (ok bool, err error) {
L:
	for {
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.Email == input.Email {
				return true
			}
			return false
		})

		if user == nil {
			// TODO:
			err = fmt.Errorf("...%s", "")
			break L
		}

		switch user.Status {
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal]:
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending]:
			// TODO:
			err = fmt.Errorf("...%s", "")
			break L
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked]:
			// TODO:
			err = fmt.Errorf("...%s", "")
			break L
		default:
			// TODO:
			err = fmt.Errorf("...%s", "")
			break L
		}

		events.Fire(events.EventUserRegistered, &events.UserRegisteredArgs{
			User: *user,
		})

		break L
	}

	return
}

// ResetPasswd 重置密码
func (r mutationRootResolver) ResetPasswd(ctx context.Context, input models.ResetPasswdInput) (ok bool, err error) {
L:
	for {
		// TODO:
		events.Fire(events.EventUserRegistered, &events.UserRegisteredArgs{
			// UserUUID: user.UUID,
		})

		break L
	}
	return
}

// ResendActivationMail 重发验证激活邮件
func (r mutationRootResolver) ResendActivationMail(ctx context.Context, input models.ActivateInput) (ok bool, err error) {
L:
	for {
		// 1, 参考Activate 1/2/3/4，不同是没有过期则不允许
		// TODO:
		events.Fire(events.EventUserRegistered, &events.UserRegisteredArgs{
			// UserUUID: user.UUID,
		})

		break L
	}
	return
}

// Login 用户登陆
func (r mutationRootResolver) Login(ctx context.Context, input models.UserLoginInput) (payload *gqlapi.LoginPayload, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK

		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.Email == input.Email {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeEmailPasswordMismatch
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)和密码(password=%s)不匹配。", rc, input.Email, input.Password)
			break L
		}
		switch user.Status {
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending]:
			rc = gqlapi.ReturnCodeUserStatusPending
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)当前状态是等待验证。", rc, user.Email)
			break L
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal]:
		case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked]:
			rc = gqlapi.ReturnCodeUserStatusBlocked
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)已经被禁用。", rc, user.Email)
			break L
		default:
			rc = gqlapi.ReturnCodeUserStatusUnknown
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)当前状态是未知。", rc, user.Email)
			break L
		}
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			rc = gqlapi.ReturnCodeEmailPasswordMismatch
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)和密码(password=%s)不匹配。", rc, input.Email, input.Password)
			break L
		}

		var token string
		token, err = func(user *models.User) (string, error) {
			token := tools.New()
			return token.CreateToken(user)
		}(user)

		if err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		payload = &gqlapi.LoginPayload{
			Me:    user,
			Token: token,
		}

		events.Fire(events.EventUserSignedIn, &events.UserSignedInArgs{
			User: *user,
			IP:   "127.0.0.1", // TODO:
		})

		// 退出for循环
		break
	}

	if err != nil {
		payload = nil
	}

	return
}

// Logout 退出登陆
func (r mutationRootResolver) Logout(ctx context.Context) (bool, error) {
	// 通过上下文获得用户凭证，此处凭证经过auth的检测，必然是有效的
	credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
	user := credential.User
	tools.New().Delete(user.UUID)

	events.Fire(events.EventUserLogout, &events.UserLogoutArgs{
		User: *user,
	})

	return true, nil
}

// UpdateProfile 用户自行修改个人信息
func (r mutationRootResolver) UpdateProfile(ctx context.Context, input models.UpdateProfileInput) (user *models.User, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user = credential.User

		avatar := caches.AvatarsMap.Any(func(elem *models.Avatar) bool {
			if elem.UUID == input.AvatarUUID {
				return true
			}
			return false
		})
		if avatar != nil {
			user.AvatarID = avatar.AvatarID
		} else {
			user.AvatarID = 1
		}
		user.Name = strings.TrimSpace(input.Name)
		user.Phone = input.Phone

		if _, err = g.Engine.ID(user.UserID).Update(user); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		events.Fire(events.EventProfileUpdated, &events.ProfileUpdatedArgs{
			User: *user,
		})

		break
	}

	if err != nil {
		user = nil
	}

	return
}

// UpdatePassword 用户自行修改密码
func (r mutationRootResolver) UpdatePassword(ctx context.Context, input models.PatchPasswordInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		// 通过上下文获得用户凭证，此处凭证经过auth的检测，必然是有效的
		if strings.EqualFold(strings.TrimSpace(input.OldPassword), strings.TrimSpace(input.NewPassword)) {
			// TODO:
			// rc = gqlapi.
			err = fmt.Errorf("错误代码: %s, 错误信息: 新密码(password=%s)和旧密码(password=%s)相同。", rc, input.NewPassword, input.OldPassword)
			break L
		}

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
			rc = gqlapi.ReturnCodeWrongPassword
			err = fmt.Errorf("错误代码: %s, 错误信息: 旧密码(password=%s)不正确。", rc, input.OldPassword)
			break L
		}
		bs, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		credential.User.Password = string(bs)
		if _, err = g.Engine.ID(user.UserID).Update(user); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		events.Fire(events.EventPasswordUpdated, &events.PasswordUpdatedArgs{
			User: *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// UpdateEmail 用户自行更新邮件（登陆账号）
func (r mutationRootResolver) UpdateEmail(ctx context.Context, input models.PatchEmailInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		// 通过上下文获得用户凭证，此处凭证经过auth的检测，必然是有效的
		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		user := credential.User

		if user.Email == input.NewEmail {
			// TODO: 梳理错误代码
			// rc = g.ReturnCodeExitEmailError
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)无需修改。", rc, input.NewEmail)
			break L
		}

		// 此处用户必然提供了一个不同的邮件地址
		if caches.UsersMap.Include(func(elem *models.User) bool {
			if elem.Email == input.NewEmail {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeUserEmailTaken
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)已经存在。", rc, input.NewEmail)
			break L
		}

		user.Email = input.NewEmail
		// TODO: 需要单元测试
		if _, err = g.Engine.ID(user.UserID).Where("`email` <> ?", input.NewEmail).Update(user); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		events.Fire(events.EventEmailUpdated, &events.EmailUpdatedArgs{
			User:  *user,
			Email: input.NewEmail,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// GrantReviewers 关联用户到审核人
func (r mutationRootResolver) GrantReviewers(ctx context.Context, input models.GrantReviewersInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		edges := []*models.Edge{}
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		for _, reviewerUUID := range input.ReviewerUUIDs {
			reviewer := caches.UsersMap.Any(func(elem *models.User) bool {
				if elem.UUID == reviewerUUID {
					return true
				}
				return false
			})
			if reviewer == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, reviewerUUID)
				break L
			}

			if reviewer.Status != gqlapi.UserStatusEnumMap["NORMAL"] {
				rc = gqlapi.ReturnCodeUserNotAvailable
				err = fmt.Errorf("错误代码: %s, 错误信息: 审核用户(uuid=%s)的当前状态异常。", rc, reviewerUUID)
				break L
			}

			r := models.Edge{
				Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer], //
				AncestorID:   user.UserID,                                       //
				DescendantID: reviewer.UserID,                                   //
			}
			edges = append(edges, &r)
		}

		if err = updateEdges(gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer], user.UserID, edges); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 主动刷新缓存
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventReviewerGranted, &events.ReviewerGrantedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// RevokeReviewers 收回用户的审核人
func (r mutationRootResolver) RevokeReviewers(ctx context.Context, input models.RevokeReviewersInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		// TODO: 合并成数组一次删除
		for _, reviewerUUID := range input.ReviewerUUIDs {
			reviewer := caches.UsersMap.Any(func(elem *models.User) bool {
				if elem.UUID == reviewerUUID {
					return true
				}
				return false
			})
			if reviewer == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, reviewerUUID)
				break L
			}
			edge := caches.EdgesMap.Any(func(elem *models.Edge) bool {
				if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer] &&
					elem.AncestorID == user.UserID &&
					elem.DescendantID == reviewer.UserID {
					return true
				}
				return false
			})
			if edge != nil {
				if _, err = g.Engine.ID(edge.EdgeID).Delete(edge); err != nil {
					rc = gqlapi.ReturnCodeUnknowError
					err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
					break L
				}
			}
		}

		// 主动刷新缓存
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventReviewerRevoked, &events.ReviewerRevokedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// GrantClusters 关联用户到群集
func (r mutationRootResolver) GrantClusters(ctx context.Context, input models.GrantClustersInput) (ok bool, err error) {
L:
	for {
		edges := []*models.Edge{}
		rc := gqlapi.ReturnCodeOK
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		for _, clusterUUID := range input.ClusterUUIDs {
			cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
				if elem.UUID == clusterUUID {
					return true
				}
				return false
			})
			if cluster == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, clusterUUID)
				break L
			}

			if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
				rc = gqlapi.ReturnCodeClusterNotAvailable
				err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, clusterUUID)
				break L
			}

			r := models.Edge{
				Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster], //
				AncestorID:   user.UserID,                                      //
				DescendantID: cluster.ClusterID,                                //
			}
			edges = append(edges, &r)
		}

		if err = updateEdges(gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster], user.UserID, edges); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 主动刷新缓存
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventClusterGranted, &events.ClusterGrantedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// RevokeClusters 收回用户的群集授权
func (r mutationRootResolver) RevokeClusters(ctx context.Context, input models.RevokeClustersInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		// TODO: 合并成数组一次删除
		for _, clusterUUID := range input.ClusterUUIDs {
			cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
				if elem.UUID == clusterUUID {
					return true
				}
				return false
			})
			if cluster == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, clusterUUID)
				break L
			}
			edge := caches.EdgesMap.Any(func(elem *models.Edge) bool {
				if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
					elem.AncestorID == user.UserID &&
					elem.DescendantID == cluster.ClusterID {
					return true
				}
				return false
			})
			if edge != nil {
				if _, err = g.Engine.ID(edge.EdgeID).Delete(edge); err != nil {
					rc = gqlapi.ReturnCodeUnknowError
					err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
					break L
				}
			}
		}

		// 主动刷新缓存
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventClusterRevoked, &events.ClusterRevokedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// GrantRoles 关联用户到角色
func (r mutationRootResolver) GrantRoles(ctx context.Context, input models.GrantRolesInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		edges := []*models.Edge{}
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		for _, roleUUID := range input.RoleUUIDs {
			role := caches.RolesMap.Any(func(elem *models.Role) bool {
				if elem.UUID == roleUUID {
					return true
				}
				return false
			})
			if role == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 角色(uuid=%s)不存在。", rc, roleUUID)
				break L
			}
			r := models.Edge{
				Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole], //
				AncestorID:   user.UserID,                                   //
				DescendantID: role.RoleID,                                   //
			}
			edges = append(edges, &r)
		}

		if err = updateEdges(gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole], user.UserID, edges); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 主动刷新缓存
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventRoleGranted, &events.RoleGrantedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// RevokeRoles 收回用户的角色授权
func (r mutationRootResolver) RevokeRoles(ctx context.Context, input models.RevokeRolesInput) (ok bool, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		user := caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})

		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		// TODO: 合并成数组一次删除
		for _, roleUUID := range input.RoleUUIDs {
			role := caches.RolesMap.Any(func(elem *models.Role) bool {
				if elem.UUID == roleUUID {
					return true
				}
				return false
			})
			if role == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 角色(uuid=%s)不存在。", rc, roleUUID)
				break L
			}
			edge := caches.EdgesMap.Any(func(elem *models.Edge) bool {
				if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole] &&
					elem.AncestorID == user.UserID &&
					elem.DescendantID == role.RoleID {
					return true
				}
				return false
			})

			if edge != nil {
				if _, err = g.Engine.ID(edge.EdgeID).Delete(edge); err != nil {
					rc = gqlapi.ReturnCodeUnknowError
					err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
					break L
				}
			}
		}

		// 主动刷新缓存
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventRoleRevoked, &events.RoleRevokedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

// CreateUser 创建一个用户
func (r mutationRootResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (user *models.User, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK

		if caches.UsersMap.Include(func(elem *models.User) bool {
			if elem.Email == input.Email {
				return true
			}
			return false
		}) {
			rc = gqlapi.ReturnCodeUserEmailTaken
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)已经注册。", rc, input.Email)
			break L
		}

		if input.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending] &&
			input.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal] &&
			input.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked] {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 输入参数状态(status=%d)无效。", rc, input.Status)
			break L
		}

		edges := []*models.Edge{}

		for _, roleUUID := range input.RoleUUIDs {
			role := caches.RolesMap.Any(func(elem *models.Role) bool {
				if elem.UUID == roleUUID {
					return true
				}
				return false
			})
			if role == nil {
				rc = gqlapi.ReturnCodeNotFound
				err = fmt.Errorf("错误代码: %s, 错误信息: 角色(uuid=%s)不存在。", rc, roleUUID)
				break L
			}
			edge := &models.Edge{
				Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole],
				AncestorID:   0,
				DescendantID: role.RoleID,
			}
			edges = append(edges, edge)
		}

		// 处理参数可选问题
		if input.ReviewerUUIDs != nil {
			for _, reviewerUUID := range input.ReviewerUUIDs {
				reviewer := caches.UsersMap.Any(func(elem *models.User) bool {
					if elem.UUID == reviewerUUID {
						return true
					}
					return false
				})

				if reviewer == nil {
					rc = gqlapi.ReturnCodeNotFound
					err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, reviewerUUID)
					break L
				}

				if reviewer.Status != gqlapi.UserStatusEnumMap["NORMAL"] {
					rc = gqlapi.ReturnCodeUserNotAvailable
					err = fmt.Errorf("错误代码: %s, 错误信息: 审核用户(uuid=%s)的当前状态异常。", rc, reviewerUUID)
					break L
				}

				edge := &models.Edge{
					Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer],
					AncestorID:   0,
					DescendantID: reviewer.UserID,
				}
				edges = append(edges, edge)
			}
		}

		// 处理参数可选问题
		if input.ClusterUUIDs != nil {
			for _, clusterUUID := range input.ClusterUUIDs {
				cluster := caches.ClustersMap.Any(func(elem *models.Cluster) bool {
					if elem.UUID == clusterUUID {
						return true
					}
					return false
				})
				if cluster == nil {
					rc = gqlapi.ReturnCodeNotFound
					err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不存在。", rc, clusterUUID)
					break L
				}

				if cluster.Status != gqlapi.ClusterStatusEnumMap["NORMAL"] {
					rc = gqlapi.ReturnCodeClusterNotAvailable
					err = fmt.Errorf("错误代码: %s, 错误信息: 群集(uuid=%s)不可用。", rc, clusterUUID)
					break L
				}

				edge := &models.Edge{
					Type:         gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster],
					AncestorID:   0,
					DescendantID: cluster.ClusterID,
				}
				edges = append(edges, edge)
			}
		}

		bs, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user = &models.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: string(bs),
			Status:   input.Status,
			Phone:    input.Phone,
		}
		avatar := caches.AvatarsMap.Any(func(elem *models.Avatar) bool {
			if elem.UUID == input.AvatarUUID {
				return true
			}
			return false
		})
		if avatar != nil {
			user.AvatarID = avatar.AvatarID
		} else {
			user.AvatarID = 1
		}

		// 开启事务
		session := g.Engine.NewSession()
		defer session.Close()
		session.Begin()
		if _, err = session.Insert(user); err != nil {
			session.Rollback()
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 写用户关联信息
		if len(edges) > 0 {
			for _, edge := range edges {
				edge.AncestorID = user.UserID
			}
			if _, err = session.Insert(&edges); err != nil {
				session.Rollback()
				rc = gqlapi.ReturnCodeUnknowError
				err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
				break L
			}
		}

		if err = session.Commit(); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 主动刷新缓存
		caches.UsersMap.Append(user)
		caches.EdgesMap.Reload()

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventUserCreated, &events.UserCreatedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		break
	}

	if err != nil {
		user = nil
	}

	return
}

// UpdateUser 更新用户信息
func (r mutationRootResolver) UpdateUser(ctx context.Context, input models.UpdateUserInput) (user *models.User, err error) {
L:
	for {
		rc := gqlapi.ReturnCodeOK
		user = caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})
		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}

		if input.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending] &&
			input.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal] &&
			input.Status != gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked] {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 输入参数状态(status=%d)无效。", rc, input.Status)
			break L
		}

		// 找不到 => 此次更新了用户登录邮箱
		// 找到，且UserID相同 => 此次没有更新用户邮箱
		// 找到，且UserID不同 => 此次更新了用户登录邮箱，但是该邮箱被其他账号使用
		exists := func(uid uint, email string) bool {
			user := caches.UsersMap.Any(func(elem *models.User) bool {
				if elem.Email == email {
					return true
				}
				return false
			})
			if user != nil && user.UserID != uid {
				return true
			}
			return false
		}(user.UserID, input.Email)

		if exists {
			rc = gqlapi.ReturnCodeUserEmailTaken
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(email=%s)已经注册。", rc, input.Email)
			break L
		}

		bs, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		avatar := caches.AvatarsMap.Any(func(elem *models.Avatar) bool {
			if elem.UUID == input.AvatarUUID {
				return true
			}
			return false
		})
		user.Email = input.Email
		user.Password = string(bs)
		user.Status = input.Status
		user.Name = input.Name
		if avatar == nil {
			user.AvatarID = 1
		} else {
			user.AvatarID = avatar.AvatarID
		}

		if _, err = g.Engine.ID(user.UserID).Update(user); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 后台更新用户，同步清理该用户的令牌，如果有的话
		tools.New().Delete(input.UserUUID)

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		defer events.Fire(events.EventUserUpdated, &events.UserUpdatedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		break L
	}

	return
}

// PatchUserStatus 修改用户状态
func (r mutationRootResolver) PatchUserStatus(ctx context.Context, input models.PatchUserStatusInput) (ok bool, err error) {
L:
	for {
		var user *models.User
		rc := gqlapi.ReturnCodeOK

		if input.Status < 0 || input.Status > 2 {
			rc = gqlapi.ReturnCodeInvalidParams
			err = fmt.Errorf("错误代码: %s, 错误信息: 参数(status=%d)无效。", rc, input.Status)
			break L
		}
		user = caches.UsersMap.Any(func(elem *models.User) bool {
			if elem.UUID == input.UserUUID {
				return true
			}
			return false
		})
		if user == nil {
			rc = gqlapi.ReturnCodeNotFound
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)不存在。", rc, input.UserUUID)
			break L
		}
		// 检查user.Status
		if user.Status == input.Status {
			// rc = g.ReturnCodeModifyUserWarn // TODO: 梳理错误代码
			err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)状态无需更新。", rc, input.UserUUID)
			break L
		}

		user.Status = input.Status
		if _, err = g.Engine.ID(user.UserID).Update(user); err != nil {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: %s", rc, err.Error())
			break L
		}

		// 后台更新用户，同步清理该用户的令牌，如果有的话
		tools.New().Delete(input.UserUUID)

		credential := ctx.Value(g.CREDENTIAL_KEY).(tools.Credential)
		events.Fire(events.EventUserStatusPatched, &events.UserStatusPatchedArgs{
			Manager: *credential.User,
			User:    *user,
		})

		// 退出for循环
		ok = true
		break
	}

	return
}

type userResolver struct{ *Resolver }

// Avatar 用户的头像
func (r userResolver) Avatar(ctx context.Context, obj *models.User) (avatar *models.Avatar, err error) {
	rc := gqlapi.ReturnCodeOK
	avatar = caches.AvatarsMap.Any(func(elem *models.Avatar) bool {
		if elem.AvatarID == obj.AvatarID {
			return true
		}
		return false
	})
	if avatar == nil {
		rc = gqlapi.ReturnCodeNotFound
		err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)的头像不存在。", rc, obj.UUID)
	}
	return
}

// Roles 用户的角色
func (r userResolver) Roles(ctx context.Context, obj *models.User) (L []*models.Role, err error) {
	edges := caches.EdgesMap.Filter(func(elem *models.Edge) bool {
		if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToRole] &&
			elem.AncestorID == obj.UserID {
			return true
		}
		return false
	})

	if edges == nil {
		return
	}

	L = caches.RolesMap.Filter(func(elem *models.Role) bool {
		for _, r := range edges {
			if elem.RoleID == r.DescendantID {
				return true
			}
		}
		return false
	})
	return
}

// Reviewers 用户的审核人
func (r userResolver) Reviewers(ctx context.Context, obj *models.User) (L []*models.User, err error) {
	edges := caches.EdgesMap.Filter(func(elem *models.Edge) bool {
		if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToReviewer] &&
			elem.AncestorID == obj.UserID {
			return true
		}
		return false
	})

	if edges == nil {
		return
	}

	L = caches.UsersMap.Filter(func(elem *models.User) bool {
		for _, r := range edges {
			if elem.UserID == r.DescendantID {
				return true
			}
		}
		return false
	})
	return
}

// Clusters 用户可以访问的群集，TODO: 分页未完成
func (r userResolver) Clusters(ctx context.Context, obj *models.User, after *string, before *string, first *int, last *int) (*gqlapi.ClusterConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return &gqlapi.ClusterConnection{}, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return &gqlapi.ClusterConnection{}, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	edges := caches.EdgesMap.Filter(func(elem *models.Edge) bool {
		if elem.Type == gqlapi.EdgeEnumMap[gqlapi.EdgeEnumUserToCluster] &&
			elem.AncestorID == obj.UserID {
			return true
		}
		return false
	})

	if edges == nil {
		return &gqlapi.ClusterConnection{}, nil
	}

	clusters := caches.ClustersMap.Filter(func(elem *models.Cluster) bool {
		for _, r := range edges {
			if elem.ClusterID == r.DescendantID {
				return true
			}
		}
		return false
	})
	clusterEdges := []*gqlapi.ClusterEdge{}
	for _, cluster := range clusters {
		clusterEdges = append(clusterEdges, &gqlapi.ClusterEdge{
			Node:   cluster,
			Cursor: EncodeCursor(fmt.Sprintf("%d", cluster.ClusterID)),
		})
	}
	if len(clusterEdges) == 0 {
		return &gqlapi.ClusterConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     "",
		EndCursor:       "",
	}

	return &gqlapi.ClusterConnection{
		PageInfo:   pageInfo,
		Edges:      clusterEdges,
		TotalCount: len(edges),
	}, nil
}

// Tickets 用户发起的工单，TODO: 分页未完成
func (r userResolver) Tickets(ctx context.Context, obj *models.User, after *string, before *string, first *int, last *int) (*gqlapi.TicketConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	tickets := []*models.Ticket{}
	if err := g.Engine.Where("user_id = ?", obj.UserID).Find(tickets); err != nil {
		return &gqlapi.TicketConnection{}, err
	}
	edges := []*gqlapi.TicketEdge{}
	for _, ticket := range tickets {
		edges = append(edges, &gqlapi.TicketEdge{
			Node:   ticket,
			Cursor: EncodeCursor(fmt.Sprintf("%d", ticket.TicketID)),
		})
	}
	if len(edges) == 0 {
		return &gqlapi.TicketConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     "",
		EndCursor:       "",
	}

	return &gqlapi.TicketConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// Queries 用户发起的查询，TODO: 分页未完成
func (r userResolver) Queries(ctx context.Context, obj *models.User, after *string, before *string, first *int, last *int) (*gqlapi.QueryConnection, error) {
	rc := gqlapi.ReturnCodeOK
	// 参数判断，只允许 first/before first/after last/before last/after 模式
	if first != nil && last != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`first`和`last`只能选择一种。", rc)
	}
	if after != nil && before != nil {
		rc = gqlapi.ReturnCodeInvalidParams
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 参数`after`和`before`只能选择一种。", rc)
	}

	queries := []*models.Query{}
	g.Engine.Where("user_id = ?", obj.UserID).Find(queries)
	edges := []*gqlapi.QueryEdge{}
	for _, query := range queries {
		edges = append(edges, &gqlapi.QueryEdge{
			Node:   query,
			Cursor: EncodeCursor(fmt.Sprintf("%d", query.QueryID)),
		})
	}
	if len(edges) == 0 {
		return &gqlapi.QueryConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := &gqlapi.PageInfo{
		HasPreviousPage: false,
		HasNextPage:     false,
		StartCursor:     "",
		EndCursor:       "",
	}

	return &gqlapi.QueryConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}

// Statistics 用户维度的统计信息
func (r userResolver) Statistics(ctx context.Context, obj *models.User) (L []*models.Statistic, err error) {
	L = caches.StatisticsMap.Filter(func(elem *models.Statistic) bool {
		if elem.Group == obj.UUID {
			return true
		}
		return false
	})

	return
}
