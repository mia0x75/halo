package domain

// CronStatus is the run state of a task. Stored in mm_crons.status.
// Legacy rows also carry H/C/D/E/R; they are not enumerated until the
// cron engine is ported to this layer.
type CronStatus string

const (
	CronStatusPending CronStatus = "P"
	CronStatusSuccess CronStatus = "S"
	CronStatusFailure CronStatus = "F"
)

// Cron is a row of mm_crons.
type Cron struct {
	CronID    uint       `db:"cron_id"`
	UUID      string     `db:"uuid"`
	Status    CronStatus `db:"status"`
	Name      string     `db:"name"`
	Cmd       string     `db:"cmd"`
	Params    string     `db:"params"`
	LastRun   *string    `db:"last_run"`
	NextRun   *string    `db:"next_run"`
	Interval  *string    `db:"interval"`
	Duration  *string    `db:"duration"`
	Recurrent *uint8     `db:"recurrent"`
	Hash      *string    `db:"hash"`
	UserID    uint       `db:"user_id"`
	ClusterID uint       `db:"cluster_id"`
	Version   int        `db:"version"`
	UpdateAt  *uint      `db:"update_at"`
	CreateAt  uint       `db:"create_at"`
}
