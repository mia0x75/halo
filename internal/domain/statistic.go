package domain

// Statistic is a row of mm_statistics. The primary key is (group, key); UUID is
// the stable lookup key used by the repository.
type Statistic struct {
	Group    string  `db:"group"`
	Key      string  `db:"key"`
	UUID     string  `db:"uuid"`
	Value    float64 `db:"value"`
	Version  int     `db:"version"`
	UpdateAt *uint   `db:"update_at"`
	CreateAt uint    `db:"create_at"`
}
