package domain

import (
	"fmt"
	"strconv"
)

// Option is a row of mm_options. The primary key is name; UUID is the stable
// lookup key used by the repository.
//
// The presentation metadata (control type, group, order, dependency) lives in
// the row itself, so the settings page and the API always agree on how an
// option is edited.
type Option struct {
	Name        string `db:"name"`
	UUID        string `db:"uuid"`
	Parameter   string `db:"parameter"`
	Description string `db:"description"`
	Element     string `db:"element"`
	Writable    uint8  `db:"writable"`
	GroupName   string `db:"group_name"`
	SortOrder   int    `db:"sort_order"`
	DependsOn   string `db:"depends_on"`
	MinValue    *int   `db:"min_value"`
	MaxValue    *int   `db:"max_value"`
	Version     int    `db:"version"`
	UpdateAt    *uint  `db:"update_at"`
	CreateAt    uint   `db:"create_at"`
}

// The control types the settings page knows how to render.
const (
	OptionElementInput    = "input"
	OptionElementNumber   = "number"
	OptionElementSwitch   = "switch"
	OptionElementPassword = "password"
)

// Secret reports whether the value must never be returned in clear text.
func (o *Option) Secret() bool {
	return o != nil && o.Element == OptionElementPassword
}

// Writable reports whether the option may be changed from the settings page.
func (o *Option) IsWritable() bool {
	return o != nil && o.Writable == 1
}

// ValidateOptionValue checks a new value against the option metadata.
func ValidateOptionValue(o *Option) error {
	switch o.Element {
	case OptionElementSwitch:
		if o.Parameter != "true" && o.Parameter != "false" {
			return fmt.Errorf("%w: 选项 %s 只接受 true 或 false", ErrInvalidInput, o.Name)
		}
	case OptionElementNumber:
		n, err := strconv.Atoi(o.Parameter)
		if err != nil {
			return fmt.Errorf("%w: 选项 %s 的值必须是整数", ErrInvalidInput, o.Name)
		}
		if o.MinValue != nil && n < *o.MinValue {
			return fmt.Errorf("%w: 选项 %s 的值不能小于 %d", ErrInvalidInput, o.Name, *o.MinValue)
		}
		if o.MaxValue != nil && n > *o.MaxValue {
			return fmt.Errorf("%w: 选项 %s 的值不能大于 %d", ErrInvalidInput, o.Name, *o.MaxValue)
		}
	case OptionElementInput, OptionElementPassword:
		if len(o.Parameter) > 75 {
			return fmt.Errorf("%w: 选项 %s 的值超过 75 个字符", ErrInvalidInput, o.Name)
		}
	}
	return nil
}

// MaskOptionValue hides a secret value behind a fixed placeholder.
func MaskOptionValue(o *Option) string {
	if o == nil {
		return ""
	}
	if !o.Secret() {
		return o.Parameter
	}
	if o.Parameter == "" {
		return ""
	}
	return "********"
}
