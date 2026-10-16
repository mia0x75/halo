package event

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// Handler processes one event. Errors from audit handlers are aggregated and
// returned by Emit/Fire; errors from notify handlers are logged and dropped.
type Handler func(ctx context.Context, ev Event) error

// Bus is an in-process event bus. Audit handlers run synchronously so a write
// is recorded before the caller returns; notification and statistics handlers
// run asynchronously under a bounded limiter and may be dropped when saturated.
type Bus struct {
	mu     sync.RWMutex
	audit  map[string][]Handler
	notify map[string][]Handler
	sem    chan struct{}
	wg     sync.WaitGroup
	closed atomic.Bool
}

// NewBus returns a bus whose notification concurrency is bounded by
// notifyConcurrency. A non-positive value falls back to 100.
func NewBus(notifyConcurrency int) *Bus {
	if notifyConcurrency <= 0 {
		notifyConcurrency = 100
	}
	return &Bus{
		audit:  make(map[string][]Handler),
		notify: make(map[string][]Handler),
		sem:    make(chan struct{}, notifyConcurrency),
	}
}

// OnAudit registers a synchronous audit handler for a topic. The wildcard "*"
// matches every topic.
func (b *Bus) OnAudit(topic string, h Handler) {
	if h == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.audit[topic] = append(b.audit[topic], h)
}

// OnNotify registers an asynchronous notification handler for a topic. The
// wildcard "*" matches every topic.
func (b *Bus) OnNotify(topic string, h Handler) {
	if h == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.notify[topic] = append(b.notify[topic], h)
}

// Emit runs audit handlers synchronously and returns their aggregated error.
// Notifications are not triggered.
func (b *Bus) Emit(ctx context.Context, topic string, ev Event) error {
	if ev == nil {
		return errors.New("event: nil event")
	}
	var errs []error
	for _, h := range b.auditHandlers(topic) {
		if err := h(ctx, ev); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Fire runs audit handlers synchronously and notifications asynchronously under
// the bounded limiter. It returns the audit error, if any.
func (b *Bus) Fire(ctx context.Context, topic string, ev Event) error {
	err := b.Emit(ctx, topic, ev)
	b.dispatchNotify(topic, ev)
	return err
}

// FireSync runs both audit and notification handlers synchronously.
func (b *Bus) FireSync(ctx context.Context, topic string, ev Event) error {
	err := b.Emit(ctx, topic, ev)
	if ev == nil {
		return err
	}
	for _, h := range b.notifyHandlers(topic) {
		if nerr := h(ctx, ev); nerr != nil {
			log.Errorf("[E] notify handler failed, topic=%s: %s", topic, nerr.Error())
		}
	}
	return err
}

// dispatchNotify queues notifications on a bounded limiter. When saturated the
// event is dropped rather than blocking the caller or spawning unbounded
// goroutines.
func (b *Bus) dispatchNotify(topic string, ev Event) {
	if ev == nil || b.closed.Load() {
		return
	}
	handlers := b.notifyHandlers(topic)
	if len(handlers) == 0 {
		return
	}
	select {
	case b.sem <- struct{}{}:
	default:
		log.Warnf("[W] event notify limiter saturated, dropping topic=%s", topic)
		return
	}
	b.wg.Add(1)
	go func() {
		defer func() { <-b.sem; b.wg.Done() }()
		ctx := context.Background()
		for _, h := range handlers {
			if err := h(ctx, ev); err != nil {
				log.Errorf("[E] notify handler failed, topic=%s: %s", topic, err.Error())
			}
		}
	}()
}

// Close waits for in-flight notifications to finish. After Close, Fire drops
// notifications instead of queuing them.
func (b *Bus) Close() {
	b.closed.Store(true)
	b.wg.Wait()
}

func (b *Bus) auditHandlers(topic string) []Handler {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return collect(b.audit, topic)
}

func (b *Bus) notifyHandlers(topic string) []Handler {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return collect(b.notify, topic)
}

// collect returns the handlers registered for the exact topic followed by the
// wildcard handlers, preserving registration order within each bucket.
func collect(m map[string][]Handler, topic string) []Handler {
	out := make([]Handler, 0, len(m[topic])+len(m["*"]))
	out = append(out, m[topic]...)
	out = append(out, m["*"]...)
	return out
}
