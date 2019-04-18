package routers

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/tools"
)

const (
	HeaderAccept                          = "Accept"
	HeaderAcceptEncoding                  = "Accept-Encoding"
	HeaderAllow                           = "Allow"
	HeaderAuthorization                   = "Authorization"
	HeaderContentDisposition              = "Content-Disposition"
	HeaderContentEncoding                 = "Content-Encoding"
	HeaderContentLength                   = "Content-Length"
	HeaderContentType                     = "Content-Type"
	HeaderCookie                          = "Cookie"
	HeaderSetCookie                       = "Set-Cookie"
	HeaderIfModifiedSince                 = "If-Modified-Since"
	HeaderLastModified                    = "Last-Modified"
	HeaderLocation                        = "Location"
	HeaderUpgrade                         = "Upgrade"
	HeaderVary                            = "Vary"
	HeaderWWWAuthenticate                 = "WWW-Authenticate"
	HeaderXForwardedFor                   = "X-Forwarded-For"
	HeaderXForwardedProto                 = "X-Forwarded-Proto"
	HeaderXForwardedProtocol              = "X-Forwarded-Protocol"
	HeaderXForwardedSsl                   = "X-Forwarded-Ssl"
	HeaderXUrlScheme                      = "X-Url-Scheme"
	HeaderXHTTPMethodOverride             = "X-HTTP-Method-Override"
	HeaderXRealIP                         = "X-Real-IP"
	HeaderXRequestID                      = "X-Request-ID"
	HeaderXRequestedWith                  = "X-Requested-With"
	HeaderServer                          = "Server"
	HeaderOrigin                          = "Origin"
	HeaderAccessControlRequestMethod      = "Access-Control-Request-Method" // Access control
	HeaderAccessControlRequestHeaders     = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin        = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods       = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders       = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials   = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders      = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge             = "Access-Control-Max-Age"
	HeaderStrictTransportSecurity         = "Strict-Transport-Security" // Security
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
)

type CompressResponseWriter struct {
	io.Writer
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
}

func (w *CompressResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *CompressResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent { // Issue #489
		w.ResponseWriter.Header().Del(HeaderContentEncoding)
	}
	w.Header().Del(HeaderContentLength) // Issue #444
	w.ResponseWriter.WriteHeader(code)
}

func (w *CompressResponseWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()
	if w.Header().Get(HeaderContentType) == "" {
		w.Header().Set(HeaderContentType, http.DetectContentType(b))
	}
	h.Del("Content-Length")
	return w.Writer.Write(b)
}

func (w *CompressResponseWriter) Flush() {
	w.Writer.(*gzip.Writer).Flush()
}

func (w *CompressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	L:
		for _, enc := range strings.Split(r.Header.Get(HeaderAcceptEncoding), ",") {
			switch strings.TrimSpace(enc) {
			case "gzip":
				w.Header().Set(HeaderContentEncoding, "gzip")
				w.Header().Add(HeaderVary, HeaderAcceptEncoding)

				gw, _ := gzip.NewWriterLevel(w, 6)
				defer gw.Close()

				h, hok := w.(http.Hijacker)
				if !hok { /* w is not Hijacker... oh well... */
					h = nil
				}

				f, fok := w.(http.Flusher)
				if !fok {
					f = nil
				}

				cn, cnok := w.(http.CloseNotifier)
				if !cnok {
					cn = nil
				}

				w = &CompressResponseWriter{
					Writer:         gw,
					ResponseWriter: w,
					Hijacker:       h,
					Flusher:        f,
					CloseNotifier:  cn,
				}

				break L
			case "deflate":
				w.Header().Set(HeaderContentEncoding, "deflate")
				w.Header().Add(HeaderVary, HeaderAcceptEncoding)

				fw, _ := flate.NewWriter(w, 6)
				defer fw.Close()

				h, hok := w.(http.Hijacker)
				if !hok { /* w is not Hijacker... oh well... */
					h = nil
				}

				f, fok := w.(http.Flusher)
				if !fok {
					f = nil
				}

				cn, cnok := w.(http.CloseNotifier)
				if !cnok {
					cn = nil
				}

				w = &CompressResponseWriter{
					Writer:         fw,
					ResponseWriter: w,
					Hijacker:       h,
					Flusher:        f,
					CloseNotifier:  cn,
				}

				break L
			}
		}

		next.ServeHTTP(w, r)
	})
}

func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderServer, fmt.Sprintf("Golang/mux %s", g.Version))
		next.ServeHTTP(w, r)
	})
}

func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(g.Config().Secret.Jwt.TokenName)
		for {
			// 没有提供令牌不处理
			if strings.TrimSpace(token) == "" {
				break
			}
			credential, err := tools.New().ParseToken(token)
			if err != nil {
				break
			}
			ctx := context.WithValue(r.Context(), g.CREDENTIAL_KEY, credential)
			r = r.WithContext(ctx)

			break
		}

		next.ServeHTTP(w, r)
	})
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(HeaderVary, HeaderOrigin)
		w.Header().Add(HeaderVary, HeaderAcceptEncoding)
		w.Header().Set(HeaderAccessControlAllowOrigin, "*")
		w.Header().Set(HeaderAccessControlAllowMethods, "*")
		w.Header().Set(HeaderAccessControlAllowHeaders, "Origin, X-Requested-With, Content-Type, Accept, Authentication, Allow-Credentials")
		w.Header().Set(HeaderAccessControlAllowCredentials, "true")

		next.ServeHTTP(w, r)
	})
}

// TODO:
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rc := gqlapi.ReturnCodeUnknowError
				log.Errorf("错误代码: %s, 错误信息: %s", rc, tools.PanicDetail())
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pool := &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 256))
			},
		}
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		defer pool.Put(buf)
		start := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		resp := NewResponse(w)
		next.ServeHTTP(resp, r)
		stop := time.Now()
		buf.WriteString(fmt.Sprintf("[I] %s ", time.Now().Format("2006-01-02 15:04:05")))
		buf.WriteString(fmt.Sprintf(`"remote_ip":"%s",`, r.RemoteAddr))
		buf.WriteString(fmt.Sprintf(`"host":"%s",`, r.Host))
		buf.WriteString(fmt.Sprintf(`"uri":"%s",`, r.RequestURI))
		buf.WriteString(fmt.Sprintf(`"method":"%s",`, r.Method))
		p := r.URL.Path
		if p == "" {
			p = "/"
		}
		buf.WriteString(fmt.Sprintf(`"path":"%s",`, p))
		buf.WriteString(fmt.Sprintf(`"protocal":"%s",`, r.Proto))
		buf.WriteString(fmt.Sprintf(`"referer":"%s",`, r.Referer()))
		buf.WriteString(fmt.Sprintf(`"user_agent":"%s",`, r.UserAgent()))
		buf.WriteString(fmt.Sprintf(`"status":"%d",`, resp.Status))
		buf.WriteString(fmt.Sprintf(`"bytes_in":"%d",`, r.ContentLength))
		buf.WriteString(fmt.Sprintf(`"bytes_out":"%d",`, resp.Size))
		buf.WriteString(fmt.Sprintf(`"latency":"%s"`, stop.Sub(start).String()))
		log.Println(buf.String())
	})
}

type (
	// Response wraps an http.ResponseWriter and implements its interface to be used
	// by an HTTP handler to construct an HTTP response.
	// See: https://golang.org/pkg/net/http/#ResponseWriter
	Response struct {
		beforeFuncs []func()
		afterFuncs  []func()
		Writer      http.ResponseWriter
		Status      int
		Size        int64
		Committed   bool
	}
)

// NewResponse creates a new instance of Response.
func NewResponse(w http.ResponseWriter) (r *Response) {
	return &Response{Writer: w}
}

// Header returns the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

// Before registers a function which is called just before the response is written.
func (r *Response) Before(fn func()) {
	r.beforeFuncs = append(r.beforeFuncs, fn)
}

// After registers a function which is called just after the response is written.
// If the `Content-Length` is unknown, none of the after function is executed.
func (r *Response) After(fn func()) {
	r.afterFuncs = append(r.afterFuncs, fn)
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) WriteHeader(code int) {
	if r.Committed {
		return
	}
	for _, fn := range r.beforeFuncs {
		fn()
	}
	r.Status = code
	r.Writer.WriteHeader(code)
	r.Committed = true
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *Response) Write(b []byte) (n int, err error) {
	if !r.Committed {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	for _, fn := range r.afterFuncs {
		fn()
	}
	return
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	r.Writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}

func (r *Response) reset(w http.ResponseWriter) {
	r.beforeFuncs = nil
	r.afterFuncs = nil
	r.Writer = w
	r.Size = 0
	r.Status = http.StatusOK
	r.Committed = false
}
