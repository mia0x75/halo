package domain

// TicketStatus is the workflow state of a ticket. Stored in mm_tickets.status.
type TicketStatus int

const (
	TicketStatusWaitingForVld TicketStatus = 1
	TicketStatusVldFailure    TicketStatus = 2
	TicketStatusVldWarning    TicketStatus = 3
	TicketStatusWaitingForMrv TicketStatus = 4
	TicketStatusMrvFailure    TicketStatus = 5
	TicketStatusLgtm          TicketStatus = 6
	TicketStatusDone          TicketStatus = 7
	TicketStatusExecFailure   TicketStatus = 8
	TicketStatusClosed        TicketStatus = 9
)

// Ticket is a row of mm_tickets. CronID is nullable: manual tickets have none.
type Ticket struct {
	TicketID   uint         `db:"ticket_id"`
	UUID       string       `db:"uuid"`
	ClusterID  uint         `db:"cluster_id"`
	Database   string       `db:"database"`
	Subject    string       `db:"subject"`
	Content    string       `db:"content"`
	Status     TicketStatus `db:"status"`
	UserID     uint         `db:"user_id"`
	ReviewerID uint         `db:"reviewer_id"`
	CronID     *uint        `db:"cron_id"`
	Version    int          `db:"version"`
	UpdateAt   *uint        `db:"update_at"`
	CreateAt   uint         `db:"create_at"`
}

// TicketStatusCount counts the tickets of one status inside a time window.
type TicketStatusCount struct {
	Status TicketStatus
	Count  uint
}
