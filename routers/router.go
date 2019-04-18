package routers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/mia0x75/halo/directives"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/resolvers"
)

var cfg = gqlapi.Config{}

func init() {
	cfg = gqlapi.Config{
		Resolvers: &resolvers.Resolver{},
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

func Routes(r *mux.Router) {
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

	r.HandleFunc("/", handler.Playground("GraphQL playground", "/api/query"))

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}

	options := []handler.Option{
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			// notify bug tracker...
			return fmt.Errorf("错误代码: %s, 错误信息: 内部服务器错误，%v", gqlapi.ReturnCodeUnknowError, err)
		}),
		handler.WebsocketUpgrader(upgrader),
	}
	if strings.EqualFold(g.Config().Log.Level, "debug") {
		options = append(options,
			// 性能跟踪，DEBUG模式时打开
			// interface conversion: interface {} is nil, not *gqlapollotracing.tracingData
			// handler.RequestMiddleware(gqlapollotracing.RequestMiddleware()),
			// handler.Tracer(gqlapollotracing.NewTracer()),
			// 自省功能，DEBUG模式时打开
			handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
				// if userForContext(ctx).IsAdmin {
				// 	graphql.GetRequestContext(ctx).DisableIntrospection = true
				// }

				graphql.GetRequestContext(ctx).DisableIntrospection = false
				return next(ctx)
			}),
			// TODO: 处理复杂度
			handler.ComplexityLimit(1000),
		)
	}

	r.HandleFunc("/api/query", handler.GraphQL(
		gqlapi.NewExecutableSchema(cfg),
		options...,
	))
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
