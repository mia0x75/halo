package event

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

// testEvent is a minimal Event for bus tests.
type testEvent struct{ TopicName string }

func (e testEvent) Topic() string            { return e.TopicName }
func (e testEvent) AuditRecord() *domain.Log { return nil }

// TestFireNoDeadlockNoLeak guards the E1/E2 fix: firing an event far past the
// notification concurrency bound must neither deadlock nor leak goroutines, and
// the audit handler must run for every event (audit is never dropped).
func TestFireNoDeadlockNoLeak(t *testing.T) {
	bus := NewBus(10)
	var calls int64
	bus.OnAudit("test.fire.audit", func(context.Context, Event) error {
		atomic.AddInt64(&calls, 1)
		return nil
	})

	before := runtime.NumGoroutine()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 5000; i++ {
			_ = bus.Fire(context.Background(), "test.fire.audit", testEvent{"test.fire.audit"})
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Fire deadlocked")
	}

	bus.Close()
	time.Sleep(100 * time.Millisecond)
	if after := runtime.NumGoroutine(); after-before > 10+10 {
		t.Fatalf("goroutine leak: before=%d after=%d", before, after)
	}
	if got := atomic.LoadInt64(&calls); got != 5000 {
		t.Fatalf("audit handler ran %d times, want 5000", got)
	}
}

// TestAuditSynchronous guards the E3 fix: audit handlers run in the caller
// goroutine, so Emit returns only after the audit handler ran.
func TestAuditSynchronous(t *testing.T) {
	bus := NewBus(4)
	ran := false
	bus.OnAudit("test.fire.sync", func(context.Context, Event) error {
		ran = true
		return nil
	})
	if err := bus.Emit(context.Background(), "test.fire.sync", testEvent{"test.fire.sync"}); err != nil {
		t.Fatalf("Emit: %v", err)
	}
	if !ran {
		t.Fatal("audit handler did not run synchronously")
	}
}
