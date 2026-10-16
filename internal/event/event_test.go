package event

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// TestFireNoDeadlockNoLeak guards the E1/E2 fix: firing an event far past the
// emitter channel capacity must neither deadlock nor leak goroutines, and the
// audit handler must run for every event (audit is never dropped).
func TestFireNoDeadlockNoLeak(t *testing.T) {
	var calls int64
	auditEE.On("test.fire.audit", func(e *Event) { atomic.AddInt64(&calls, 1) })

	before := runtime.NumGoroutine()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 5000; i++ {
			Fire("test.fire.audit", i)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Fire deadlocked: unbounded channel delivery not suppressed")
	}

	time.Sleep(200 * time.Millisecond)
	if after := runtime.NumGoroutine(); after-before > fireConcurrency+10 {
		t.Fatalf("goroutine leak: before=%d after=%d", before, after)
	}
	if got := atomic.LoadInt64(&calls); got != 5000 {
		t.Fatalf("audit handler ran %d times, want 5000", got)
	}
}

// TestAuditSynchronous guards the E3 fix: audit handlers run in the caller
// goroutine, so FireSync returns only after the audit write is attempted.
func TestAuditSynchronous(t *testing.T) {
	ran := false
	auditEE.On("test.fire.sync", func(e *Event) { ran = true })
	FireSync("test.fire.sync", nil)
	if !ran {
		t.Fatal("audit handler did not run synchronously")
	}
}
