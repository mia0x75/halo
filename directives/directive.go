package directives

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/graphql"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/tools"
)

// Auth 权限验证的directive
func Auth(ctx context.Context, obj interface{}, next graphql.Resolver, requires []gqlapi.RoleEnum) (res interface{}, err error) {
	rc := gqlapi.ReturnCodeOK

	if rctx := graphql.GetResolverContext(ctx); obj != rctx.Parent.Result {
		// TODO: 处理rc
		return nil, fmt.Errorf("错误代码: %s, 错误信息: 父对象类型不匹配。", rc)
	}

	// TODO: 考虑给tokens增加一个方法，从Context中返回*tokens.Credential
	if value := ctx.Value(g.CREDENTIAL_KEY); value != nil {
		if credential, ok := value.(tools.Credential); ok {
			switch credential.User.Status {
			case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumNormal]:
				// 有令牌，但是找不到匹配的角色
				rc = gqlapi.ReturnCodeForbidden
				err = fmt.Errorf("错误代码: %s, 错误信息: 权限不足，访问被拒绝。", rc)
				for _, require := range requires {
					for _, role := range credential.Roles {
						if gqlapi.RoleEnumMap[require] == role.RoleID {
							rc = gqlapi.ReturnCodeOK
							err = nil
							return next(ctx)
						}
					}
				}
			// 用户状态不允许，通常是用户登录后，令牌过期前，管理员操作了用户的状态
			case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumPending]:
				rc = gqlapi.ReturnCodeUserStatusPending
				err = fmt.Errorf("错误代码: %s, 错误信息: 用户(uuid=%s)当前状态是等待验证。", rc, credential.User.UUID)
			// 用户状态不允许，通常是用户登录后，令牌过期前，管理员操作了用户的状态
			case gqlapi.UserStatusEnumMap[gqlapi.UserStatusEnumBlocked]:
				rc = gqlapi.ReturnCodeUserStatusBlocked
				err = fmt.Errorf("错误代码: %s, 错误信息: 账号(uuid=%s)已经被禁用。", rc, credential.User.UUID)
			}
		} else {
			rc = gqlapi.ReturnCodeUnknowError
			err = fmt.Errorf("错误代码: %s, 错误信息: 未知错误，无效的访问凭证(value=%v)。", rc, value)
		}
	} else {
		// 没有令牌
		rc = gqlapi.ReturnCodeUnauthorized
		err = fmt.Errorf("错误代码: %s, 错误信息: 没有提供访问令牌、访问令牌被收回或已过期，访问被拒绝。", rc)
	}

	return nil, err
}

// Date 日期格式化的directive
func Date(ctx context.Context, obj interface{}, next graphql.Resolver, format string) (res interface{}, err error) {
	return next(ctx)
}

// EnumInt 枚举转数字的directive
func EnumInt(ctx context.Context, obj interface{}, next graphql.Resolver, value int) (res interface{}, err error) {
	return value, nil
}

// Length 字符串长度限制的directive
func Length(ctx context.Context, obj interface{}, next graphql.Resolver, max int) (res interface{}, err error) {
	switch obj.(type) {
	case string:
		src := obj.(string)
		if len(src) > max {
			return nil, fmt.Errorf("字符串`%s`长度已经超过允许的最大值(%d)。", src, max)
		}
	}
	return next(ctx)
}

// Lower 字符串转小写的directive
func Lower(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	switch obj.(type) {
	case string:
		return strings.ToLower(obj.(string)), nil
	}
	return next(ctx)
}

// Matches 字符串正则匹配的directive
func Matches(ctx context.Context, obj interface{}, next graphql.Resolver, pattern string) (res interface{}, err error) {
	re := &regexp.Regexp{}
	if re, err = regexp.Compile(pattern); err != nil {
		return
	}
	switch obj.(type) {
	case string:
		if !re.MatchString(obj.(string)) {
			return nil, fmt.Errorf("The given string `%s` does not match against the regexp pattern `%s`", obj.(string), pattern)
		}
	}
	return next(ctx)
}

// Range 数字范围限制的directive
func Range(ctx context.Context, obj interface{}, next graphql.Resolver, begin int, end int) (res interface{}, err error) {
	msg := "The int value `%d` is out of range, a value between the ranges of: %d to %d is accepted."
	switch obj.(type) {
	case int, int8, int16, int32, int64:
		value := obj.(int64)
		if value < int64(begin) || value > int64(end) {
			return nil, fmt.Errorf(msg, value, begin, end)
		}
	case uint, uint8, uint16, uint32, uint64:
		value := obj.(uint64)
		if value < uint64(begin) || value > uint64(end) {
			return nil, fmt.Errorf(msg, value, begin, end)
		}
	}
	return next(ctx)
}

// Rename 改名的directive
func Rename(ctx context.Context, obj interface{}, next graphql.Resolver, to string) (res interface{}, err error) {
	return next(ctx)
}

// Timestamp 事件戳的directive
func Timestamp(ctx context.Context, obj interface{}, next graphql.Resolver, format string) (res interface{}, err error) {
	return next(ctx)
}

// Trim 字符串去首尾空白字符的directive
func Trim(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return next(ctx)
}

// Upper 字符串转大写的directive
func Upper(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return next(ctx)
}

// Uuid UUID的directive
func Uuid(ctx context.Context, obj interface{}, next graphql.Resolver, name *string, from []*string) (res interface{}, err error) {
	return next(ctx)
}
