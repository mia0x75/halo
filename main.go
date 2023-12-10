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
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/crons"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/routers"
)

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

	fmt.Println(g.Banner)
	fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n",
		"Version", g.Version,
		"Git commit", g.Git,
		"Compile", g.Compile,
		"Distro", g.Distro,
		"Kernel", g.Kernel,
		"Branch", g.Branch,
	)
	fmt.Println()
	if *version {
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	os.Setenv("HALO_CFG", g.ConfigFile)

	g.InitLog()
	if err := g.InitDB(); err != nil {
		os.Exit(0)
	}
	caches.Init()
	crons.NewScheduler()

	addr := g.Config().Listen
	log.Infof("[I] http listening %s", addr)

	r := mux.NewRouter()
	routers.Routes(r)
	if strings.EqualFold(g.Config().Log.Level, "debug") {
		routers.ProfilerRoutes(r)
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
		if err := srv.ListenAndServeTLS(g.Config().Cert, g.Config().Key); err != http.ErrServerClosed {
			log.Fatalf("错误代码: 1500, 错误信息: %s", err.Error())
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
	log.Println("[I] Shutting down")
	os.Exit(0)
}
