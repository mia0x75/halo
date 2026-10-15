package routers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/mia0x75/halo/directives"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/resolvers"
)

var cfg = gqlapi.Config{}

func init() {
	cfg = gqlapi.Config{
		Directives: gqlapi.DirectiveRoot{
			Auth:    directives.Auth,
			Date:    directives.Date,
			EnumInt: directives.EnumInt,
			Length:  directives.Length,
			Lower:   directives.Lower,
			Matches: directives.Matches,
			Range:   directives.Range,
			Rename:  directives.Rename,
			Trim:    directives.Trim,
			Upper:   directives.Upper,
			Uuid:    directives.Uuid,
		},
	}
	// countComplexity := func(childComplexity, count int) int {
	// 	return count * childComplexity
	// }
	// c.Complexity.User.Reviewers = countComplexity
	// cfg.Complexity.User.Reviewers = func(childComplexity int) int {
	// 	return childComplexity
	// }
}

// 示例代码 - 开始
// Define our struct
type AuthMiddleware struct {
	tokenUsers map[string]string
}

// Initialize it somewhere
func (m *AuthMiddleware) Populate() {
	m.tokenUsers["00000000"] = "user0"
	m.tokenUsers["aaaaaaaa"] = "userA"
	m.tokenUsers["05f717e5"] = "randomUser"
	m.tokenUsers["deadbeef"] = "user0"
}

// Middleware function, which will be called for each request
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := m.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

// 示例代码 - 结束

func Routes(r *mux.Router, application *app.App) {
	cfg.Resolvers = resolvers.NewResolver(application)
	// TODO: 暂时保留
	// CSRF := csrf.Protect(
	// 	[]byte("a-32-byte-long-key-goes-here"),
	// 	csrf.RequestHeader("Authenticity-Token"),
	// 	csrf.FieldName("authenticity_token"),
	// 	csrf.ErrorHandler(http.HandlerFunc(Error403)),
	// )
	r.Use(Headers)
	r.Use(Context)
	r.Use(Cors)
	r.Use(Logging)
	r.Use(Compress)

	r.HandleFunc("/", playground.Handler("GraphQL playground", "/api/query"))

	srv := handler.New(gqlapi.NewExecutableSchema(cfg))

	// WebSocket 传输（替代旧版 handler.WebsocketUpgrader，默认允许所有来源的等价行为由 coder/websocket 实现承担；
	// 如需自定义来源校验，可通过 transport.CoderWebsocketImplementation{AcceptOptions} 传入）。
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		// notify bug tracker...
		return fmt.Errorf("错误代码: %s, 错误信息: 内部服务器错误，%v", gqlapi.ReturnCodeUnknowError, err)
	})

	srv.Use(extension.Introspection{})

	/* RequestMiddleware signature error ??
	if strings.EqualFold(config.Config().Log.Level, "debug") {
		// 性能跟踪，DEBUG模式时打开
		// 自省功能，DEBUG模式时打开
		// TODO: 处理复杂度
		// 可按需在此添加 extension.FixedComplexityLimit(1000) 等扩展
	}
	*/

	r.HandleFunc("/api/query", srv.ServeHTTP)
}

func Error403(w http.ResponseWriter, r *http.Request) {
	return
}

func ProfilerRoutes(r *mux.Router) {
	p := r.Path("/debug").Subrouter()
	p.HandleFunc("/pprof/", pprof.Index)
	p.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	p.HandleFunc("/pprof/profile", pprof.Profile)
	p.HandleFunc("/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	p.HandleFunc("/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	p.HandleFunc("/pprof/heap", pprof.Handler("heap").ServeHTTP)
	p.HandleFunc("/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	p.HandleFunc("/pprof/block", pprof.Handler("block").ServeHTTP)
}
