package domain

// Comment is a row of mm_comments.
type Comment struct {
	CommentID uint   `db:"comment_id"`
	UUID      string `db:"uuid"`
	Content   string `db:"content"`
	TicketID  uint   `db:"ticket_id"`
	UserID    uint   `db:"user_id"`
	Version   int    `db:"version"`
	UpdateAt  *uint  `db:"update_at"`
	CreateAt  uint   `db:"create_at"`
}
