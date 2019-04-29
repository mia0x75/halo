package events

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
)

// 事件列表
const (
	EventTicketCreated        = "OnTicketCreated"        // 工单创建成功 - PASS
	EventTicketUpdated        = "OnTicketUpdated"        // 工单修改成功 - PASS
	EventTicketRemoved        = "OnTicketRemoved"        // 工单删除成功 - PASS
	EventTicketExecuted       = "OnTicketExecuted"       // 工单执行成功
	EventTicketFailed         = "OnTicketFailed"         // 工单执行失败
	EventTicketScheduled      = "OnTicketScheduled"      // 工单预约成功
	EventTicketStatusPatched  = "OnTicketStatusPatched"  // 工单状态修改成功
	EventQueryCreated         = "OnQueryCreated"         // 创建执行查询 - PASS
	EventQueryAnalyzed        = "OnQueryAnalyzed"        // 查询分析成功 - PASS
	EventQueryRewrited        = "OnQueryRewrited"        // 查询重写成功 - PASS
	EventUserRegistered       = "OnUserReigstered"       // 用户注册成功 - PASS
	EventUserSignedIn         = "OnUserSignedIn"         // 用户登录成功 - PASS
	EventPasswordUpdated      = "OnPasswordUpdated"      // 用户修改密码成功 - PASS
	EventEmailUpdated         = "OnEmailUpdated"         // 用户修改账号成功 - PASS
	EventProfileUpdated       = "OnProfileUpdated"       // 用户修改个人资料成功 - PASS
	EventUserLogout           = "OnUserLogout"           // 用户退出登录 - PASS
	EventUserCreated          = "OnUserCreated"          // 用户创建成功 - PASS
	EventUserUpdated          = "OnUserUpdated"          // 用户更新成功 - PASS
	EventUserStatusPatched    = "OnUserStatusPatched"    // 用户状态修改成功 - PASS
	EventRuleValuesPatched    = "OnRuleValuesPatched"    // 规则取值修改成功 - PASS
	EventRuleBitwisePatched   = "OnRuleBitwisePatched"   // 规则执行标志修改成功 - PASS
	EventOptionValuePatched   = "OnOptionValuePatched"   // 系统选项修改成功 - PASS
	EventCommentCreated       = "OnCommentCreated"       // 添加审核意见成功 - PASS
	EventCronCancelled        = "OnCronCancelled"        // 计划任务取消成功
	EventClusterStatusPatched = "OnClusterStatusPatched" // 群集状态修改成功 - PASS
	EventClusterRemoved       = "OnClusterRemoved"       // 群集移除成功 - PASS
	EventClusterUpdated       = "OnClusterUpdated"       // 群集修改成功 - PASS
	EventClusterCreated       = "OnClusterCreated"       // 群集创建成功 - PASS
	EventReviewerGranted      = "OnReviewerGranted"      // 授权审核人成功 - PASS
	EventReviewerRevoked      = "OnReviewerRevoked"      // 收回审核人成功 - PASS
	EventClusterGranted       = "OnClusterGranted"       // 授权群集成功 - PASS
	EventClusterRevoked       = "OnClusterRevoked"       // 收回群集成功 - PASS
	EventRoleGranted          = "OnRoleGranted"          // 授权角色成功 - PASS
	EventRoleRevoked          = "OnRoleRevoked"          // 收回角色成功 - PASS
)

// TicketCreatedArgs 工单创建成功事件参数
type TicketCreatedArgs struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// TicketCreatedLogWriter 创建工单的日志记录
func TicketCreatedLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketCreatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)创建了一个新工单(uuid=%s)。\n", args.User.UUID, args.Ticket.UUID))
	}
}

// TicketCreatedMailSender 创建工单的邮件通知
func TicketCreatedMailSender(e *Event) {
	if args, ok := e.Args.(*TicketCreatedArgs); ok {
		subject, body := renderer(g.TplTicketCreated, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// TicketCreatedStatisticUpdater 创建工单的统计更新
func TicketCreatedStatisticUpdater(e *Event) {
	if args, ok := e.Args.(*TicketCreatedArgs); ok {
		if err := UpdateStatistics("overall", "total-tickets", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		today := time.Now().Format("2006-01-02")
		if err := UpdateStatistics("tickets-daily", today, 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "total-tickets", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// TicketUpdatedArgs 工单更新成功事件参数
type TicketUpdatedArgs struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// TicketUpdatedLogWriter 工单更新日志记录
func TicketUpdatedLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketUpdatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)更新了工单(uuid=%s)。\n", args.User.UUID, args.Ticket.UUID))
	}
}

// TicketUpdatedMailSender 更新工单的邮件通知
func TicketUpdatedMailSender(e *Event) {
	if args, ok := e.Args.(*TicketUpdatedArgs); ok {
		subject, body := renderer(g.TplTicketUpdated, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// TicketRemovedArgs 工单删除成功事件参数
type TicketRemovedArgs struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// TicketRemovedLogWriter 删除工单的日志记录
func TicketRemovedLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketRemovedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)删除了工单(uuid=%s)。\n", args.User.UUID, args.Ticket.UUID))
	}
}

// TicketRemovedMailSender 删除工单的邮件通知
func TicketRemovedMailSender(e *Event) {
	if args, ok := e.Args.(*TicketRemovedArgs); ok {
		// 这个用于邮件模版的时间显示
		args.Ticket.UpdateAt = uint(time.Now().UTC().Unix())
		subject, body := renderer(g.TplTicketRemoved, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// TicketRemovedStatisticsUpdater 删除工单的统计更新
func TicketRemovedStatisticsUpdater(e *Event) {
	if args, ok := e.Args.(*TicketRemovedArgs); ok {
		if err := UpdateStatistics("overall", "total-tickets", -1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		today := time.Now().Format("2006-01-02")
		if err := UpdateStatistics("tickets-daily", today, -1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "total-tickets", -1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// TicketExecutedArgs 工单执行成功事件参数
type TicketExecutedArgs struct {
	Ticket  models.Ticket
	Cluster models.Cluster
}

// TicketExecutedLogWriter 工单执行成功的日志记录
func TicketExecutedLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketExecutedArgs); ok {
		agent := &models.User{}
		if _, err := g.Engine.Where("`uuid` = ?", "00000000-0000-0000-0000-000000000000").Get(agent); err != nil {
			return
		}
		LogWriter(agent.UserID, fmt.Sprintf("用户(uuid=%s)执行工单(uuid=%s)成功。\n", agent.UUID, args.Ticket.UUID))
	}
}

// TicketExecutedMailSender 工单执行成功的邮件通知
func TicketExecutedMailSender(e *Event) {
	if args, ok := e.Args.(*TicketExecutedArgs); ok {
		users := []models.User{}
		if err := g.Engine.In("user_id", args.Ticket.UserID, args.Ticket.ReviewerID).Find(&users); err != nil {
			return
		}
		subject, body := renderer(g.TplTicketExecuted, args)
		for _, user := range users {
			MailSender(MailSendArgs{
				To:      mail.Address{Name: user.Name, Address: user.Email},
				Subject: subject,
				Body:    body,
			})
		}
	}
}

// TicketFailedArgs 工单执行失败事件参数
type TicketFailedArgs struct {
	Ticket  models.Ticket
	Cluster models.Cluster
}

// TicketFailedLogWriter 工单执行失败的日志记录
func TicketFailedLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketFailedArgs); ok {
		agent := &models.User{}
		if _, err := g.Engine.Where("`uuid` = ?", "00000000-0000-0000-0000-000000000000").Get(agent); err != nil {
			return
		}
		LogWriter(agent.UserID, fmt.Sprintf("用户(uuid=%s)执行工单(uuid=%s)失败。\n", agent.UUID, args.Ticket.UUID))
	}
}

// TicketFailedMailSender 工单执行失败的邮件通知
func TicketFailedMailSender(e *Event) {
	if args, ok := e.Args.(*TicketFailedArgs); ok {
		users := []models.User{}
		if err := g.Engine.In("user_id", args.Ticket.UserID, args.Ticket.ReviewerID).Find(&users); err != nil {
			return
		}
		subject, body := renderer(g.TplTicketFailed, args)
		for _, user := range users {
			MailSender(MailSendArgs{
				To:      mail.Address{Name: user.Name, Address: user.Email},
				Subject: subject,
				Body:    body,
			})
		}
	}
}

// TicketScheduledArgs 工单预约成功事件参数
type TicketScheduledArgs struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
	Cron    models.Cron
}

// TicketScheduledLogWriter 工单预约成功的日志记录
func TicketScheduledLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketScheduledArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)在预约工单(uuid=%s)成功，预约时间(time=%v)。\n", args.User.UUID, args.Ticket.UUID, args.Cron.NextRun))
	}
}

// TicketScheduledMailSender 工单预约成功的邮件通知
func TicketScheduledMailSender(e *Event) {
	if args, ok := e.Args.(*TicketScheduledArgs); ok {
		users := caches.UsersMap.Filter(func(elem *models.User) bool {
			if elem.UserID == args.Ticket.UserID || elem.UserID == args.Ticket.ReviewerID {
				return true
			}
			return false
		})
		subject, body := renderer(g.TplTicketScheduled, args)
		for _, user := range users {
			MailSender(MailSendArgs{
				To:      mail.Address{Name: user.Name, Address: user.Email},
				Subject: subject,
				Body:    body,
			})
		}
	}
}

// TicketStatusPatchedArgs 工单状态更新成功事件参数
type TicketStatusPatchedArgs struct {
	User    models.User
	Ticket  models.Ticket
	Cluster models.Cluster
}

// TicketStatusPatchedLogWriter 工单状态更新成功日志记录
func TicketStatusPatchedLogWriter(e *Event) {
	if args, ok := e.Args.(*TicketStatusPatchedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)修改工单(uuid=%s)状态成功。\n", args.User.UUID, args.Ticket.UUID))
	}
}

// TicketStatusPatchedMailSender 工单状态更新成功邮件通知
func TicketStatusPatchedMailSender(e *Event) {
	if args, ok := e.Args.(*TicketStatusPatchedArgs); ok {
		users := caches.UsersMap.Filter(func(elem *models.User) bool {
			if elem.UserID == args.Ticket.UserID || elem.UserID == args.Ticket.ReviewerID {
				return true
			}
			return false
		})
		var subject, body string
		switch args.Ticket.Status {
		case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumClosed]:
			users = []*models.User{&args.User}
			subject, body = renderer(g.TplTicketClosed, args)
		case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumMrvFailure]:
			subject, body = renderer(g.TplTicketMrvFailure, args)
		case gqlapi.TicketStatusEnumMap[gqlapi.TicketStatusEnumLgtm]:
			subject, body = renderer(g.TplTicketLgtm, args)
		default:
			return
		}
		for _, user := range users {
			MailSender(MailSendArgs{
				To:      mail.Address{Name: user.Name, Address: user.Email},
				Subject: subject,
				Body:    body,
			})
		}
	}
}

// QueryCreatedArgs 查询创建成功事件参数
type QueryCreatedArgs struct {
	User  models.User
	Query models.Query
}

// QueryCreatedLogWriter 创建查询的日志记录
func QueryCreatedLogWriter(e *Event) {
	if args, ok := e.Args.(*QueryCreatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)执行数据查询(uuid=%s)成功。\n", args.User.UUID, args.Query.UUID))
	}
}

// QueryCreatedStatisticsUpdater 创建查询的统计更新
func QueryCreatedStatisticsUpdater(e *Event) {
	if args, ok := e.Args.(*QueryCreatedArgs); ok {
		if err := UpdateStatistics("overall", "total-queries", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		today := time.Now().Format("2006-01-02")
		if err := UpdateStatistics("queries-daily", today, 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "total-queries", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// QueryAnalyzedArgs 查询分析成功事件参数
type QueryAnalyzedArgs struct {
	User  models.User
	Query models.Query
}

// QueryAnalyzedLogWriter 查询分析的日志记录
func QueryAnalyzedLogWriter(e *Event) {
	if args, ok := e.Args.(*QueryAnalyzedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)执行查询分析(uuid=%s)成功。\n", args.User.UUID, args.Query.UUID))
	}
}

// QueryAnalyzedStatisticsUpdater 查询分析的统计更新
func QueryAnalyzedStatisticsUpdater(e *Event) {
	if args, ok := e.Args.(*QueryAnalyzedArgs); ok {
		if err := UpdateStatistics("overall", "total-queries", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		today := time.Now().Format("2006-01-02")
		if err := UpdateStatistics("queries-daily", today, 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "total-analyses", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// QueryRewritedArgs 查询分析成功事件参数
type QueryRewritedArgs struct {
	User  models.User
	Query models.Query
}

// QueryRewritedLogWriter 查询重写的日志记录
func QueryRewritedLogWriter(e *Event) {
	if args, ok := e.Args.(*QueryRewritedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)执行查询重写(uuid=%s)成功。\n", args.User.UUID, args.Query.UUID))
	}
}

// QueryRewritedStatisticsUpdater 查询重写的统计更新
func QueryRewritedStatisticsUpdater(e *Event) {
	if args, ok := e.Args.(*QueryRewritedArgs); ok {
		if err := UpdateStatistics("overall", "total-queries", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		today := time.Now().Format("2006-01-02")
		if err := UpdateStatistics("queries-daily", today, 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "total-rewrites", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// UserRegisteredArgs 用户注册成功事件参数
type UserRegisteredArgs struct {
	User models.User
}

// UserRegisteredLogWriter 用户注册的日志记录
func UserRegisteredLogWriter(e *Event) {
	if args, ok := e.Args.(*UserRegisteredArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)注册成功。\n", args.User.UUID))
	}
}

// UserRegisteredMailSender 用户注册的邮件通知（激活）
func UserRegisteredMailSender(e *Event) {
	if args, ok := e.Args.(*UserRegisteredArgs); ok {
		subject, body := renderer(g.TplUserRegistered, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// UserRegisteredStatisticsUpdater 用户注册的统计更新
func UserRegisteredStatisticsUpdater(e *Event) {
	if _, ok := e.Args.(*UserRegisteredArgs); ok {
		if err := UpdateStatistics("overall", "total-users", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// UserSignedInArgs 用户登录成功事件参数
type UserSignedInArgs struct {
	User models.User
	IP   string
}

// UserSignedInLogWriter 用户登录的日志记录
func UserSignedInLogWriter(e *Event) {
	if args, ok := e.Args.(*UserSignedInArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)通过远程主机(ip=%s)登录系统。", args.User.UUID, args.IP))
	}
}

// UserSignedInStatisticsUpdater 用户登录的统计更新
func UserSignedInStatisticsUpdater(e *Event) {
	if args, ok := e.Args.(*UserSignedInArgs); ok {
		if err := UpdateStatistics(args.User.UUID, "total-logins", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "latest-login", float64(time.Now().Unix())); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// PasswordUpdatedArgs 密码更新成功事件参数
type PasswordUpdatedArgs struct {
	User models.User
}

// PasswordUpdatedLogWriter 密码更新日志记录
func PasswordUpdatedLogWriter(e *Event) {
	if args, ok := e.Args.(*PasswordUpdatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)更新密码成功。\n", args.User.UUID))
	}
}

// PasswordUpdatedMailSender 密码更新邮件通知
func PasswordUpdatedMailSender(e *Event) {
	if args, ok := e.Args.(*PasswordUpdatedArgs); ok {
		subject, body := renderer(g.TplPasswordUpdated, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// EmailUpdatedArgs 账号更新成功事件参数
type EmailUpdatedArgs struct {
	User  models.User
	Email string
}

// EmailUpdatedLogWriter 账号更新日志记录
func EmailUpdatedLogWriter(e *Event) {
	if args, ok := e.Args.(*EmailUpdatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)更新账号成功。\n", args.User.UUID))
	}
}

// EmailUpdatedMailSender 账号更新邮件通知
func EmailUpdatedMailSender(e *Event) {
	if args, ok := e.Args.(*EmailUpdatedArgs); ok {
		subject, body := renderer(g.TplEmailUpdated, args)
		// TODO: 把新Email和当前时间加密后，放在验证邮件中
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// ProfileUpdatedArgs 用户信息更新成功事件参数
type ProfileUpdatedArgs struct {
	User models.User
}

// ProfileUpdatedLogWriter 用户信息更新日志记录
func ProfileUpdatedLogWriter(e *Event) {
	if args, ok := e.Args.(*ProfileUpdatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)更新个人资料成功。\n", args.User.UUID))
	}
}

// ProfileUpdatedMailSender 用户信息更新邮件通知
func ProfileUpdatedMailSender(e *Event) {
	if args, ok := e.Args.(*ProfileUpdatedArgs); ok {
		subject, body := renderer(g.TplProfileUpdated, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// UserLogoutArgs 用户退出登录成功事件参数
type UserLogoutArgs struct {
	User models.User
}

// UserLogoutLogWriter 用户退出登录日志记录
func UserLogoutLogWriter(e *Event) {
	if args, ok := e.Args.(*UserLogoutArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)退出登录成功。\n", args.User.UUID))
	}
}

// UserCreatedArgs 用户创建成功事件参数
type UserCreatedArgs struct {
	Manager models.User
	User    models.User
}

// UserCreatedLogWriter 用户创建成功日志记录
func UserCreatedLogWriter(e *Event) {
	if args, ok := e.Args.(*UserCreatedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)创建用户(uuid=%s)成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// UserCreatedMailSender 用户创建成功邮件通知
func UserCreatedMailSender(e *Event) {
	if args, ok := e.Args.(*UserCreatedArgs); ok {
		subject, body := renderer(g.TplUserCreated, args)
		MailSender(MailSendArgs{
			To:      mail.Address{Name: args.User.Name, Address: args.User.Email},
			Subject: subject,
			Body:    body,
		})
	}
}

// UserCreatedStatisticsUpdater 用户创建成功统计更新
func UserCreatedStatisticsUpdater(e *Event) {
	if _, ok := e.Args.(*UserCreatedArgs); ok {
		if err := UpdateStatistics("overall", "total-users", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// UserUpdatedArgs 用户更新成功事件参数
type UserUpdatedArgs struct {
	Manager models.User
	User    models.User
}

// UserUpdatedLogWriter 用户更新成功日志记录
func UserUpdatedLogWriter(e *Event) {
	if args, ok := e.Args.(*UserUpdatedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新用户(uuid=%s)成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// UserStatusPatchedArgs 用户状态更新成功事件参数
type UserStatusPatchedArgs struct {
	Manager models.User
	User    models.User
}

// UserStatusPatchedLogWriter 用户状态更新成功日志记录
func UserStatusPatchedLogWriter(e *Event) {
	if args, ok := e.Args.(*UserStatusPatchedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新用户(uuid=%s)状态成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// RuleValuesPatchedArgs 规则值更新成功事件参数
type RuleValuesPatchedArgs struct {
	Manager models.User
	Rule    models.Rule
}

// RuleValuesPatchedLogWriter 规则值更新成功日志记录
func RuleValuesPatchedLogWriter(e *Event) {
	if args, ok := e.Args.(*RuleValuesPatchedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新规则(uuid=%s)值成功。\n", args.Manager.UUID, args.Rule.UUID))
	}
}

// RuleBitwisePatchedArgs 规则状态位更新成功事件参数
type RuleBitwisePatchedArgs struct {
	Manager models.User
	Rule    models.Rule
}

// RuleBitwisePatchedLogWriter 规则状态位更新成功日志记录
func RuleBitwisePatchedLogWriter(e *Event) {
	if args, ok := e.Args.(*RuleBitwisePatchedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新规则(uuid=%s)状态位成功。\n", args.Manager.UUID, args.Rule.UUID))
	}
}

// OptionValuePatchedArgs 系统选项更新事件参数
type OptionValuePatchedArgs struct {
	Manager models.User
	Option  models.Option
}

// OptionValuePatchedLogWriter 系统选项更新日志记录
func OptionValuePatchedLogWriter(e *Event) {
	if args, ok := e.Args.(*OptionValuePatchedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)更新系统选项(uuid=%s)成功。\n", args.Manager.UUID, args.Option.UUID))
	}
}

// CommentCreatedArgs 审核意见添加成功事件参数
type CommentCreatedArgs struct {
	User    models.User
	Ticket  models.Ticket
	Comment models.Comment
}

// CommentCreatedLogWriter 审核意见添加成功日记记录
func CommentCreatedLogWriter(e *Event) {
	if args, ok := e.Args.(*CommentCreatedArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)添加工单(uuid=%s)审核意见成功。\n", args.User.UUID, args.Ticket.UUID))
	}
}

// CommentCreatedStatisticsUpdater 审核意见添加成功统计更新
func CommentCreatedStatisticsUpdater(e *Event) {
	if args, ok := e.Args.(*CommentCreatedArgs); ok {
		if err := UpdateStatistics("overall", "total-comments", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
		if err := UpdateStatistics(args.User.UUID, "total-comments", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// CommentCreatedMailSender 审核意见添加成功邮件通知
func CommentCreatedMailSender(e *Event) {
	if args, ok := e.Args.(*CommentCreatedArgs); ok {
		users := caches.UsersMap.Filter(func(elem *models.User) bool {
			if elem.UserID == args.Ticket.UserID || elem.UserID == args.Ticket.ReviewerID {
				return true
			}
			return false
		})
		subject, body := renderer(g.TplCommentCreated, args)
		for _, user := range users {
			MailSender(MailSendArgs{
				To:      mail.Address{Name: user.Name, Address: user.Email},
				Subject: subject,
				Body:    body,
			})
		}
	}
}

// CronCancelledArgs 工单预约取消事件参数
type CronCancelledArgs struct {
	User   models.User
	Ticket models.Ticket
	Cron   models.Cron
}

// CronCancelledLogWriter 工单预约取消日志记录
func CronCancelledLogWriter(e *Event) {
	if args, ok := e.Args.(*CronCancelledArgs); ok {
		LogWriter(args.User.UserID, fmt.Sprintf("用户(uuid=%s)取消工单(uuid=%s)执行成功。\n", args.User.UUID, args.Ticket.UUID))
	}
}

// CronCancelledMailSender 工单预约取消邮件通知
func CronCancelledMailSender(e *Event) {
	if args, ok := e.Args.(*CronCancelledArgs); ok {
		users := caches.UsersMap.Filter(func(elem *models.User) bool {
			if elem.UserID == args.Ticket.UserID || elem.UserID == args.Ticket.ReviewerID {
				return true
			}
			return false
		})
		subject, body := renderer(g.TplCronCancelled, args)
		for _, user := range users {
			MailSender(MailSendArgs{
				To:      mail.Address{Name: user.Name, Address: user.Email},
				Subject: subject,
				Body:    body,
			})
		}
	}
}

// ClusterStatusPatchedArgs 群集状态更新事件参数
type ClusterStatusPatchedArgs struct {
	Manager models.User
	Cluster models.Cluster
}

// ClusterStatusPatchedLogWriter 群集状态更新日志记录
func ClusterStatusPatchedLogWriter(e *Event) {
	if args, ok := e.Args.(*ClusterStatusPatchedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("用户(uuid=%s)修改群集(uuid=%s)状态成功。\n", args.Manager.UUID, args.Cluster.UUID))
	}
}

// ClusterRemovedArgs 群集移除事件参数
type ClusterRemovedArgs struct {
	Manager models.User
	Cluster models.Cluster
}

// ClusterRemovedLogWriter 群集移除日志记录
func ClusterRemovedLogWriter(e *Event) {
	if args, ok := e.Args.(*ClusterRemovedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("用户(uuid=%s)删除群集(uuid=%s)成功。\n", args.Manager.UUID, args.Cluster.UUID))
	}
}

// ClusterRemovedStatisticsUpdater 群集移除统计更新
func ClusterRemovedStatisticsUpdater(e *Event) {
	if _, ok := e.Args.(*ClusterRemovedArgs); ok {
		if err := UpdateStatistics("overall", "total-clusters", -1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// ClusterUpdatedArgs 群集更新事件参数
type ClusterUpdatedArgs struct {
	Manager models.User
	Cluster models.Cluster
}

// ClusterUpdatedLogWriter 群集更新日志记录
func ClusterUpdatedLogWriter(e *Event) {
	if args, ok := e.Args.(*ClusterUpdatedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("用户(uuid=%s)更新群集(uuid=%s)成功。\n", args.Manager.UUID, args.Cluster.UUID))
	}
}

// ClusterCreatedArgs 群集创建事件参数
type ClusterCreatedArgs struct {
	Manager models.User
	Cluster models.Cluster
}

// ClusterCreatedLogWriter 群集创建日志记录
func ClusterCreatedLogWriter(e *Event) {
	if args, ok := e.Args.(*ClusterCreatedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("用户(uuid=%s)创建群集(uuid=%s)成功。\n", args.Manager.UUID, args.Cluster.UUID))
	}
}

// ClusterCreatedStatisticsUpdater 群集创建统计更新
func ClusterCreatedStatisticsUpdater(e *Event) {
	if _, ok := e.Args.(*ClusterCreatedArgs); ok {
		if err := UpdateStatistics("overall", "total-clusters", 1); err != nil {
			log.Warnf("[W] An unexpected error occurred: %s", err.Error())
		}
	}
}

// ReviewerGrantedArgs 授权审核人事件参数
type ReviewerGrantedArgs struct {
	Manager models.User
	User    models.User
}

// ReviewerGrantedLogWriter 授权审核人日志记录
func ReviewerGrantedLogWriter(e *Event) {
	if args, ok := e.Args.(*ReviewerGrantedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)授权用户(uuid=%s)审核人成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// ReviewerRevokedArgs 收回审核人事件参数
type ReviewerRevokedArgs struct {
	Manager models.User
	User    models.User
}

// ReviewerRevokedLogWriter 收回审核人日志记录
func ReviewerRevokedLogWriter(e *Event) {
	if args, ok := e.Args.(*ReviewerRevokedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)收回用户(uuid=%s)审核人成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// ClusterGrantedArgs 授权群集事件参数
type ClusterGrantedArgs struct {
	Manager models.User
	User    models.User
}

// ClusterGrantedLogWriter 授权群集日志记录
func ClusterGrantedLogWriter(e *Event) {
	if args, ok := e.Args.(*ClusterGrantedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)授权用户(uuid=%s)群集成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// ClusterRevokedArgs 收回群集事件参数
type ClusterRevokedArgs struct {
	Manager models.User
	User    models.User
}

// ClusterRevokedLogWriter 收回群集日志记录
func ClusterRevokedLogWriter(e *Event) {
	if args, ok := e.Args.(*ClusterRevokedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)收回用户(uuid=%s)群集成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// RoleGrantedArgs 授权角色事件参数
type RoleGrantedArgs struct {
	Manager models.User
	User    models.User
}

// RoleGrantedLogWriter 授权角色日志记录
func RoleGrantedLogWriter(e *Event) {
	if args, ok := e.Args.(*RoleGrantedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)授权用户(uuid=%s)角色成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// RoleRevokedArgs 收回角色事件参数
type RoleRevokedArgs struct {
	Manager models.User
	User    models.User
}

// RoleRevokedLogWriter 收回角色日志记录
func RoleRevokedLogWriter(e *Event) {
	if args, ok := e.Args.(*RoleRevokedArgs); ok {
		LogWriter(args.Manager.UserID, fmt.Sprintf("管理员(uuid=%s)收回用户(uuid=%s)角色成功。\n", args.Manager.UUID, args.User.UUID))
	}
}

// UpdateStatistics 统计信息更新
func UpdateStatistics(group string, key string, value float64) error {
	session := g.Engine.NewSession()
	stats := &models.Statistic{}
	if ok, err := session.Where("`group` = ? AND `key` = ?", group, key).ForUpdate().Get(stats); ok {
		stats.Value = stats.Value + value
		if _, err := session.Where("`group` = ? AND `key` = ?", group, key).Update(stats); err != nil {
			session.Rollback()
			return err
		}
	} else if err != nil {
		session.Rollback()
		return err
	} else if !ok {
		stats.Group = group
		stats.Key = key
		stats.Value = value
		if _, err := session.Insert(stats); err != nil {
			session.Rollback()
			return err
		}
	}
	if err := session.Commit(); err != nil {
		return err
	}
	return nil
}

// LogWriter 日志记录表操作
func LogWriter(userID uint, operation string) error {
	if _, err := g.Engine.Insert(&models.Log{
		UserID:    userID,
		Operation: operation,
	}); err != nil {
		return fmt.Errorf("[E] 错误代码: %s, 错误信息: %s", gqlapi.ReturnCodeUnknowError, err.Error())
	}
	return nil
}

// MailSendArgs 邮件发送事件参数
type MailSendArgs struct {
	To      mail.Address
	Subject string
	Body    string
}

// MailSender 邮件发送
func MailSender(args MailSendArgs) {
	if !g.Config().Mail.Enabled {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				log.Errorf("[E] Send mail failed: %s", e.Error())
			} else {
				log.Errorf("[E] Send mail failed: %+v", err)
			}
		}
	}()

	from := mail.Address{Name: "系统用户", Address: g.Config().Mail.User}
	to := args.To

	// set headers for html email
	header := textproto.MIMEHeader{}
	header.Set(textproto.CanonicalMIMEHeaderKey("from"), fmt.Sprintf("%s <%s>", from.Name, from.Address))
	header.Set(textproto.CanonicalMIMEHeaderKey("to"), fmt.Sprintf("%s <%s>", to.Name, to.Address))
	header.Set(textproto.CanonicalMIMEHeaderKey("content-type"), "text/plain; charset=UTF-8")
	header.Set(textproto.CanonicalMIMEHeaderKey("mime-version"), "1.0")
	header.Set(textproto.CanonicalMIMEHeaderKey("subject"), args.Subject)

	// init empty message
	var buffer bytes.Buffer

	// write header
	for key, value := range header {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	// write body
	buffer.WriteString(fmt.Sprintf("\r\n%s", args.Body))

	// Connect to the SMTP Server
	servername := g.Config().Mail.Addr

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", g.Config().Mail.User, g.Config().Mail.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		panic(err)
	}

	_, err = w.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	c.Quit()
}

func renderer(id g.Template, data interface{}) (string, string) {
	tpl := caches.TemplatesMap.Any(func(elem *models.Template) bool {
		if elem.UUID == string(id) {
			return true
		}
		return false
	})
	if tpl == nil {
		return "", ""
	}
	fmap := template.FuncMap{
		"formatDate": formatDate,
	}
	tplSubject := template.Must(template.New("subject").Parse(tpl.Subject))
	tplBody := template.Must(template.New("body").Funcs(fmap).Parse(tpl.Body))
	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)
	tplSubject.Execute(subject, data)
	tplBody.Execute(body, data)
	return subject.String(), body.String()
}

func formatDate(ts uint) string {
	return time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05")
}
