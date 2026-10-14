package domain

// Glossary is a row of mm_glossaries. The primary key is (group, key); UUID is
// the stable lookup key used by the repository.
type Glossary struct {
	Group       string `db:"group"`
	Key         uint8  `db:"key"`
	Value       string `db:"value"`
	UUID        string `db:"uuid"`
	Description string `db:"description"`
	Version     int    `db:"version"`
	UpdateAt    *uint  `db:"update_at"`
	CreateAt    uint   `db:"create_at"`
}
