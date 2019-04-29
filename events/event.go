package events

var ee = New(1024, DefaultMatcher())

// 注册事件处理函数
func init() {
	ee.On(EventTicketCreated, TicketCreatedLogWriter)
	ee.On(EventTicketCreated, TicketCreatedMailSender)
	ee.On(EventTicketCreated, TicketCreatedStatisticUpdater)

	ee.On(EventTicketUpdated, TicketUpdatedLogWriter)
	ee.On(EventTicketUpdated, TicketUpdatedMailSender)

	ee.On(EventTicketRemoved, TicketRemovedLogWriter)
	ee.On(EventTicketRemoved, TicketRemovedMailSender)
	ee.On(EventTicketRemoved, TicketRemovedStatisticsUpdater)

	ee.On(EventTicketExecuted, TicketExecutedLogWriter)
	ee.On(EventTicketExecuted, TicketExecutedMailSender)

	ee.On(EventTicketFailed, TicketFailedLogWriter)
	ee.On(EventTicketFailed, TicketFailedMailSender)

	ee.On(EventTicketScheduled, TicketScheduledLogWriter)
	ee.On(EventTicketScheduled, TicketScheduledMailSender)

	ee.On(EventTicketStatusPatched, TicketStatusPatchedLogWriter)
	ee.On(EventTicketStatusPatched, TicketStatusPatchedMailSender)

	ee.On(EventQueryCreated, QueryCreatedLogWriter)
	ee.On(EventQueryCreated, QueryCreatedStatisticsUpdater)

	ee.On(EventQueryAnalyzed, QueryAnalyzedLogWriter)
	ee.On(EventQueryAnalyzed, QueryAnalyzedStatisticsUpdater)

	ee.On(EventQueryRewrited, QueryRewritedLogWriter)
	ee.On(EventQueryRewrited, QueryRewritedStatisticsUpdater)

	ee.On(EventUserRegistered, UserRegisteredLogWriter)
	ee.On(EventUserRegistered, UserRegisteredMailSender)
	ee.On(EventUserRegistered, UserRegisteredStatisticsUpdater)

	ee.On(EventUserSignedIn, UserSignedInLogWriter)
	ee.On(EventUserSignedIn, UserSignedInStatisticsUpdater)

	ee.On(EventPasswordUpdated, PasswordUpdatedLogWriter)
	ee.On(EventPasswordUpdated, PasswordUpdatedMailSender)

	ee.On(EventEmailUpdated, EmailUpdatedLogWriter)
	ee.On(EventEmailUpdated, EmailUpdatedMailSender)

	ee.On(EventProfileUpdated, ProfileUpdatedLogWriter)
	ee.On(EventProfileUpdated, ProfileUpdatedMailSender)

	ee.On(EventUserLogout, UserLogoutLogWriter)

	ee.On(EventUserCreated, UserCreatedLogWriter)
	ee.On(EventUserCreated, UserCreatedMailSender)
	ee.On(EventUserCreated, UserCreatedStatisticsUpdater)

	ee.On(EventUserUpdated, UserUpdatedLogWriter)

	ee.On(EventUserStatusPatched, UserStatusPatchedLogWriter)

	ee.On(EventRuleValuesPatched, RuleValuesPatchedLogWriter)

	ee.On(EventRuleBitwisePatched, RuleBitwisePatchedLogWriter)

	ee.On(EventOptionValuePatched, OptionValuePatchedLogWriter)

	ee.On(EventCommentCreated, CommentCreatedLogWriter)
	ee.On(EventCommentCreated, CommentCreatedMailSender)
	ee.On(EventCommentCreated, CommentCreatedStatisticsUpdater)

	ee.On(EventCronCancelled, CronCancelledLogWriter)
	ee.On(EventCronCancelled, CronCancelledMailSender)

	ee.On(EventClusterStatusPatched, ClusterStatusPatchedLogWriter)

	ee.On(EventClusterRemoved, ClusterRemovedLogWriter)
	ee.On(EventClusterRemoved, ClusterRemovedStatisticsUpdater)

	ee.On(EventClusterUpdated, ClusterUpdatedLogWriter)

	ee.On(EventClusterCreated, ClusterCreatedLogWriter)
	ee.On(EventClusterCreated, ClusterCreatedStatisticsUpdater)

	ee.On(EventReviewerGranted, ReviewerGrantedLogWriter)

	ee.On(EventReviewerRevoked, ReviewerRevokedLogWriter)

	ee.On(EventClusterGranted, ClusterGrantedLogWriter)

	ee.On(EventClusterRevoked, ClusterRevokedLogWriter)

	ee.On(EventRoleGranted, RoleGrantedLogWriter)

	ee.On(EventRoleRevoked, RoleRevokedLogWriter)
}

// Fire 触发事件
func Fire(topic string, args interface{}) {
	go ee.Emit(topic, args)
}

// FireSync 触发事件
func FireSync(topic string, args interface{}) {
	ee.Emit(topic, args)
}
