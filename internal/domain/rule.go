package domain

// RuleLevel is the severity of a validation rule. Matches the values written
// to mm_rules.level by the seed data.
type RuleLevel int

const (
	RuleLevelMessage RuleLevel = 0
	RuleLevelError   RuleLevel = 1
	RuleLevelWarning RuleLevel = 2
)

// Rule is a row of mm_rules. The primary key is name; UUID is the stable lookup
// key used by the repository. Group, VldrGroup and Bitwise are seed-driven
// categories, so they stay numeric until the rule engine is ported.
type Rule struct {
	Name        string    `db:"name"`
	UUID        string    `db:"uuid"`
	Group       uint8     `db:"group"`
	Description string    `db:"description"`
	Level       RuleLevel `db:"level"`
	VldrGroup   uint8     `db:"vldr_group"`
	Operator    string    `db:"operator"`
	Values      string    `db:"values"`
	Bitwise     uint8     `db:"bitwise"`
	Func        string    `db:"func"`
	Message     string    `db:"message"`
	Element     string    `db:"element"`
	Version     int       `db:"version"`
	UpdateAt    *uint     `db:"update_at"`
	CreateAt    uint      `db:"create_at"`
}
