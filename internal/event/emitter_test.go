package event

import (
	"testing"
	"time"
)

// The vendored Emitter is a stripped-down pubsub: Args is the raw value passed
// to Emit and there are no Event helper accessors, so the old upstream test
// suite no longer compiles. These tests cover the subset the app relies on.

func receive(t *testing.T, ch <-chan Event) Event {
	t.Helper()
	select {
	case e := <-ch:
		return e
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
		return Event{}
	}
}

func TestEmitterEmitOn(t *testing.T) {
	ee := New(0, DefaultMatcher())
	ch := ee.On("test")

	go ee.Emit("test", 42)

	e := receive(t, ch)
	if e.Topic != "test" {
		t.Fatalf("topic = %q, want %q", e.Topic, "test")
	}
	if e.Args != 42 {
		t.Fatalf("args = %#v, want 42", e.Args)
	}
}

func TestEmitterOff(t *testing.T) {
	ee := New(0, DefaultMatcher())
	ee.On("test")
	if got := len(ee.Listeners("test")); got != 1 {
		t.Fatalf("listeners = %d, want 1", got)
	}
	if got := len(ee.Topics()); got != 1 {
		t.Fatalf("topics = %d, want 1", got)
	}

	ee.Off("test")
	if got := len(ee.Listeners("test")); got != 0 {
		t.Fatalf("listeners after Off = %d, want 0", got)
	}
	if got := len(ee.Topics()); got != 0 {
		t.Fatalf("topics after Off = %d, want 0", got)
	}
}
