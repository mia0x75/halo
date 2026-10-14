package domain

// UserStatus is the account state. Stored in mm_users.status.
type UserStatus int

const (
	UserStatusNormal  UserStatus = 1
	UserStatusBlocked UserStatus = 2
	UserStatusPending UserStatus = 3
)

// User is a row of mm_users.
type User struct {
	UserID   uint       `db:"user_id"`
	UUID     string     `db:"uuid"`
	Email    string     `db:"email"`
	Password string     `db:"password"`
	Status   UserStatus `db:"status"`
	Name     string     `db:"name"`
	Phone    uint64     `db:"phone"`
	AvatarID uint       `db:"avatar_id"`
	Version  int        `db:"version"`
	UpdateAt *uint      `db:"update_at"`
	CreateAt uint       `db:"create_at"`
}
