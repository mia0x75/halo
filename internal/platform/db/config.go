package db

// Config holds metadata connection pool settings. Field names mirror
// cfg.Database so callers can map it without a conversion shim.
type Config struct {
	Addr           string
	MaxIdle        int
	MaxConnections int
	WaitTimeout    int
}
