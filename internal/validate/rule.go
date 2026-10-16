package validate

import (
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

// ruleFromDomain wraps a domain entity with an empty parse cache.
func ruleFromDomain(d *domain.Rule) *Rule {
	if d == nil {
		return nil
	}
	return &Rule{Rule: *d}
}
