package event

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/service"
)

// GlobalAuditLogWriter returns an audit handler that persists the event's
// AuditRecord. A nil record means the service layer already wrote the audit
// line inside its transaction, so the handler is a no-op for that event.
func GlobalAuditLogWriter(repo domain.LogRepository) Handler {
	return func(ctx context.Context, ev Event) error {
		rec := ev.AuditRecord()
		if rec == nil {
			return nil
		}
		if _, err := repo.Create(ctx, rec); err != nil {
			log.Errorf("[E] audit log write failed, topic=%s: %s", ev.Topic(), err.Error())
			return fmt.Errorf("event: audit %s: %w", ev.Topic(), err)
		}
		return nil
	}
}

// GlobalMessageWriter turns the events a user cares about into in-app
// messages.
//
// The channel split is deliberate: the in-app message goes to the party that
// has to act or to know the result, while the mail channel keeps the offline
// record (see ADR 0014). The actor never receives a message about their own
// action.
func GlobalMessageWriter(repo domain.MessageRepository) Handler {
	return func(ctx context.Context, ev Event) error {
		var (
			ownerID    uint
			reviewerID uint
			actorID    uint
			ticket     models.Ticket
			title      string
			content    string
		)

		switch e := ev.(type) {
		case *TicketCreated:
			// 新工单是审核人的待办，提交人自己知道。
			reviewerID, actorID, ticket = e.Ticket.ReviewerID, e.User.UserID, e.Ticket
			title = "待审核工单"
			content = fmt.Sprintf("工单「%s」等待审核。", ticket.Subject)
		case *CommentCreated:
			// 评论是给提交人的补充说明。
			ownerID, actorID, ticket = e.Ticket.UserID, e.User.UserID, e.Ticket
			title = "工单有新评论"
			content = fmt.Sprintf("工单「%s」收到新的评论。", ticket.Subject)
		case *TicketStatusPatched:
			ownerID, actorID, ticket = e.Ticket.UserID, e.User.UserID, e.Ticket
			title = "工单状态更新"
			content = fmt.Sprintf("工单「%s」的状态已更新。", ticket.Subject)
		case *TicketExecuted:
			ownerID, actorID, ticket = e.Ticket.UserID, e.UserID, e.Ticket
			title = "工单执行完成"
			content = fmt.Sprintf("工单「%s」已执行完成。", ticket.Subject)
		case *TicketFailed:
			ownerID, actorID, ticket = e.Ticket.UserID, e.UserID, e.Ticket
			title = "工单执行失败"
			content = fmt.Sprintf("工单「%s」执行失败。", ticket.Subject)
		default:
			return nil
		}

		recipient := ownerID
		if recipient == 0 {
			recipient = reviewerID
		}
		if recipient == 0 || recipient == actorID {
			return nil
		}

		if err := repo.Create(ctx, &domain.Message{
			UUID:     uuid.NewString(),
			UserID:   recipient,
			Category: "ticket",
			Title:    title,
			Content:  content,
			Link:     "/tickets/" + ticket.UUID,
		}); err != nil {
			log.Errorf("[E] message write failed, topic=%s: %s", ev.Topic(), err.Error())
			return fmt.Errorf("event: message %s: %w", ev.Topic(), err)
		}
		return nil
	}
}

// GlobalMetricsUpdater returns a metrics handler covering every event that
// maintains counters or gauges.
func GlobalMetricsUpdater(repo domain.MetricRepository) Handler {
	return func(ctx context.Context, ev Event) error {
		switch e := ev.(type) {
		case *TicketCreated:
			// 每日计数在事件发生的当下累加，历史数值才不会随状态推进而回溯变化。
			if err := bumpDaily(ctx, repo, service.MetricTicketsCreated); err != nil {
				return err
			}
			return ticketStats(ctx, repo, e.User.UUID, 1)
		case *TicketRemoved:
			if err := bumpDaily(ctx, repo, service.MetricTicketsCreated, -1); err != nil {
				return err
			}
			return ticketStats(ctx, repo, e.User.UUID, -1)
		case *QueryCreated:
			return queryStats(ctx, repo, e.User.UUID, "total-queries")
		case *QueryAnalyzed:
			return queryStats(ctx, repo, e.User.UUID, "total-analyses")
		case *QueryRewrited:
			return queryStats(ctx, repo, e.User.UUID, "total-rewrites")
		case *UserRegistered:
			return bump(ctx, repo, overallKey("total-users"), 1)
		case *UserCreated:
			return bump(ctx, repo, overallKey("total-users"), 1)
		case *UserSignedIn:
			if err := bump(ctx, repo, userKey(e.User.UUID, "total-logins"), 1); err != nil {
				return err
			}
			// latest-login is a gauge: it holds the last sign-in time instead of
			// accumulating it.
			return setGauge(ctx, repo, userKey(e.User.UUID, "latest-login"), float64(time.Now().Unix()))
		case *CommentCreated:
			if err := bump(ctx, repo, overallKey("total-comments"), 1); err != nil {
				return err
			}
			return bump(ctx, repo, userKey(e.User.UUID, "total-comments"), 1)
		case *TicketExecuted:
			if err := bumpDaily(ctx, repo, service.MetricTicketsExecuted); err != nil {
				return err
			}
			return ticketStats(ctx, repo, "", 1)
		case *TicketFailed:
			if err := bumpDaily(ctx, repo, service.MetricTicketsFailed); err != nil {
				return err
			}
			return ticketStats(ctx, repo, "", 1)
		case *ClusterCreated:
			return bump(ctx, repo, overallKey("total-clusters"), 1)
		case *ClusterRemoved:
			return bump(ctx, repo, overallKey("total-clusters"), -1)
		default:
			// Events without counters are a no-op.
			return nil
		}
	}
}

// ticketStats updates the overall and per-user ticket counters. The daily
// counters are written by the event handler, where the day of the fact is known.
func ticketStats(ctx context.Context, repo domain.MetricRepository, userUUID string, delta float64) error {
	if err := bump(ctx, repo, overallKey("total-tickets"), delta); err != nil {
		return err
	}
	if userUUID == "" {
		return nil
	}
	return bump(ctx, repo, userKey(userUUID, "total-tickets"), delta)
}

// queryStats updates the overall, daily and per-user query counters. counter is
// the per-user metric name ("total-queries" / "total-analyses" / "total-rewrites").
func queryStats(ctx context.Context, repo domain.MetricRepository, userUUID, counter string) error {
	if err := bump(ctx, repo, overallKey("total-queries"), 1); err != nil {
		return err
	}
	if err := bump(ctx, repo, dailyKey("total-queries", time.Now().Format("2006-01-02")), 1); err != nil {
		return err
	}
	return bump(ctx, repo, userKey(userUUID, counter), 1)
}

// overallKey builds the key of a platform wide counter.
func overallKey(name string) domain.MetricKey {
	return domain.MetricKey{Name: name, Scope: domain.MetricScopeOverall}
}

// dailyKey builds the key of a counter inside a time bucket.
func dailyKey(name, bucket string) domain.MetricKey {
	return domain.MetricKey{Name: name, Scope: domain.MetricScopeDaily, Bucket: bucket}
}

// userKey builds the key of a counter or gauge scoped to a user.
func userKey(userUUID, name string) domain.MetricKey {
	return domain.MetricKey{Name: name, Scope: domain.MetricScopeUser, ScopeID: userUUID}
}

// bump applies a single counter increment, logging (not failing) errors to
// match the former best-effort behaviour.
func bump(ctx context.Context, repo domain.MetricRepository, key domain.MetricKey, delta float64) error {
	if err := repo.Increment(ctx, key, delta); err != nil {
		log.Warnf("[W] metric increment failed, name=%s scope=%d: %s", key.Name, key.Scope, err.Error())
	}
	return nil
}

// setGauge overwrites a gauge metric with the given value.
func setGauge(ctx context.Context, repo domain.MetricRepository, key domain.MetricKey, value float64) error {
	if err := repo.SetGauge(ctx, key, value); err != nil {
		log.Warnf("[W] metric gauge write failed, name=%s scope=%d: %s", key.Name, key.Scope, err.Error())
	}
	return nil
}

// bumpDaily accumulates a daily counter for the current day. The metric name
// lives in the service package, where the rollup job writes the other daily
// buckets.
func bumpDaily(ctx context.Context, repo domain.MetricRepository, name string, delta ...float64) error {
	value := 1.0
	if len(delta) > 0 {
		value = delta[0]
	}
	key := domain.MetricKey{
		Name:   name,
		Scope:  domain.MetricScopeDaily,
		Bucket: time.Now().Format("2006-01-02"),
	}
	if err := repo.Increment(ctx, key, value); err != nil {
		log.Warnf("[W] daily metric increment failed, name=%s: %s", name, err.Error())
	}
	return nil
}
