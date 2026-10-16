package event

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/notify"
)

// Deps are the collaborators the event handlers need. They are injected so the
// event package depends only on domain repositories and the notifier
// abstraction, never on the service or transport layers.
type Deps struct {
	Log       domain.LogRepository
	Statistic domain.StatisticRepository
	User      domain.UserRepository
	Notifier  notify.Notifier
}

// Register wires the standard handlers onto the bus and records it as the
// default bus used by the package-level Fire helpers.
func Register(bus *Bus, deps Deps) {
	// Global audit: one wildcard handler replaces the 34 per-event log writers.
	// Each event decides via AuditRecord whether it needs a line.
	bus.OnAudit("*", GlobalAuditLogWriter(deps.Log))

	// Global statistics: one wildcard handler replaces the 11 per-event
	// statistic updaters, dispatching on the concrete event type.
	bus.OnNotify("*", GlobalMetricsUpdater(deps.Statistic))

	// Mail: per-topic handlers (template + recipients differ per event). The
	// registry is local to Register so no package-level state leaks.
	registry := bind(deps.Notifier, deps.User)
	for topic, handler := range registry {
		bus.OnNotify(topic, handler)
	}

	setDefaultBus(bus)
}

// bind returns the mail routing table. It is a free function (not a method) and
// keeps its state in a local map so nothing is parked at package scope.
func bind(n notify.Notifier, userRepo domain.UserRepository) map[string]Handler {
	registry := make(map[string]Handler)
	for topic, h := range mailRoutes(n, userRepo) {
		if h != nil {
			registry[topic] = h
		}
	}
	return registry
}
