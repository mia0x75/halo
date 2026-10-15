// Package app is the dependency container. It owns the metadata connection
// pool and wires the repositories and services built on top of it. During the
// 2D migration it runs alongside the legacy xorm engine (g.Engine); main.go
// initialises both and this container is additive until 2D-B.11 removes xorm.
package app

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/analyzer"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/db"
	mysqlrepo "github.com/mia0x75/halo/internal/repository/mysql"
	"github.com/mia0x75/halo/internal/service"
	"github.com/mia0x75/halo/validate"
)

// App is the dependency container.
type App struct {
	DB       *db.DB
	Repos    service.Repos
	Tx       service.TxManager
	Agent    agent.Executor
	Services Services
}

// Services bundles the application services.
type Services struct {
	Log       *service.LogService
	Ticket    *service.TicketService
	Statement *service.StatementService
	Review    *service.ReviewService
	Agent     *service.AgentService
	Avatar    *service.AvatarService
	Cluster   *service.ClusterService
	Comment   *service.CommentService
	Cron      *service.CronService
	Edge      *service.EdgeService
	Glossary  *service.GlossaryService
	Option    *service.OptionService
	Query     *service.QueryService
	Role      *service.RoleService
	Rule      *service.RuleService
	Statistic *service.StatisticService
	Template  *service.TemplateService
	User      *service.UserService
}

// New opens the metadata pool and wires repositories and services.
func New(ctx context.Context, cfg db.Config) (*App, error) {
	database, err := db.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("app: %w", err)
	}

	// database.DB is the embedded *sqlx.DB; its own embedded *sql.DB feeds the
	// tx manager, while the *sqlx.DB satisfies the sqlc DBTX interface.
	tx := mysqlrepo.NewTxManager(database.DB.DB)
	q := mysqlrepo.New(database.DB)

	repos := service.Repos{
		Avatar:    mysqlrepo.NewAvatarRepo(q),
		Cluster:   mysqlrepo.NewClusterRepo(q),
		Comment:   mysqlrepo.NewCommentRepo(q),
		Cron:      mysqlrepo.NewCronRepo(q),
		Edge:      mysqlrepo.NewEdgeRepo(q),
		Glossary:  mysqlrepo.NewGlossaryRepo(q),
		Log:       mysqlrepo.NewLogRepo(q),
		Option:    mysqlrepo.NewOptionRepo(q),
		Query:     mysqlrepo.NewQueryRepo(q),
		Role:      mysqlrepo.NewRoleRepo(q),
		Rule:      mysqlrepo.NewRuleRepo(q),
		Statement: mysqlrepo.NewStatementRepo(q),
		Statistic: mysqlrepo.NewStatisticRepo(q),
		Template:  mysqlrepo.NewTemplateRepo(q),
		Ticket:    mysqlrepo.NewTicketRepo(q),
		User:      mysqlrepo.NewUserRepo(q),
	}

	logSvc := service.NewLogService(repos.Log)
	// The schema source is deferred to 2D-B, so the adapter is still a skeleton.
	validator := validate.NewAdapter(mysqlrepo.NewRuleSource(repos.Rule), nil)
	analyzerSvc, err := analyzer.New(domain.ClusterTypeMySQL, validator)
	if err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("app: analyzer: %w", err)
	}

	services := Services{
		Log:       logSvc,
		Ticket:    service.NewTicketService(tx, repos.Ticket),
		Statement: service.NewStatementService(repos.Statement, logSvc, analyzerSvc),
		Review:    service.NewReviewService(tx, repos.Ticket, validator),
		Agent:     service.NewAgentService(tx, agent.LocalExecutor{}),
		Avatar:    service.NewAvatarService(tx, repos.Avatar),
		Cluster:   service.NewClusterService(tx, repos.Cluster),
		Comment:   service.NewCommentService(tx, repos.Comment),
		Cron:      service.NewCronService(tx, repos.Cron),
		Edge:      service.NewEdgeService(tx, repos.Edge),
		Glossary:  service.NewGlossaryService(tx, repos.Glossary),
		Option:    service.NewOptionService(tx, repos.Option),
		Query:     service.NewQueryService(tx, repos.Query),
		Role:      service.NewRoleService(tx, repos.Role),
		Rule:      service.NewRuleService(tx, repos.Rule),
		Statistic: service.NewStatisticService(tx, repos.Statistic),
		Template:  service.NewTemplateService(tx, repos.Template),
		User:      service.NewUserService(tx, repos.User),
	}

	return &App{DB: database, Repos: repos, Tx: tx, Agent: agent.LocalExecutor{}, Services: services}, nil
}

// Close releases the metadata pool.
func (a *App) Close() error {
	if a.DB == nil {
		return nil
	}
	return a.DB.Close()
}
