package event

import (
	"context"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/domain"
)

// Event is the interface every domain event implements. Handlers dispatch on
// the concrete type (type switch) so statistics and mail keep their per-event
// behaviour without an untyped payload.
type Event interface {
	// Topic returns the routing key (one of the Event* constants).
	Topic() string
	// AuditRecord returns the audit row to persist. A nil return means the
	// audit line is already written by the service layer inside its
	// transaction; see each implementation's comment for the source.
	AuditRecord() *domain.Log
}

// 事件列表
const (
	EventTicketCreated        = "OnTicketCreated"        // 工单创建成功
	EventTicketUpdated        = "OnTicketUpdated"        // 工单修改成功
	EventTicketRemoved        = "OnTicketRemoved"        // 工单删除成功
	EventTicketExecuted       = "OnTicketExecuted"       // 工单执行成功
	EventTicketFailed         = "OnTicketFailed"         // 工单执行失败
	EventTicketScheduled      = "OnTicketScheduled"      // 工单预约成功
	EventTicketStatusPatched  = "OnTicketStatusPatched"  // 工单状态修改成功
	EventQueryCreated         = "OnQueryCreated"         // 创建执行查询
	EventQueryAnalyzed        = "OnQueryAnalyzed"        // 查询分析成功
	EventQueryRewrited        = "OnQueryRewrited"        // 查询重写成功
	EventUserRegistered       = "OnUserReigstered"       // 用户注册成功
	EventUserSignedIn         = "OnUserSignedIn"         // 用户登录成功
	EventPasswordUpdated      = "OnPasswordUpdated"      // 用户修改密码成功
	EventEmailUpdated         = "OnEmailUpdated"         // 用户修改账号成功
	EventProfileUpdated       = "OnProfileUpdated"       // 用户修改个人资料成功
	EventUserLogout           = "OnUserLogout"           // 用户退出登录成功
	EventUserCreated          = "OnUserCreated"          // 用户创建成功
	EventUserUpdated          = "OnUserUpdated"          // 用户更新成功
	EventUserStatusPatched    = "OnUserStatusPatched"    // 用户状态修改成功
	EventRuleValuesPatched    = "OnRuleValuesPatched"    // 规则取值修改成功
	EventRuleBitwisePatched   = "OnRuleBitwisePatched"   // 规则执行标志修改成功
	EventOptionValuePatched   = "OnOptionValuePatched"   // 系统选项修改成功
	EventCommentCreated       = "OnCommentCreated"       // 添加审核意见成功
	EventCronCancelled        = "OnCronCancelled"        // 计划任务取消成功
	EventClusterStatusPatched = "OnClusterStatusPatched" // 群集状态修改成功
	EventClusterRemoved       = "OnClusterRemoved"       // 群集移除成功
	EventClusterUpdated       = "OnClusterUpdated"       // 群集修改成功
	EventClusterCreated       = "OnClusterCreated"       // 群集创建成功
	EventReviewerGranted      = "OnReviewerGranted"      // 授权审核人成功
	EventReviewerRevoked      = "OnReviewerRevoked"      // 收回审核人成功
	EventClusterGranted       = "OnClusterGranted"       // 授权群集成功
	EventClusterRevoked       = "OnClusterRevoked"       // 收回群集成功
	EventRoleGranted          = "OnRoleGranted"          // 授权角色成功
	EventRoleRevoked          = "OnRoleRevoked"          // 收回角色成功
)

// SystemUserUUID is the seeded system account used as the actor of automated
// actions such as ticket execution.
const SystemUserUUID = "00000000-0000-0000-0000-000000000000"

// defaultBus is the process-wide bus set by Register. The package-level Fire
// helpers keep the ~39 call sites unchanged; they are only valid after the app
// container has been built.
var defaultBus atomic.Pointer[Bus]

// setDefaultBus records the bus used by the package-level helpers.
func setDefaultBus(b *Bus) { defaultBus.Store(b) }

// Fire 触发事件：审计同步（必达），通知/统计有界异步（可丢）。
func Fire(topic string, ev Event) {
	if b := defaultBus.Load(); b != nil {
		if err := b.Fire(context.Background(), topic, ev); err != nil {
			log.Errorf("[E] event fire failed, topic=%s: %s", topic, err.Error())
		}
	}
}

// FireSync 触发事件：审计与通知都同步执行。
func FireSync(topic string, ev Event) {
	if b := defaultBus.Load(); b != nil {
		if err := b.FireSync(context.Background(), topic, ev); err != nil {
			log.Errorf("[E] event fire_sync failed, topic=%s: %s", topic, err.Error())
		}
	}
}
