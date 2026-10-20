package mysqlrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mia0x75/halo/internal/service"
)

// TxManager implements service.TxManager on the metadata MySQL pool. It lives
// here because the repository layer owns the driver; the interface lives in
// service because that is where the orchestration demand originates.
type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

var _ service.TxManager = (*TxManager)(nil)

// WithTx begins a transaction, binds every repository to it and commits only
// when fn returns nil. Any error (including a rollback path) discards the lot.
func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context, repos service.Repos) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	q := New(tx)
	repos := service.Repos{
		Avatar:    NewAvatarRepo(q),
		Cluster:   NewClusterRepo(q),
		Comment:   NewCommentRepo(q),
		Cron:      NewCronRepo(q),
		Edge:      NewEdgeRepo(q),
		Glossary:  NewGlossaryRepo(q),
		Log:       NewLogRepo(q),
		Option:    NewOptionRepo(q),
		Query:     NewQueryRepo(q),
		Role:      NewRoleRepo(q),
		Rule:      NewRuleRepo(q),
		Statement: NewStatementRepo(q),
		Metric:    NewMetricRepo(q),
		RuleHit:   NewRuleHitRepo(q),
		Template:  NewTemplateRepo(q),
		Ticket:    NewTicketRepo(q),
		User:      NewUserRepo(q),
	}

	if err := fn(ctx, repos); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
