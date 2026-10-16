package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// newAudit builds an audit row with the shared defaults (fresh UUID, version 1,
// current time). userID is the actor of the operation.
func newAudit(userID uint, operation string) *domain.Log {
	return &domain.Log{
		UUID:      uuid.NewString(),
		UserID:    userID,
		Operation: operation,
		Version:   1,
		CreateAt:  uint(time.Now().Unix()),
	}
}

// ---------------------------------------------------------------------------
// Ticket events
// ---------------------------------------------------------------------------

// TicketCreated 工单创建成功事件
type TicketCreated struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// Topic routes the ticket-created event.
func (e *TicketCreated) Topic() string { return EventTicketCreated }

// AuditRecord returns nil because the audit line is written by
// service.TicketService.Create inside its transaction
// (Operation: "ticket.create"). See service.TicketService.
func (e *TicketCreated) AuditRecord() *domain.Log { return nil }

// TicketUpdated 工单更新成功事件
type TicketUpdated struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// Topic routes the ticket-updated event.
func (e *TicketUpdated) Topic() string { return EventTicketUpdated }

// AuditRecord returns nil because the audit line is written by
// service.TicketService.Update inside its transaction
// (Operation: "ticket.update"). See service.TicketService.
func (e *TicketUpdated) AuditRecord() *domain.Log { return nil }

// TicketRemoved 工单删除成功事件
type TicketRemoved struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// Topic routes the ticket-removed event.
func (e *TicketRemoved) Topic() string { return EventTicketRemoved }

// AuditRecord returns nil because the audit line is written by
// service.TicketService.Remove inside its transaction
// (Operation: "ticket.remove"). See service.TicketService.
func (e *TicketRemoved) AuditRecord() *domain.Log { return nil }

// TicketExecuted 工单执行成功事件。UserID is the actor resolved by the CLI
// (the seeded system user); the CLI path has no service-level audit.
type TicketExecuted struct {
	UserID  uint
	Ticket  models.Ticket
	Cluster models.Cluster
}

// Topic routes the ticket-executed event.
func (e *TicketExecuted) Topic() string { return EventTicketExecuted }

// AuditRecord records the execution; the CLI execution path does not go through
// any service that writes an audit row.
func (e *TicketExecuted) AuditRecord() *domain.Log {
	return newAudit(e.UserID, fmt.Sprintf("用户(uuid=%s)执行工单(uuid=%s)成功。\n", SystemUserUUID, e.Ticket.UUID))
}

// TicketFailed 工单执行失败事件
type TicketFailed struct {
	UserID  uint
	Ticket  models.Ticket
	Cluster models.Cluster
}

// Topic routes the ticket-failed event.
func (e *TicketFailed) Topic() string { return EventTicketFailed }

// AuditRecord records the failure; the CLI execution path has no service audit.
func (e *TicketFailed) AuditRecord() *domain.Log {
	return newAudit(e.UserID, fmt.Sprintf("用户(uuid=%s)执行工单(uuid=%s)失败。\n", SystemUserUUID, e.Ticket.UUID))
}

// TicketScheduled 工单预约成功事件
type TicketScheduled struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
	Cron    models.Cron
}

// Topic routes the ticket-scheduled event.
func (e *TicketScheduled) Topic() string { return EventTicketScheduled }

// AuditRecord records the schedule; the service only writes the ticket row, not
// an audit line.
func (e *TicketScheduled) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)在预约工单(uuid=%s)成功，预约时间(time=%v)。\n", e.User.UUID, e.Ticket.UUID, e.Cron.NextRun))
}

// TicketStatusPatched 工单状态更新成功事件
type TicketStatusPatched struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// Topic routes the ticket-status-patched event.
func (e *TicketStatusPatched) Topic() string { return EventTicketStatusPatched }

// AuditRecord returns nil because the audit line is written by
// service.TicketService.PatchStatus inside its transaction
// (Operation: "ticket.patch_status"). See service.TicketService.
func (e *TicketStatusPatched) AuditRecord() *domain.Log { return nil }

// ---------------------------------------------------------------------------
// Query events
// ---------------------------------------------------------------------------

// QueryCreated 查询创建成功事件
type QueryCreated struct {
	User  models.User
	Query models.Query
}

// Topic routes the query-created event.
func (e *QueryCreated) Topic() string { return EventQueryCreated }

// AuditRecord returns nil because the audit line is written by
// service.QueryService.Create inside its transaction
// (Operation: "query.create"). See service.QueryService.
func (e *QueryCreated) AuditRecord() *domain.Log { return nil }

// QueryAnalyzed 查询分析成功事件
type QueryAnalyzed struct {
	User  models.User
	Query models.Query
}

// Topic routes the query-analyzed event.
func (e *QueryAnalyzed) Topic() string { return EventQueryAnalyzed }

// AuditRecord returns nil because the resolver persists the analyzed query via
// service.QueryService.Create, which writes "query.create" in its transaction.
func (e *QueryAnalyzed) AuditRecord() *domain.Log { return nil }

// QueryRewrited 查询重写成功事件
type QueryRewrited struct {
	User  models.User
	Query models.Query
}

// Topic routes the query-rewritten event.
func (e *QueryRewrited) Topic() string { return EventQueryRewrited }

// AuditRecord returns nil because the resolver persists the rewritten query via
// service.QueryService.Create, which writes "query.create" in its transaction.
func (e *QueryRewrited) AuditRecord() *domain.Log { return nil }

// ---------------------------------------------------------------------------
// User events
// ---------------------------------------------------------------------------

// UserRegistered 用户注册成功事件
type UserRegistered struct {
	User models.User
}

// Topic routes the user-registered event.
func (e *UserRegistered) Topic() string { return EventUserRegistered }

// AuditRecord records the registration; no service writes this audit line.
func (e *UserRegistered) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)注册成功。\n", e.User.UUID))
}

// UserSignedIn 用户登录成功事件
type UserSignedIn struct {
	User models.User
	IP   string
}

// Topic routes the user-signed-in event.
func (e *UserSignedIn) Topic() string { return EventUserSignedIn }

// AuditRecord records the sign-in; login performs no service call.
func (e *UserSignedIn) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)通过远程主机(ip=%s)登录系统。", e.User.UUID, e.IP))
}

// PasswordUpdated 密码更新成功事件
type PasswordUpdated struct {
	User models.User
}

// Topic routes the password-updated event.
func (e *PasswordUpdated) Topic() string { return EventPasswordUpdated }

// AuditRecord records the password change; the user service writes only the
// generic "user.update" line, so the specific line is kept here.
func (e *PasswordUpdated) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)更新密码成功。\n", e.User.UUID))
}

// EmailUpdated 账号更新成功事件
type EmailUpdated struct {
	User  models.User
	Email string
}

// Topic routes the email-updated event.
func (e *EmailUpdated) Topic() string { return EventEmailUpdated }

// AuditRecord records the account change; the user service writes only the
// generic "user.update" line.
func (e *EmailUpdated) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)更新账号成功。\n", e.User.UUID))
}

// ProfileUpdated 用户信息更新成功事件
type ProfileUpdated struct {
	User models.User
}

// Topic routes the profile-updated event.
func (e *ProfileUpdated) Topic() string { return EventProfileUpdated }

// AuditRecord records the profile change; the user service writes only the
// generic "user.update" line.
func (e *ProfileUpdated) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)更新个人资料成功。\n", e.User.UUID))
}

// UserLogout 用户退出登录成功事件
type UserLogout struct {
	User models.User
}

// Topic routes the user-logout event.
func (e *UserLogout) Topic() string { return EventUserLogout }

// AuditRecord records the logout; logout performs no service call.
func (e *UserLogout) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)退出登录成功。\n", e.User.UUID))
}

// UserCreated 用户创建成功事件
type UserCreated struct {
	Manager models.User
	User    models.User
}

// Topic routes the user-created event.
func (e *UserCreated) Topic() string { return EventUserCreated }

// AuditRecord records the creation; the resolver writes the user via
// TxManager directly rather than UserService.Create, so there is no
// service-level audit line.
func (e *UserCreated) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)创建用户(uuid=%s)成功。\n", e.Manager.UUID, e.User.UUID))
}

// UserUpdated 用户更新成功事件
type UserUpdated struct {
	Manager models.User
	User    models.User
}

// Topic routes the user-updated event.
func (e *UserUpdated) Topic() string { return EventUserUpdated }

// AuditRecord returns nil because the audit line is written by
// service.UserService.Update inside its transaction
// (Operation: "user.update"). See service.UserService.
func (e *UserUpdated) AuditRecord() *domain.Log { return nil }

// UserStatusPatched 用户状态更新成功事件
type UserStatusPatched struct {
	Manager models.User
	User    models.User
}

// Topic routes the user-status-patched event.
func (e *UserStatusPatched) Topic() string { return EventUserStatusPatched }

// AuditRecord records the status change; the user service writes only the
// generic "user.update" line.
func (e *UserStatusPatched) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新用户(uuid=%s)状态成功。\n", e.Manager.UUID, e.User.UUID))
}

// ---------------------------------------------------------------------------
// Rule / Option events
// ---------------------------------------------------------------------------

// RuleParametersPatched 规则值更新成功事件
type RuleParametersPatched struct {
	Manager models.User
	Rule    models.Rule
}

// Topic routes the rule-parameters-patched event.
func (e *RuleParametersPatched) Topic() string { return EventRuleParametersPatched }

// AuditRecord records the rule value change; the rule service writes only the
// generic "rule.update" line.
func (e *RuleParametersPatched) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新规则(uuid=%s)值成功。\n", e.Manager.UUID, e.Rule.UUID))
}

// RuleBitwisePatched 规则状态位更新成功事件
type RuleBitwisePatched struct {
	Manager models.User
	Rule    models.Rule
}

// Topic routes the rule-bitwise-patched event.
func (e *RuleBitwisePatched) Topic() string { return EventRuleBitwisePatched }

// AuditRecord records the rule flag change; the rule service writes only the
// generic "rule.update" line.
func (e *RuleBitwisePatched) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新规则(uuid=%s)状态位成功。\n", e.Manager.UUID, e.Rule.UUID))
}

// OptionParameterPatched 系统选项更新成功事件
type OptionParameterPatched struct {
	Manager models.User
	Option  models.Option
}

// Topic routes the option-value-patched event.
func (e *OptionParameterPatched) Topic() string { return EventOptionParameterPatched }

// AuditRecord records the option change; the option service writes only the
// generic "option.update" line.
func (e *OptionParameterPatched) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新系统选项(uuid=%s)成功。\n", e.Manager.UUID, e.Option.UUID))
}

// ---------------------------------------------------------------------------
// Comment / Cron events
// ---------------------------------------------------------------------------

// CommentCreated 审核意见添加成功事件
type CommentCreated struct {
	User    models.User
	Ticket  models.Ticket
	Comment models.Comment
}

// Topic routes the comment-created event.
func (e *CommentCreated) Topic() string { return EventCommentCreated }

// AuditRecord returns nil because the audit line is written by
// service.CommentService.Create inside its transaction
// (Operation: "comment.create"). See service.CommentService.
func (e *CommentCreated) AuditRecord() *domain.Log { return nil }

// CronCancelled 工单预约取消成功事件
type CronCancelled struct {
	User   models.User
	Ticket models.Ticket
	Cron   models.Cron
}

// Topic routes the cron-cancelled event.
func (e *CronCancelled) Topic() string { return EventCronCancelled }

// AuditRecord records the cancellation; only read service calls precede the
// Fire, so no service audit line exists.
func (e *CronCancelled) AuditRecord() *domain.Log {
	return newAudit(e.User.UserID, fmt.Sprintf("用户(uuid=%s)取消工单(uuid=%s)执行成功。\n", e.User.UUID, e.Ticket.UUID))
}

// ---------------------------------------------------------------------------
// Cluster events
// ---------------------------------------------------------------------------

// ClusterStatusPatched 群集状态更新成功事件
type ClusterStatusPatched struct {
	Manager models.User
	Cluster models.Cluster
}

// Topic routes the cluster-status-patched event.
func (e *ClusterStatusPatched) Topic() string { return EventClusterStatusPatched }

// AuditRecord records the status change; the cluster service writes only the
// generic "cluster.update" line.
func (e *ClusterStatusPatched) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("用户(uuid=%s)修改群集(uuid=%s)状态成功。\n", e.Manager.UUID, e.Cluster.UUID))
}

// ClusterRemoved 群集移除成功事件
type ClusterRemoved struct {
	Manager models.User
	Cluster models.Cluster
}

// Topic routes the cluster-removed event.
func (e *ClusterRemoved) Topic() string { return EventClusterRemoved }

// AuditRecord returns nil because the audit line is written by
// service.ClusterService.Delete inside its transaction
// (Operation: "cluster.delete"). See service.ClusterService.
func (e *ClusterRemoved) AuditRecord() *domain.Log { return nil }

// ClusterUpdated 群集更新成功事件
type ClusterUpdated struct {
	Manager models.User
	Cluster models.Cluster
}

// Topic routes the cluster-updated event.
func (e *ClusterUpdated) Topic() string { return EventClusterUpdated }

// AuditRecord returns nil because the audit line is written by
// service.ClusterService.Update inside its transaction
// (Operation: "cluster.update"). See service.ClusterService.
func (e *ClusterUpdated) AuditRecord() *domain.Log { return nil }

// ClusterCreated 群集创建成功事件
type ClusterCreated struct {
	Manager models.User
	Cluster models.Cluster
}

// Topic routes the cluster-created event.
func (e *ClusterCreated) Topic() string { return EventClusterCreated }

// AuditRecord returns nil because the audit line is written by
// service.ClusterService.Create inside its transaction
// (Operation: "cluster.create"). See service.ClusterService.
func (e *ClusterCreated) AuditRecord() *domain.Log { return nil }

// ---------------------------------------------------------------------------
// Grant / revoke events
// ---------------------------------------------------------------------------

// ReviewerGranted 授权审核人成功事件
type ReviewerGranted struct {
	Manager models.User
	User    models.User
}

// Topic routes the reviewer-granted event.
func (e *ReviewerGranted) Topic() string { return EventReviewerGranted }

// AuditRecord records the grant; the edge service writes a technical
// "edge.replace" line, so the business-specific line is kept here.
func (e *ReviewerGranted) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)授权用户(uuid=%s)审核人成功。\n", e.Manager.UUID, e.User.UUID))
}

// ReviewerRevoked 收回审核人成功事件
type ReviewerRevoked struct {
	Manager models.User
	User    models.User
}

// Topic routes the reviewer-revoked event.
func (e *ReviewerRevoked) Topic() string { return EventReviewerRevoked }

// AuditRecord records the revoke; the edge service writes a technical
// "edge.delete" line.
func (e *ReviewerRevoked) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)收回用户(uuid=%s)审核人成功。\n", e.Manager.UUID, e.User.UUID))
}

// ClusterGranted 授权群集成功事件
type ClusterGranted struct {
	Manager models.User
	User    models.User
}

// Topic routes the cluster-granted event.
func (e *ClusterGranted) Topic() string { return EventClusterGranted }

// AuditRecord records the grant; the edge service writes a technical
// "edge.replace" line.
func (e *ClusterGranted) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)授权用户(uuid=%s)群集成功。\n", e.Manager.UUID, e.User.UUID))
}

// ClusterRevoked 收回群集成功事件
type ClusterRevoked struct {
	Manager models.User
	User    models.User
}

// Topic routes the cluster-revoked event.
func (e *ClusterRevoked) Topic() string { return EventClusterRevoked }

// AuditRecord records the revoke; the edge service writes a technical
// "edge.delete" line.
func (e *ClusterRevoked) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)收回用户(uuid=%s)群集成功。\n", e.Manager.UUID, e.User.UUID))
}

// RoleGranted 授权角色成功事件
type RoleGranted struct {
	Manager models.User
	User    models.User
}

// Topic routes the role-granted event.
func (e *RoleGranted) Topic() string { return EventRoleGranted }

// AuditRecord records the grant; the edge service writes a technical
// "edge.replace" line.
func (e *RoleGranted) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)授权用户(uuid=%s)角色成功。\n", e.Manager.UUID, e.User.UUID))
}

// RoleRevoked 收回角色成功事件
type RoleRevoked struct {
	Manager models.User
	User    models.User
}

// Topic routes the role-revoked event.
func (e *RoleRevoked) Topic() string { return EventRoleRevoked }

// AuditRecord records the revoke; the edge service writes a technical
// "edge.delete" line.
func (e *RoleRevoked) AuditRecord() *domain.Log {
	return newAudit(e.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)收回用户(uuid=%s)角色成功。\n", e.Manager.UUID, e.User.UUID))
}
