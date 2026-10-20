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

// GlobalMetricsUpdater returns a metrics handler covering every event that
// maintains counters or gauges.
func GlobalMetricsUpdater(repo domain.MetricRepository) Handler {
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

// ticketStats updates the overall, daily and per-user ticket counters.
func ticketStats(ctx context.Context, repo domain.MetricRepository, userUUID string, delta float64) error {
	if err := bump(ctx, repo, overallKey("total-tickets"), delta); err != nil {
		return err
	}
	if err := bump(ctx, repo, dailyKey("total-tickets", time.Now().Format("2006-01-02")), delta); err != nil {
		return err
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
