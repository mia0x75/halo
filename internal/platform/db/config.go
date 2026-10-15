package db

// Config holds metadata connection pool settings. Field names mirror
// cfg.Database so callers can map it without a conversion shim.
type Config struct {
	Addr           string
	MaxIdle        int
	MaxConnections int
	WaitTimeout    int
	// Driver selects the metadata engine: "mysql" or "postgres".
	// Empty means auto-detect from Addr. It mirrors cfg.Database.Driver.
	Driver string
}
