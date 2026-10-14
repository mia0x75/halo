package domain

// Role is a row of mm_roles. Role semantics come from the seed name, not a
// status column.
type Role struct {
	RoleID      uint   `db:"role_id"`
	UUID        string `db:"uuid"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Version     int    `db:"version"`
	UpdateAt    *uint  `db:"update_at"`
	CreateAt    uint   `db:"create_at"`
}
