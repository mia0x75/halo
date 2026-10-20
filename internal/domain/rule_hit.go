package domain

// Violation is a single rule hit. Level is the raw mm_rules.level value
// (0=message, 1=error, 2=warning); Rule carries the name of the rule that was
// hit so the hits can be grouped by rule.
type Violation struct {
	// Rule is the name of the rule that was hit, empty when the engine does
	// not expose it.
	Rule string
	// Message is the rendered hit message.
	Message string
	// Level is the severity of the hit.
	Level uint8
}

// RuleHit is a persisted rule hit: one row of mm_violations.
type RuleHit struct {
	HitID     uint
	UUID      string
	TicketID  uint
	Sequence  uint16
	RuleName  string
	Level     uint8
	Message   string
	CreateAt  uint
	UpdatedAt *uint
}

// RuleHitCount is an aggregated number of hits, used by the dashboard
// summaries (by level, by rule).
type RuleHitCount struct {
	// Rule is the rule name, empty for the by-level summary.
	Rule string
	// Level is the severity, empty for the by-rule summary.
	Level uint8
	// Count is the number of hits.
	Count uint
}
