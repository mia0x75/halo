package gqlapi

// GenderEnumMap 性别枚举转int8
var GenderEnumMap = map[GenderEnum]int8{
	"NA":     -1, // 未知
	"FEMALE": 0,  // 女
	"MALE":   1,  // 男
}

// ClusterStatusEnumMap 群集状态枚举转uint8
var ClusterStatusEnumMap = map[ClusterStatusEnum]uint8{
	"NORMAL":   1, // 正常
	"DISABLED": 2, // 停用
}

// RoleEnumMap 角色枚举转uint
var RoleEnumMap = map[RoleEnum]uint{
	"ADMIN":     1, // 管理员
	"REVIEWER":  2, // 审核人
	"DEVELOPER": 3, // 开发者
	"USER":      4, // 注册用户
	"GUEST":     5, // 访客
}

// StatementTypeEnumMap 语法类型枚举
var StatementTypeEnumMap = map[StatementTypeEnum]uint8{
	"ALTER_TABLE":     1,  //
	"ANALYZE_TABLE":   2,  //
	"BEGIN":           3,  //
	"COMMIT":          4,  //
	"CREATE_DATABASE": 5,  //
	"CREATE_INDEX":    6,  //
	"CREATE_TABLE":    7,  //
	"CREATE_VIEW":     8,  //
	"CREATE_USER":     9,  //
	"DELETE":          10, //
	"DROP_DATABASE":   11, //
	"DROP_INDEX":      12, //
	"DROP_TABLE":      13, //
	"EXPLAIN":         14, //
	"INSERT":          15, //
	"LOAD_DATA":       16, //
	"ROLLBACK":        17, //
	"SELECT":          18, //
	"SET":             19, //
	"SHOW":            20, //
	"TRUNCATE_TABLE":  21, //
	"UPDATE":          22, //
	"GRANT":           23, //
	"REVOKE":          24, //
	"DEALLOCATE":      25, //
	"EXECUTE":         26, //
	"PREPARE":         27, //
	"USE":             28, //
}

// TicketStatusEnumMap 工单状态枚举转uint8
var TicketStatusEnumMap = map[TicketStatusEnum]uint8{
	"WAITING_FOR_VLD": 1, // 等待系统审核
	"VLD_FAILURE":     2, // 系统审核失败
	"VLD_WARNING":     3, // 系统审核警告
	"WAITING_FOR_MRV": 4, // 等待人工审核
	"MRV_FAILURE":     5, // 人工审核失败
	"LGTM":            6, // 人工审核通过
	"DONE":            7, // 上线执行成功
	"EXEC_FAILURE":    8, // 上线执行失败
	"CLOSED":          9, // 工单手工关闭
}

// UserStatusEnumMap 用户状态枚举转uint8
var UserStatusEnumMap = map[UserStatusEnum]uint8{
	"NORMAL":  1, // 正常
	"BLOCKED": 2, // 锁定
	"PENDING": 3, // 等待验证
}

// EdgeEnumMap 多到多类型枚举转uint
var EdgeEnumMap = map[EdgeEnum]uint{
	"USER_TO_REVIEWER": 1, // 用户审核
	"USER_TO_ROLE":     2, // 用户角色
	"USER_TO_CLUSTER":  3, // 用户群集
}

// QueryTypeEnumMap 查询类型枚举转uint8
var QueryTypeEnumMap = map[QueryTypeEnum]uint8{
	"QUERY":   1,
	"ANALYZE": 2,
	"REWRITE": 3,
}
