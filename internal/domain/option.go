package domain

import (
	"fmt"
	"strconv"
)

// Option is a row of mm_options. The primary key is name; UUID is the stable
// lookup key used by the repository.
type Option struct {
	Name        string `db:"name"`
	UUID        string `db:"uuid"`
	Parameter   string `db:"parameter"`
	Description string `db:"description"`
	Element     string `db:"element"`
	Writable    uint8  `db:"writable"`
	Version     int    `db:"version"`
	UpdateAt    *uint  `db:"update_at"`
	CreateAt    uint   `db:"create_at"`
}

// OptionSpec describes how one platform option is edited and validated.
//
// The seeded options are legacy SMTP/LDAP entries that the server never reads;
// the spec still pins their control type so the settings page renders a switch
// for booleans and never echoes a password back.
type OptionSpec struct {
	// Element is the control the settings page renders: input, switch, password.
	Element string
	// Writable reports whether the option may be changed from the settings page.
	Writable bool
	// Secret marks values that must never be returned in clear text.
	Secret bool
	// Default is the value used when the stored parameter is empty.
	Default string
	// Min and Max bound the integer values of a number option.
	Min int
	Max int
}

// optionSpecs pins the behaviour of the known option names.
var optionSpecs = map[string]OptionSpec{
	"smtp.enabled":    {Element: "switch", Writable: false, Default: "false"},
	"smtp.host":       {Element: "input", Writable: false},
	"smtp.port":       {Element: "input", Writable: false, Default: "465"},
	"smtp.user":       {Element: "input", Writable: false},
	"smtp.password":   {Element: "password", Writable: false, Secret: true},
	"smtp.encryption": {Element: "input", Writable: false, Default: "ssl"},
	"ldap.enabled":    {Element: "switch", Writable: false, Default: "false"},
	"ldap.host":       {Element: "input", Writable: false},
	"ldap.domain":     {Element: "input", Writable: false},
	"ldap.type":       {Element: "input", Writable: false, Default: "ldap"},
	"ldap.user":       {Element: "input", Writable: false},
	"ldap.password":   {Element: "password", Writable: false, Secret: true},
	"ldap.sc":         {Element: "input", Writable: false},
	"ldap.ou":         {Element: "input", Writable: false},
	// 保留期以天为单位，由 halocli retention cleanup 读取。
	"retention.metrics-days":    {Element: "number", Writable: true, Default: "90", Min: 1, Max: 3650},
	"retention.hits-days":       {Element: "number", Writable: true, Default: "90", Min: 1, Max: 3650},
	"retention.logs-days":       {Element: "number", Writable: true, Default: "180", Min: 1, Max: 3650},
	"retention.messages-days":   {Element: "number", Writable: true, Default: "90", Min: 1, Max: 3650},
	"retention.tasks-days":      {Element: "number", Writable: true, Default: "30", Min: 1, Max: 3650},
	"mail.enabled":              {Element: "switch", Writable: true, Default: "false"},
	"mail.addr":                 {Element: "input", Writable: true},
	"mail.user":                 {Element: "input", Writable: true},
	"mail.password":             {Element: "password", Writable: true, Secret: true},
	"mail.insecure-skip-verify": {Element: "switch", Writable: true, Default: "false"},
}

// OptionSpecFor returns the spec of an option name. Unknown names fall back to
// a plain read-only text field.
func OptionSpecFor(name string) OptionSpec {
	if spec, ok := optionSpecs[name]; ok {
		return spec
	}
	return OptionSpec{Element: "input"}
}

// ValidateOptionValue checks a new value against the option spec.
func ValidateOptionValue(name, value string) error {
	spec := OptionSpecFor(name)
	switch spec.Element {
	case "switch":
		if value != "true" && value != "false" {
			return fmt.Errorf("%w: 选项 %s 只接受 true 或 false", ErrInvalidInput, name)
		}
	case "number":
		n, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%w: 选项 %s 的值必须是整数", ErrInvalidInput, name)
		}
		if spec.Max > spec.Min && (n < spec.Min || n > spec.Max) {
			return fmt.Errorf("%w: 选项 %s 的值必须在 %d 到 %d 之间", ErrInvalidInput, name, spec.Min, spec.Max)
		}
	case "input", "password":
		if len(value) > 75 {
			return fmt.Errorf("%w: 选项 %s 的值超过 75 个字符", ErrInvalidInput, name)
		}
	}
	return nil
}

// MaskOptionValue hides a secret value behind a fixed placeholder.
func MaskOptionValue(name, value string) string {
	if !OptionSpecFor(name).Secret {
		return value
	}
	if value == "" {
		return ""
	}
	return "********"
}
