package domain

import "errors"

// Sentinel errors let callers map repository failures to API responses without
// depending on a concrete driver.
var (
	ErrNotFound     = errors.New("domain: not found")
	ErrInvalidInput = errors.New("domain: invalid input")
	ErrConflict     = errors.New("domain: conflict")
	ErrForbidden    = errors.New("domain: forbidden")
)
