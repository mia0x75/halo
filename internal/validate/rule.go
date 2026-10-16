package validate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/mia0x75/halo/internal/domain"
)

// Rule is a validation rule as the engine sees it: the persisted domain entity
// plus lazily parsed, memoised Values. It embeds domain.Rule so call sites keep
// reading r.Name / r.Values / r.Message unchanged.
//
// The mutex guards the cache because validator groups run concurrently.
type Rule struct {
	domain.Rule

	mu    sync.Mutex
	cache map[string]any
}

// empty reports whether Values carries nothing to parse. Seeded "in"/"regexp"
// rules always have values; NULL rows (stored as the literal "nil") do not.
func (r *Rule) empty() bool {
	return r.Values == "" || r.Values == "nil"
}

// Int parses Values as a base-10 integer.
func (r *Rule) Int() (int, error) {
	if r.empty() {
		return 0, nil
	}
	if v, ok := r.cached("int"); ok {
		return v.(int), nil
	}
	n, err := strconv.Atoi(r.Values)
	if err != nil {
		return 0, fmt.Errorf("parse values %q as int: %w", r.Values, err)
	}
	r.setCache("int", n)
	return n, nil
}

// StringSlice parses Values as a JSON array of strings.
func (r *Rule) StringSlice() ([]string, error) {
	if r.empty() {
		return nil, nil
	}
	if v, ok := r.cached("strings"); ok {
		return v.([]string), nil
	}
	var s []string
	if err := json.Unmarshal([]byte(r.Values), &s); err != nil {
		return nil, fmt.Errorf("parse values %q as string slice: %w", r.Values, err)
	}
	r.setCache("strings", s)
	return s, nil
}

// Regexp compiles Values as a regular expression.
func (r *Rule) Regexp() (*regexp.Regexp, error) {
	if r.empty() {
		return nil, nil
	}
	if v, ok := r.cached("re"); ok {
		return v.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(r.Values)
	if err != nil {
		return nil, fmt.Errorf("compile values %q as regexp: %w", r.Values, err)
	}
	r.setCache("re", re)
	return re, nil
}

func (r *Rule) cached(key string) (any, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache == nil {
		return nil, false
	}
	v, ok := r.cache[key]
	return v, ok
}

func (r *Rule) setCache(key string, v any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache == nil {
		r.cache = map[string]any{}
	}
	r.cache[key] = v
}

// ruleFromDomain wraps a domain entity with an empty parse cache.
func ruleFromDomain(d *domain.Rule) *Rule {
	if d == nil {
		return nil
	}
	return &Rule{Rule: *d}
}

// ValidateRules checks that every rule is structurally sound and that its
// Values parse under its declared Operator. It runs once at startup; the
// per-Func registry check lands with V.2.
func ValidateRules(rules []*domain.Rule) error {
	seen := make(map[string]struct{}, len(rules))
	var errs []error
	for _, r := range rules {
		if r == nil {
			continue
		}
		if r.Name == "" {
			errs = append(errs, fmt.Errorf("validate: rule with empty name"))
			continue
		}
		if r.Func == "" {
			errs = append(errs, fmt.Errorf("validate: rule %s has empty func", r.Name))
		}
		if _, dup := seen[r.Name]; dup {
			errs = append(errs, fmt.Errorf("validate: duplicate rule %s", r.Name))
		}
		seen[r.Name] = struct{}{}

		wrapped := ruleFromDomain(r)
		switch r.Operator {
		case "in":
			if _, err := wrapped.StringSlice(); err != nil {
				errs = append(errs, fmt.Errorf("rule %s: %w", r.Name, err))
			}
		case "regexp":
			if _, err := wrapped.Regexp(); err != nil {
				errs = append(errs, fmt.Errorf("rule %s: %w", r.Name, err))
			}
		}
	}
	return joinErrors(errs)
}

// joinErrors collapses the collected validation errors into one, keeping the
// full list readable in a single startup log line.
func joinErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		msg := fmt.Sprintf("validate: %d rule problems", len(errs))
		for _, e := range errs {
			msg += "; " + e.Error()
		}
		return fmt.Errorf("%s", msg)
	}
}
