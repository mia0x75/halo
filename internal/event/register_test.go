package event

import (
	"context"
	"testing"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

type fakeLogRepo struct{ created []*domain.Log }

func (f *fakeLogRepo) Create(_ context.Context, l *domain.Log) (int64, error) {
	f.created = append(f.created, l)
	return int64(len(f.created)), nil
}
func (f *fakeLogRepo) GetByID(context.Context, int64) (*domain.Log, error)     { return nil, nil }
func (f *fakeLogRepo) GetByUUID(context.Context, string) (*domain.Log, error)  { return nil, nil }
func (f *fakeLogRepo) Update(context.Context, *domain.Log) error               { return nil }
func (f *fakeLogRepo) Delete(context.Context, int64) error                     { return nil }
func (f *fakeLogRepo) DeleteByUUID(context.Context, string) error              { return nil }
func (f *fakeLogRepo) List(context.Context, int64, int) ([]*domain.Log, error) { return nil, nil }

type fakeMetricRepo struct {
	keys    []string
	gauges  []string
	counter map[string]float64
}

func (f *fakeMetricRepo) Create(context.Context, *domain.Metric) error { return nil }
func (f *fakeMetricRepo) GetByUUID(context.Context, string) (*domain.Metric, error) {
	return nil, nil
}
func (f *fakeMetricRepo) GetByKey(context.Context, domain.MetricKey) (*domain.Metric, error) {
	return nil, nil
}
func (f *fakeMetricRepo) Update(context.Context, *domain.Metric) error { return nil }
func (f *fakeMetricRepo) DeleteByUUID(context.Context, string) error   { return nil }
func (f *fakeMetricRepo) List(context.Context, string, int) ([]*domain.Metric, error) {
	return nil, nil
}
func (f *fakeMetricRepo) Increment(_ context.Context, key domain.MetricKey, delta float64) error {
	f.keys = append(f.keys, key.Name)
	if f.counter == nil {
		f.counter = map[string]float64{}
	}
	f.counter[key.Name] += delta
	return nil
}
func (f *fakeMetricRepo) DeleteBucketsBefore(context.Context, string) (int64, error) {
	return 0, nil
}
func (f *fakeMetricRepo) SetGauge(_ context.Context, key domain.MetricKey, value float64) error {
	f.gauges = append(f.gauges, key.Name)
	if f.counter == nil {
		f.counter = map[string]float64{}
	}
	f.counter[key.Name] = value
	return nil
}

type fakeUserRepo struct{}

func (fakeUserRepo) Create(context.Context, *domain.User) (int64, error) { return 0, nil }
func (fakeUserRepo) GetByID(context.Context, int64) (*domain.User, error) {
	return &domain.User{UUID: "u", Name: "n", Email: "e@x"}, nil
}
func (fakeUserRepo) GetByUUID(context.Context, string) (*domain.User, error) {
	return &domain.User{UUID: "u"}, nil
}
func (fakeUserRepo) Update(context.Context, *domain.User) error               { return nil }
func (fakeUserRepo) Delete(context.Context, int64) error                      { return nil }
func (fakeUserRepo) DeleteByUUID(context.Context, string) error               { return nil }
func (fakeUserRepo) List(context.Context, int64, int) ([]*domain.User, error) { return nil, nil }

// TestRegisterAuditGating verifies that GlobalAuditLogWriter persists a row for
// events with a non-nil AuditRecord and skips events whose record is nil
// (already audited by the service layer).
func TestRegisterAuditGating(t *testing.T) {
	logs := &fakeLogRepo{}
	bus := NewBus(4)
	Register(bus, Deps{Log: logs, Metric: &fakeMetricRepo{}, User: fakeUserRepo{}})

	// Non-nil record: a login must be audited by the event layer.
	if err := bus.Emit(context.Background(), EventUserSignedIn, &UserSignedIn{User: models.User{UserID: 7, UUID: "u7"}, IP: "127.0.0.1"}); err != nil {
		t.Fatalf("emit signed-in: %v", err)
	}
	if len(logs.created) != 1 {
		t.Fatalf("signed-in audit rows = %d, want 1", len(logs.created))
	}
	if logs.created[0].UserID != 7 {
		t.Fatalf("audit user id = %d, want 7", logs.created[0].UserID)
	}

	// Nil record: ticket creation is audited by the service layer, so the event
	// layer must not write a second row.
	if err := bus.Emit(context.Background(), EventTicketCreated, &TicketCreated{Ticket: models.Ticket{UUID: "t"}}); err != nil {
		t.Fatalf("emit created: %v", err)
	}
	if len(logs.created) != 1 {
		t.Fatalf("audit rows after nil-record event = %d, want 1", len(logs.created))
	}
}

// TestRegisterMetricsDispatch verifies the wildcard metrics handler keeps the
// per-event counter behaviour.
func TestRegisterMetricsDispatch(t *testing.T) {
	metrics := &fakeMetricRepo{}
	bus := NewBus(4)
	Register(bus, Deps{Log: &fakeLogRepo{}, Metric: metrics, User: fakeUserRepo{}})

	bus.OnNotify("test.done", func(context.Context, Event) error { return nil })
	if err := bus.Fire(context.Background(), EventTicketCreated, &TicketCreated{User: models.User{UUID: "u1"}}); err != nil {
		t.Fatalf("fire: %v", err)
	}
	bus.Close()

	got := map[string]bool{}
	for _, k := range metrics.keys {
		got[k] = true
	}
	for _, want := range []string{"total-tickets"} {
		if !got[want] {
			t.Errorf("missing metrics key %q (got %v)", want, metrics.keys)
		}
	}
	// The ticket counter is maintained on the overall, daily and user scopes.
	if len(metrics.keys) != 3 {
		t.Errorf("expected 3 counter updates, got %v", metrics.keys)
	}
}
