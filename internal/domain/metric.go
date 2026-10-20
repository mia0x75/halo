package domain

// Metric scopes. A metric is identified by name + scope + scope id + bucket.
const (
	// MetricScopeOverall counts across the whole platform.
	MetricScopeOverall uint8 = iota
	// MetricScopeDaily counts per time bucket.
	MetricScopeDaily
	// MetricScopeUser counts per user.
	MetricScopeUser
	// MetricScopeCluster counts per cluster.
	MetricScopeCluster
	// MetricScopeGroup counts per cluster group.
	MetricScopeGroup
)

// Metric types. The write path differs: a counter is incremented atomically, a
// gauge is overwritten.
const (
	// MetricTypeCounter accumulates deltas.
	MetricTypeCounter uint8 = iota
	// MetricTypeGauge holds the latest value.
	MetricTypeGauge
)

// Metric is a row of mm_metrics. The natural key is
// (Name, Scope, ScopeID, Bucket): ScopeID holds the uuid of the scoped entity
// and stays empty for the overall scope, Bucket holds the time bucket key
// (for example 2026-10-20) and stays empty outside the daily scope.
type Metric struct {
	MetricID uint
	UUID     string
	Name     string
	Scope    uint8
	ScopeID  string
	Type     uint8
	Bucket   string
	Value    float64
	Version  int
	UpdateAt *uint
	CreateAt uint
}

// MetricKey identifies a metric row.
type MetricKey struct {
	Name    string
	Scope   uint8
	ScopeID string
	Bucket  string
}
