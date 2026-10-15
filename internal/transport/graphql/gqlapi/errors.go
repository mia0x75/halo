package gqlapi

import "strconv"

// ReturnCode API返回代码自定义类型
type ReturnCode float64

// String 错误代码转String方法
func (rc ReturnCode) String() string {
	return strconv.FormatFloat(float64(rc), 'f', -1, 64)
}

// API返回代码
var (
	// TODO: 错误代码梳理
	ReturnCodeOK                     ReturnCode = 1200 // 通用 ─┬─ 请求成功
	ReturnCodeInvalidParams          ReturnCode = 1400 //       ├─ 参数错误
	ReturnCodeUnauthorized           ReturnCode = 1401 //       ├─ 未经授权或令牌失效
	ReturnCodeForbidden              ReturnCode = 1403 //       ├─ 权限不足
	ReturnCodeNotFound               ReturnCode = 1404 //       ├─ 资源不存在
	ReturnCodeTimeout                ReturnCode = 1408 //       ├─ 系统响应超时
	ReturnCodeUnknowError            ReturnCode = 1500 //       └─ 未知错误,一般是系统异常
	ReturnCodeUserStatusPending      ReturnCode = 2005 // 专用 ─┬─
	ReturnCodeUserEmailTaken         ReturnCode = 2009 //       ├─
	ReturnCodeWrongPassword          ReturnCode = 2010 //       ├─
	ReturnCodeUserStatusBlocked      ReturnCode = 2016 //       ├─
	ReturnCodeUserStatusUnknown      ReturnCode = 2026 //       ├─
	ReturnCodeUserNotAvailable       ReturnCode = 2026 //       ├─
	ReturnCodeNoEdgeToReviewer       ReturnCode = 2026 //       ├─ // TODO: 名字不太好
	ReturnCodeRegistrationIncomplete ReturnCode = 2001 //       ├─
	ReturnCodeEmailPasswordMismatch  ReturnCode = 2000 //       ├─
	ReturnCodeClusterNotAvailable    ReturnCode = 2001 //       ├─
	ReturnCodeNotWritable            ReturnCode = 2001 //       ├─ // TODO: 名字不太好
	ReturnCodeDuplicateAlias         ReturnCode = 2003 //       ├─ // TODO: 名字不太好
	ReturnCodeDuplicateHost          ReturnCode = 2006 //       └─ // TODO: 名字不太好
)
