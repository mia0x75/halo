package analyzer

import "errors"

// ErrParseFailed wraps backend parse failures so callers never depend on a
// concrete parser's error type.
var ErrParseFailed = errors.New("analyzer: parse failed")

// ErrUnsupportedStmt means the SQL parsed but cannot be lowered to IR.
var ErrUnsupportedStmt = errors.New("analyzer: unsupported statement")

// ErrNotImplemented marks backend paths that are wired but not built yet.
var ErrNotImplemented = errors.New("analyzer: not implemented")
