package domain

// Log is a row of mm_logs. It has no update_at: audit rows are append-only.
type Log struct {
	LogID     uint   `db:"log_id"`
	UUID      string `db:"uuid"`
	UserID    uint   `db:"user_id"`
	Operation string `db:"operation"`
	Version   int    `db:"version"`
	CreateAt  uint   `db:"create_at"`
}
