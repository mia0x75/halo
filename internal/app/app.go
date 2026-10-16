// Package app is the dependency container. It owns the metadata connection
// pool and wires the repositories and services built on top of it. During the
// 2D migration it runs alongside the legacy engine; main.go
// initialises both and this container is additive until 2D cleanup completes.

package app

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/analyzer"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/event"
	"github.com/mia0x75/halo/internal/platform/db"
	"github.com/mia0x75/halo/internal/platform/notify"
	mysqlrepo "github.com/mia0x75/halo/internal/repository/mysql"
	postgresrepo "github.com/mia0x75/halo/internal/repository/postgres"
	"github.com/mia0x75/halo/internal/service"
	"github.com/mia0x75/halo/internal/validate"
)

// App is the dependency container.
type App struct {
	DB       *db.DB
	Repos    service.Repos
	Tx       service.TxManager
	Agent    agent.Executor
	EventBus *event.Bus
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
	// The repository package follows the configured metadata engine.
	var tx service.TxManager
	var repos service.Repos
	var ruleSource validate.RuleSource
	if database.Driver == "postgres" {
		ptx := postgresrepo.NewTxManager(database.DB.DB)
		q := postgresrepo.New(database.DB)
		tx = ptx
		repos = service.Repos{
			Avatar:    postgresrepo.NewAvatarRepo(q),
			Cluster:   postgresrepo.NewClusterRepo(q),
			Comment:   postgresrepo.NewCommentRepo(q),
			Cron:      postgresrepo.NewCronRepo(q),
			Edge:      postgresrepo.NewEdgeRepo(q),
			Glossary:  postgresrepo.NewGlossaryRepo(q),
			Log:       postgresrepo.NewLogRepo(q),
			Option:    postgresrepo.NewOptionRepo(q),
			Query:     postgresrepo.NewQueryRepo(q),
			Role:      postgresrepo.NewRoleRepo(q),
			Rule:      postgresrepo.NewRuleRepo(q),
			Statement: postgresrepo.NewStatementRepo(q),
			Statistic: postgresrepo.NewStatisticRepo(q),
			Template:  postgresrepo.NewTemplateRepo(q),
			Ticket:    postgresrepo.NewTicketRepo(q),
			User:      postgresrepo.NewUserRepo(q),
		}
		ruleSource = postgresrepo.NewRuleSource(repos.Rule)
	} else {
		mtx := mysqlrepo.NewTxManager(database.DB.DB)
		q := mysqlrepo.New(database.DB)
		tx = mtx
		repos = service.Repos{
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
		ruleSource = mysqlrepo.NewRuleSource(repos.Rule)
	}

	logSvc := service.NewLogService(repos.Log)
	// The schema source is deferred to 2D-B, so the adapter is still a skeleton.
	validator := validate.NewAdapter(ruleSource, nil)
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

	a := &App{DB: database, Repos: repos, Tx: tx, Agent: agent.LocalExecutor{}, Services: services}

	// The event bus needs the repositories and the notifier. Wiring it here
	// keeps internal/event free of any service/transport dependency and removes
	// the former package-level init().
	a.EventBus = event.NewBus(100)
	event.Register(a.EventBus, event.Deps{
		Log:       repos.Log,
		Statistic: repos.Statistic,
		User:      repos.User,
		Notifier:  notify.NewMulti(notify.EmailNotifier{}),
	})

	// The caches read through the repositories, so they can only be populated
	// once the container is fully wired.
	cache.Init(a.Repos)

	return a, nil
}

// Close releases the metadata pool.
func (a *App) Close() error {
	if a.DB == nil {
		return nil
	}
	return a.DB.Close()
}
