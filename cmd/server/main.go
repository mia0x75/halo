//go:build !windows
// +build !windows

package main

import (
	"context"
	"errors"
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
	"github.com/mia0x75/halo/internal/transport/server"
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

	var application *app.App

	fmt.Print(app.Banner)
	fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
		"Version", app.Version,
		"Git commit", app.Git,
		"BuildAt", app.BuildAt,
		"Go version", app.GoVer,
		"Distro", app.Distro,
		"Kernel", app.Kernel,
		"Branch", app.Branch,
	)
	fmt.Println()
	if *version {
		os.Exit(0)
	}

	config.ParseConfig(*cfg)
	// 子命令（halocli migrate 等）用同一个配置文件；还没写盘时指向将要写入的位置。
	if config.SavePath != "" {
		os.Setenv("HALO_CFG", config.SavePath)
	} else {
		os.Setenv("HALO_CFG", config.ConfigFile)
	}
	configureAuth()

	// 未配置元数据库与密钥时，服务器先跑在初始化模式：只提供 /welcome 向导
	// 与前端静态资源，向导完成后热切换到正常路由。
	handler := server.NewSwitchable(nil)
	defer func() {
		if application != nil {
			_ = application.Close()
		}
	}()

	if config.Configured() {
		built, err := bootApp(context.Background())
		if err != nil {
			lg.Fatalf("[F] init app fail: %v", err)
		}
		application = built
		handler.Set(normalRouter(application))
	} else {
		lg.Warnf("[W] metadata database is not configured; open /welcome to finish setup")
		setup := server.NewSetupHandler(func(built *app.App) error {
			application = built
			configureAuth()
			handler.Set(normalRouter(built))
			return nil
		})
		r := mux.NewRouter()
		server.SetupRoutes(r, setup)
		handler.Set(r)
	}

	addr := config.Config().Listen
	certFile := config.Config().Cert
	keyFile := config.Config().Key
	// TLS 是可选项：证书缺失（没有 certs 目录或目录为空）时直接跑 HTTP，
	// 让内网首次部署不必先自签证书才能打开 /welcome。
	useTLS := fileExists(certFile) && fileExists(keyFile)
	if useTLS {
		lg.Infof("[I] https listening %s", addr)
	} else {
		lg.Warnf("[W] TLS certificate not found (%s / %s); serving plain http on %s", certFile, keyFile, addr)
	}

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15, // Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:  time.Second * 15, //
		IdleTimeout:  time.Second * 60, //
		Handler:      handler,          // 初始化完成后热切换到正常路由。
	}
	srv.SetKeepAlivesEnabled(true)
	// 各种中间件
	// csrf + gzip
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		var err error
		if useTLS {
			err = srv.ListenAndServeTLS(certFile, keyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != http.ErrServerClosed {
			lg.Fatalf("code: 1500, message: %s", err.Error())
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

// bootApp 装配应用：连接元数据库、加载缓存、校验规则并启动后台作业。
func bootApp(ctx context.Context) (*app.App, error) {
	cfg := config.Config()
	if cfg.Database == nil {
		return nil, errors.New("缺少元数据库配置（database.addr）")
	}
	application, err := app.New(ctx, db.Config{
		Addr:           cfg.Database.Addr,
		MaxIdle:        cfg.Database.MaxIdle,
		MaxConnections: cfg.Database.MaxConnections,
		WaitTimeout:    cfg.Database.WaitTimeout,
		Driver:         cfg.Database.Driver,
	})
	if err != nil {
		return nil, err
	}

	// 角色名承载权限，角色表又不可增删。名字对不上意味着权限会静默失效，
	// 所以启动时就拒绝，而不是等接口 403。
	if err := transportauth.ValidateRoles(cache.RolesMap.All()); err != nil {
		_ = application.Close()
		return nil, err
	}

	cachedRules := cache.RulesMap.All()
	rules := make([]*domain.Rule, 0, len(cachedRules))
	for _, m := range cachedRules {
		rules = append(rules, bridge.RuleToDomain(m))
	}
	if err := validate.ValidateRules(rules); err != nil {
		_ = application.Close()
		return nil, err
	}
	// Surface rules whose handler is missing; these were silently skipped by
	// the former reflection dispatch.
	validate.ValidateRegistry(rules)

	worker.Init(application.Repos)
	worker.NewScheduler()
	return application, nil
}

// normalRouter builds the route table of a configured server.
func normalRouter(application *app.App) http.Handler {
	r := mux.NewRouter()
	server.Routes(r, application)
	if cfg := config.Config(); cfg != nil && cfg.Log != nil && strings.EqualFold(cfg.Log.Level, "debug") {
		server.ProfilerRoutes(r)
	}
	return r
}

// fileExists reports whether path is a regular file.
func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// configureAuth applies the JWT settings to the auth package.
func configureAuth() {
	secret := config.Config().Secret
	if secret == nil || secret.Jwt == nil {
		return
	}
	ttl := 24 * time.Hour
	if parsed, err := time.ParseDuration(secret.Jwt.TTL); err == nil && parsed > 0 {
		ttl = parsed
	}
	transportauth.Configure(secret.Jwt.Issuer, secret.Jwt.SecretKey, ttl)
}
