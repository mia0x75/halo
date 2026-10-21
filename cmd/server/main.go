//go:build !windows
// +build !windows

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/db"
	"github.com/mia0x75/halo/internal/validate"
	"github.com/mia0x75/halo/internal/worker"

	// "github.com/mia0x75/halo/internal/platform/log"
	transportauth "github.com/mia0x75/halo/internal/transport/auth"
	"github.com/mia0x75/halo/internal/transport/http"
)

var lg = logrus.New()

const ticketLoaderKey = "ticketloader"

// func DataloaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ticketloader := TicketLoader{
// 			maxBatch: 50,
// 			wait:     1 * time.Millisecond,
// 			fetch: func(ids []int) ([]*models.Ticket, []error) {
// 				tickets := []*models.Ticket{}
// 				g.Engine.In("ticket_id", ids).Omit("content").Find(tickets)
// 				return tickets, nil
// 			},
// 		}
// 		ctx := context.WithValue(r.Context(), ticketLoaderKey, &ticketloader)
// 		r = r.WithContext(ctx)
// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	cfg := flag.String("c", "", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	fmt.Print(app.Banner)
	fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
		"Version", app.Version,
		"Git commit", app.Git,
		"Compile", app.Compile,
		"Distro", app.Distro,
		"Kernel", app.Kernel,
		"Branch", app.Branch,
	)
	fmt.Println()
	if *version {
		os.Exit(0)
	}

	config.ParseConfig(*cfg)
	os.Setenv("HALO_CFG", config.ConfigFile)

	if secret := config.Config().Secret; secret != nil && secret.Jwt != nil {
		ttl := 24 * time.Hour
		if parsed, err := time.ParseDuration(secret.Jwt.TTL); err == nil && parsed > 0 {
			ttl = parsed
		}
		transportauth.Configure(secret.Jwt.Issuer, secret.Jwt.SecretKey, ttl)
	}

	// log.InitLog()

	// The new dependency container replaces the legacy ORM engine.

	application, err := app.New(context.Background(), db.Config{
		Addr:           config.Config().Database.Addr,
		MaxIdle:        config.Config().Database.MaxIdle,
		MaxConnections: config.Config().Database.MaxConnections,
		WaitTimeout:    config.Config().Database.WaitTimeout,
		Driver:         config.Config().Database.Driver,
	})
	if err != nil {
		lg.Fatalf("[F] init app fail: %v", err)
	}
	defer func() { _ = application.Close() }()

	cachedRules := cache.RulesMap.All()
	rules := make([]*domain.Rule, 0, len(cachedRules))
	for _, m := range cachedRules {
		rules = append(rules, bridge.RuleToDomain(m))
	}
	if err := validate.ValidateRules(rules); err != nil {
		lg.Fatalf("[F] invalid validation rules: %v", err)
	}
	// Surface rules whose handler is missing; these were silently skipped by
	// the former reflection dispatch.
	validate.ValidateRegistry(rules)
	worker.Init(application.Repos)
	worker.NewScheduler()

	addr := config.Config().Listen
	lg.Infof("[I] http listening %s", addr)

	r := mux.NewRouter()
	server.Routes(r, application)
	if strings.EqualFold(config.Config().Log.Level, "debug") {
		server.ProfilerRoutes(r)
	}
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15, // Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:  time.Second * 15, //
		IdleTimeout:  time.Second * 60, //
		Handler:      r,                // Pass our instance of gorilla/mux in.
	}
	srv.SetKeepAlivesEnabled(true)
	// 各种中间件
	// csrf + gzip
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServeTLS(config.Config().Cert, config.Config().Key); err != http.ErrServerClosed {
			lg.Fatalf("错误代码: 1500, 错误信息: %s", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigs

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	lg.Println("[I] Shutting down")
	os.Exit(0)
}
