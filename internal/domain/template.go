package domain

// Template is a row of mm_templates. The primary key is uuid.
type Template struct {
	UUID        string `db:"uuid"`
	Subject     string `db:"subject"`
	Body        string `db:"body"`
	Description string `db:"description"`
	Version     int    `db:"version"`
	UpdateAt    *uint  `db:"update_at"`
	CreateAt    uint   `db:"create_at"`
}
