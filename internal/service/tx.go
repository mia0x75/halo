package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// TxManager runs fn inside a single transaction. Every repository in the Repos
// bundle shares that transaction, so a service can compose several writes and
// their audit log atomically.
type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, repos Repos) error) error
}

// Repos bundles all repositories for use inside a transaction. Fields are
// populated by the TxManager implementation.
type Repos struct {
	Avatar    domain.AvatarRepository
	Cluster   domain.ClusterRepository
	Comment   domain.CommentRepository
	Cron      domain.CronRepository
	Edge      domain.EdgeRepository
	Glossary  domain.GlossaryRepository
	Log       domain.LogRepository
	Option    domain.OptionRepository
	Query     domain.QueryRepository
	Role      domain.RoleRepository
	Rule      domain.RuleRepository
	Statement domain.StatementRepository
	Metric    domain.MetricRepository
	RuleHit   domain.RuleHitRepository
	Message   domain.MessageRepository
	Group     domain.GroupRepository
	Task      domain.TaskRepository
	Template  domain.TemplateRepository
	Ticket    domain.TicketRepository
	User      domain.UserRepository
}
