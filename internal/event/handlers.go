package event

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/domain"
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

// GlobalMetricsUpdater returns a statistics handler covering every event that
// maintains counters. The per-event group/key/delta behaviour mirrors the
// former per-event *StatisticUpdater handlers.
func GlobalMetricsUpdater(repo domain.StatisticRepository) Handler {
	return func(ctx context.Context, ev Event) error {
		switch e := ev.(type) {
		case *TicketCreated:
			return ticketStats(ctx, repo, e.User.UUID, 1)
		case *TicketRemoved:
			return ticketStats(ctx, repo, e.User.UUID, -1)
		case *QueryCreated:
			return queryStats(ctx, repo, e.User.UUID, "total-queries")
		case *QueryAnalyzed:
			return queryStats(ctx, repo, e.User.UUID, "total-analyses")
		case *QueryRewrited:
			return queryStats(ctx, repo, e.User.UUID, "total-rewrites")
		case *UserRegistered:
			return bump(ctx, repo, "overall", "total-users", 1)
		case *UserCreated:
			return bump(ctx, repo, "overall", "total-users", 1)
		case *UserSignedIn:
			if err := bump(ctx, repo, e.User.UUID, "total-logins", 1); err != nil {
				return err
			}
			return bump(ctx, repo, e.User.UUID, "latest-login", float64(time.Now().Unix()))
		case *CommentCreated:
			if err := bump(ctx, repo, "overall", "total-comments", 1); err != nil {
				return err
			}
			return bump(ctx, repo, e.User.UUID, "total-comments", 1)
		case *ClusterCreated:
			return bump(ctx, repo, "overall", "total-clusters", 1)
		case *ClusterRemoved:
			return bump(ctx, repo, "overall", "total-clusters", -1)
		default:
			// Events without counters are a no-op.
			return nil
		}
	}
}

// ticketStats updates the overall, daily and per-user ticket counters.
func ticketStats(ctx context.Context, repo domain.StatisticRepository, userUUID string, delta float64) error {
	if err := bump(ctx, repo, "overall", "total-tickets", delta); err != nil {
		return err
	}
	today := time.Now().Format("2006-01-02")
	if err := bump(ctx, repo, "tickets-daily", today, delta); err != nil {
		return err
	}
	return bump(ctx, repo, userUUID, "total-tickets", delta)
}

// queryStats updates the overall, daily and per-user query counters. counter is
// the per-user key ("total-queries" / "total-analyses" / "total-rewrites").
func queryStats(ctx context.Context, repo domain.StatisticRepository, userUUID, counter string) error {
	if err := bump(ctx, repo, "overall", "total-queries", 1); err != nil {
		return err
	}
	today := time.Now().Format("2006-01-02")
	if err := bump(ctx, repo, "queries-daily", today, 1); err != nil {
		return err
	}
	return bump(ctx, repo, userUUID, counter, 1)
}

// bump applies a single statistic increment, logging (not failing) errors to
// match the former best-effort behaviour.
func bump(ctx context.Context, repo domain.StatisticRepository, group, key string, delta float64) error {
	if err := repo.UpsertIncrement(ctx, group, key, delta); err != nil {
		log.Warnf("[W] statistics upsert failed, group=%s key=%s: %s", group, key, err.Error())
	}
	return nil
}
