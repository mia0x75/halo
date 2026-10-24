package domain

import "context"

// Repository contracts. UUID is the stable lookup/delete key for every entity
// (ADR 0002 keeps the mm_ prefix out of these names). Entities with an
// auto-increment integer key additionally expose the int64 variants.
//
// List uses keyset pagination: cursor is the last seen key (int64 for entities
// with an integer key, string UUID otherwise), limit caps the page size.

type AvatarRepository interface {
	Create(ctx context.Context, a *Avatar) (int64, error)
	GetByID(ctx context.Context, id int64) (*Avatar, error)
	GetByUUID(ctx context.Context, uuid string) (*Avatar, error)
	Update(ctx context.Context, a *Avatar) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Avatar, error)
}

type CommentRepository interface {
	Create(ctx context.Context, c *Comment) (int64, error)
	GetByID(ctx context.Context, id int64) (*Comment, error)
	GetByUUID(ctx context.Context, uuid string) (*Comment, error)
	Update(ctx context.Context, c *Comment) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Comment, error)
	ListByTicket(ctx context.Context, ticketID uint) ([]*Comment, error)
}

type CronRepository interface {
	Create(ctx context.Context, c *Cron) (int64, error)
	GetByID(ctx context.Context, id int64) (*Cron, error)
	GetByUUID(ctx context.Context, uuid string) (*Cron, error)
	Update(ctx context.Context, c *Cron) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Cron, error)
	ListSchedulable(ctx context.Context) ([]*Cron, error)
	// Count returns the total number of rows, independent of any page cursor.
	Count(ctx context.Context) (int64, error)
}

type ClusterRepository interface {
	Create(ctx context.Context, c *Cluster) (int64, error)
	GetByID(ctx context.Context, id int64) (*Cluster, error)
	GetByUUID(ctx context.Context, uuid string) (*Cluster, error)
	Update(ctx context.Context, c *Cluster) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Cluster, error)
	// CountByGroup counts the clusters still assigned to a group.
	CountByGroup(ctx context.Context, groupID uint) (int64, error)
}

type LogRepository interface {
	Create(ctx context.Context, l *Log) (int64, error)
	GetByID(ctx context.Context, id int64) (*Log, error)
	GetByUUID(ctx context.Context, uuid string) (*Log, error)
	Update(ctx context.Context, l *Log) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Log, error)
	// DeleteBefore removes the logs created before the given time.
	DeleteBefore(ctx context.Context, before uint) (int64, error)
}

type QueryRepository interface {
	Create(ctx context.Context, q *Query) (int64, error)
	GetByID(ctx context.Context, id int64) (*Query, error)
	GetByUUID(ctx context.Context, uuid string) (*Query, error)
	Update(ctx context.Context, q *Query) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Query, error)
	// Count returns the total number of rows, independent of any page cursor.
	Count(ctx context.Context) (int64, error)
	ListByUser(ctx context.Context, userID uint) ([]*Query, error)
}

type EdgeRepository interface {
	Create(ctx context.Context, e *Edge) (int64, error)
	GetByID(ctx context.Context, id int64) (*Edge, error)
	GetByUUID(ctx context.Context, uuid string) (*Edge, error)
	Update(ctx context.Context, e *Edge) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	DeleteByTypeAncestor(ctx context.Context, edgeType uint, ancestorID uint) error
	List(ctx context.Context, cursor int64, limit int) ([]*Edge, error)
}

type RoleRepository interface {
	Create(ctx context.Context, r *Role) (int64, error)
	GetByID(ctx context.Context, id int64) (*Role, error)
	GetByUUID(ctx context.Context, uuid string) (*Role, error)
	Update(ctx context.Context, r *Role) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Role, error)
}

type TicketRepository interface {
	Create(ctx context.Context, t *Ticket) (int64, error)
	GetByID(ctx context.Context, id int64) (*Ticket, error)
	GetByUUID(ctx context.Context, uuid string) (*Ticket, error)
	GetByCronID(ctx context.Context, cronID uint) (*Ticket, error)
	Update(ctx context.Context, t *Ticket) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*Ticket, error)
	// ListFiltered lists the tickets matching the filters. A zero status means
	// every status, an empty keyword matches every subject and database.
	ListFiltered(ctx context.Context, cursor int64, limit int, status uint8, keyword string) ([]*Ticket, error)
	// CountFiltered counts the tickets matching the same filters as
	// ListFiltered, independent of any page cursor.
	CountFiltered(ctx context.Context, status uint8, keyword string) (int64, error)
	ListByUser(ctx context.Context, userID uint) ([]*Ticket, error)
	// ListByReviewer returns every ticket the user reviews, oldest first.
	ListByReviewer(ctx context.Context, userID uint) ([]*Ticket, error)
	// CountByStatus counts the tickets per status inside a time window.
	CountByStatus(ctx context.Context, from, to uint) ([]TicketStatusCount, error)
}

type UserRepository interface {
	Create(ctx context.Context, u *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUUID(ctx context.Context, uuid string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor int64, limit int) ([]*User, error)
	// Count returns the total number of rows, independent of any page cursor.
	Count(ctx context.Context) (int64, error)
}

type GlossaryRepository interface {
	Create(ctx context.Context, g *Glossary) error
	GetByUUID(ctx context.Context, uuid string) (*Glossary, error)
	GetByKey(ctx context.Context, group string, key uint8) (*Glossary, error)
	Update(ctx context.Context, g *Glossary) error
	DeleteByUUID(ctx context.Context, uuid string) error
	DeleteByKey(ctx context.Context, group string, key uint8) error
	List(ctx context.Context, cursor string, limit int) ([]*Glossary, error)
}

type MetricRepository interface {
	Create(ctx context.Context, m *Metric) error
	GetByUUID(ctx context.Context, uuid string) (*Metric, error)
	GetByKey(ctx context.Context, key MetricKey) (*Metric, error)
	Update(ctx context.Context, m *Metric) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor string, limit int) ([]*Metric, error)
	// Increment atomically adds delta to a counter metric, creating the row
	// when it does not exist yet.
	Increment(ctx context.Context, key MetricKey, delta float64) error
	// SetGauge overwrites the value of a gauge metric.
	SetGauge(ctx context.Context, key MetricKey, value float64) error
	// DeleteBucketsBefore removes the daily buckets older than the given date.
	// The bucket is a YYYY-MM-DD string, so it compares lexicographically.
	DeleteBucketsBefore(ctx context.Context, bucket string) (int64, error)
}

type RuleHitRepository interface {
	Create(ctx context.Context, hit *RuleHit) error
	CreateBatch(ctx context.Context, hits []*RuleHit) error
	ListByTicket(ctx context.Context, ticketID uint) ([]*RuleHit, error)
	// List returns the most recent hits, optionally of one ticket.
	List(ctx context.Context, ticketID uint, limit int) ([]*RuleHit, error)
	// CountByLevel counts the hits per severity inside a time window.
	CountByLevel(ctx context.Context, from, to uint) ([]RuleHitCount, error)
	// CountByRule counts the hits per rule inside a time window.
	CountByRule(ctx context.Context, from, to uint, limit int) ([]RuleHitCount, error)
	// DeleteBefore removes the hits created before the given time.
	DeleteBefore(ctx context.Context, before uint) (int64, error)
	// DeleteByStatement removes the hits of one statement, so a revalidation
	// replaces the previous hits instead of adding a second copy.
	DeleteByStatement(ctx context.Context, ticketID uint, sequence uint16) (int64, error)
	// DeleteByTicket removes every hit of one ticket.
	DeleteByTicket(ctx context.Context, ticketID uint) (int64, error)
}

type StatementRepository interface {
	Create(ctx context.Context, s *Statement) error
	GetByUUID(ctx context.Context, uuid string) (*Statement, error)
	GetByTicketSequence(ctx context.Context, ticketID uint, sequence uint16) (*Statement, error)
	Update(ctx context.Context, s *Statement) error
	DeleteByUUID(ctx context.Context, uuid string) error
	DeleteByTicketSequence(ctx context.Context, ticketID uint, sequence uint16) error
	DeleteByTicket(ctx context.Context, ticketID uint) error
	List(ctx context.Context, cursor string, limit int) ([]*Statement, error)
}

type OptionRepository interface {
	Create(ctx context.Context, o *Option) error
	GetByUUID(ctx context.Context, uuid string) (*Option, error)
	GetByName(ctx context.Context, name string) (*Option, error)
	Update(ctx context.Context, o *Option) error
	DeleteByUUID(ctx context.Context, uuid string) error
	// List returns the options ordered by group and by the display order.
	List(ctx context.Context, limit int) ([]*Option, error)
}

type RuleRepository interface {
	Create(ctx context.Context, r *Rule) error
	GetByUUID(ctx context.Context, uuid string) (*Rule, error)
	GetByName(ctx context.Context, name string) (*Rule, error)
	Update(ctx context.Context, r *Rule) error
	DeleteByUUID(ctx context.Context, uuid string) error
	// List returns rules ordered by rule ID (name) ascending. The cursor is the
	// rule name to continue from; an empty cursor starts at the first rule.
	List(ctx context.Context, cursor string, limit int) ([]*Rule, error)
}

type TemplateRepository interface {
	Create(ctx context.Context, t *Template) error
	GetByUUID(ctx context.Context, uuid string) (*Template, error)
	Update(ctx context.Context, t *Template) error
	DeleteByUUID(ctx context.Context, uuid string) error
	List(ctx context.Context, cursor string, limit int) ([]*Template, error)
}
