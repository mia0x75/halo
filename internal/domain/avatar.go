package domain

// Avatar is a row of mm_avatars. The mm_ prefix is stripped in code (ADR 0002).
type Avatar struct {
	AvatarID uint   `db:"avatar_id"`
	UUID     string `db:"uuid"`
	URL      string `db:"url"`
	Version  int    `db:"version"`
	UpdateAt *uint  `db:"update_at"`
	CreateAt uint   `db:"create_at"`
}
