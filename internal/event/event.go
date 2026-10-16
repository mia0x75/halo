package event

import "github.com/sirupsen/logrus"

// Audit and notification handlers are registered on separate emitters so the
// two delivery guarantees can differ: audit is synchronous and mandatory,
// notifications are asynchronous and droppable.
//
// Both use Void, which makes Emit run handlers synchronously in the caller
// goroutine without touching an (undrained) listener channel. Without it the
// buffered channels fill up and Emit deadlocks (E1).
var (
	auditEE  = New(1024, DefaultMatcher())
	notifyEE = New(1024, DefaultMatcher())
)

// fireConcurrency bounds concurrent notification goroutines so an event burst
// cannot spawn unbounded goroutines (E2).
const fireConcurrency = 100

var fireSem = make(chan struct{}, fireConcurrency)

// 注册事件处理函数
func init() {
	auditEE.Use("*", Void)
	notifyEE.Use("*", Void)

	auditEE.On(EventTicketCreated, TicketCreatedLogWriter)
	notifyEE.On(EventTicketCreated, TicketCreatedMailSender)
	notifyEE.On(EventTicketCreated, TicketCreatedStatisticUpdater)

	auditEE.On(EventTicketUpdated, TicketUpdatedLogWriter)
	notifyEE.On(EventTicketUpdated, TicketUpdatedMailSender)

	auditEE.On(EventTicketRemoved, TicketRemovedLogWriter)
	notifyEE.On(EventTicketRemoved, TicketRemovedMailSender)
	notifyEE.On(EventTicketRemoved, TicketRemovedStatisticsUpdater)

	auditEE.On(EventTicketExecuted, TicketExecutedLogWriter)
	notifyEE.On(EventTicketExecuted, TicketExecutedMailSender)

	auditEE.On(EventTicketFailed, TicketFailedLogWriter)
	notifyEE.On(EventTicketFailed, TicketFailedMailSender)

	auditEE.On(EventTicketScheduled, TicketScheduledLogWriter)
	notifyEE.On(EventTicketScheduled, TicketScheduledMailSender)

	auditEE.On(EventTicketStatusPatched, TicketStatusPatchedLogWriter)
	notifyEE.On(EventTicketStatusPatched, TicketStatusPatchedMailSender)

	auditEE.On(EventQueryCreated, QueryCreatedLogWriter)
	notifyEE.On(EventQueryCreated, QueryCreatedStatisticsUpdater)

	auditEE.On(EventQueryAnalyzed, QueryAnalyzedLogWriter)
	notifyEE.On(EventQueryAnalyzed, QueryAnalyzedStatisticsUpdater)

	auditEE.On(EventQueryRewrited, QueryRewritedLogWriter)
	notifyEE.On(EventQueryRewrited, QueryRewritedStatisticsUpdater)

	auditEE.On(EventUserRegistered, UserRegisteredLogWriter)
	notifyEE.On(EventUserRegistered, UserRegisteredMailSender)
	notifyEE.On(EventUserRegistered, UserRegisteredStatisticsUpdater)

	auditEE.On(EventUserSignedIn, UserSignedInLogWriter)
	notifyEE.On(EventUserSignedIn, UserSignedInStatisticsUpdater)

	auditEE.On(EventPasswordUpdated, PasswordUpdatedLogWriter)
	notifyEE.On(EventPasswordUpdated, PasswordUpdatedMailSender)

	auditEE.On(EventEmailUpdated, EmailUpdatedLogWriter)
	notifyEE.On(EventEmailUpdated, EmailUpdatedMailSender)

	auditEE.On(EventProfileUpdated, ProfileUpdatedLogWriter)
	notifyEE.On(EventProfileUpdated, ProfileUpdatedMailSender)

	auditEE.On(EventUserLogout, UserLogoutLogWriter)

	auditEE.On(EventUserCreated, UserCreatedLogWriter)
	notifyEE.On(EventUserCreated, UserCreatedMailSender)
	notifyEE.On(EventUserCreated, UserCreatedStatisticsUpdater)

	auditEE.On(EventUserUpdated, UserUpdatedLogWriter)

	auditEE.On(EventUserStatusPatched, UserStatusPatchedLogWriter)

	auditEE.On(EventRuleValuesPatched, RuleValuesPatchedLogWriter)

	auditEE.On(EventRuleBitwisePatched, RuleBitwisePatchedLogWriter)

	auditEE.On(EventOptionValuePatched, OptionValuePatchedLogWriter)

	auditEE.On(EventCommentCreated, CommentCreatedLogWriter)
	notifyEE.On(EventCommentCreated, CommentCreatedMailSender)
	notifyEE.On(EventCommentCreated, CommentCreatedStatisticsUpdater)

	auditEE.On(EventCronCancelled, CronCancelledLogWriter)
	notifyEE.On(EventCronCancelled, CronCancelledMailSender)

	auditEE.On(EventClusterStatusPatched, ClusterStatusPatchedLogWriter)

	auditEE.On(EventClusterRemoved, ClusterRemovedLogWriter)
	notifyEE.On(EventClusterRemoved, ClusterRemovedStatisticsUpdater)

	auditEE.On(EventClusterUpdated, ClusterUpdatedLogWriter)

	auditEE.On(EventClusterCreated, ClusterCreatedLogWriter)
	notifyEE.On(EventClusterCreated, ClusterCreatedStatisticsUpdater)

	auditEE.On(EventReviewerGranted, ReviewerGrantedLogWriter)

	auditEE.On(EventReviewerRevoked, ReviewerRevokedLogWriter)

	auditEE.On(EventClusterGranted, ClusterGrantedLogWriter)

	auditEE.On(EventClusterRevoked, ClusterRevokedLogWriter)

	auditEE.On(EventRoleGranted, RoleGrantedLogWriter)

	auditEE.On(EventRoleRevoked, RoleRevokedLogWriter)
}

// Fire 触发事件。审计同步执行（必达），通知/统计异步执行（有界、可丢）。
func Fire(topic string, args interface{}) {
	// Audit first: synchronous, so a write is recorded before the caller
	// returns and its failure is observable (E3).
	auditEE.Emit(topic, args)
	fireAsync(topic, args)
}

// FireSync 触发事件。审计与通知都同步执行，用于 CLI 执行等要求落库后再返回的路径。
func FireSync(topic string, args interface{}) {
	auditEE.Emit(topic, args)
	notifyEE.Emit(topic, args)
}

// fireAsync delivers notifications under a bounded limiter. Notifications are
// best-effort: when the limiter is saturated the event is dropped rather than
// blocking the caller or spawning unbounded goroutines (E2).
func fireAsync(topic string, args interface{}) {
	select {
	case fireSem <- struct{}{}:
	default:
		logrus.Warnf("[W] event fire limiter saturated, dropping topic=%s", topic)
		return
	}
	go func() {
		defer func() { <-fireSem }()
		notifyEE.Emit(topic, args)
	}()
}
