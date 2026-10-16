package domain

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
