package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/bootstrap"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/transport/web"
)

// SetupHandler serves the first-run setup API. It is only mounted while the
// server is not configured yet, so the endpoints disappear once the wizard is
// done.
type SetupHandler struct {
	// OnDone is called after the setup succeeded so the server can switch to
	// the normal routes.
	OnDone func(*app.App) error

	mu      sync.Mutex
	running bool
}

// NewSetupHandler builds the setup handler.
func NewSetupHandler(onDone func(*app.App) error) *SetupHandler {
	return &SetupHandler{OnDone: onDone}
}

// SetupRoutes mounts the setup endpoints and the frontend on the router. It is
// the whole route table of a server that is not configured yet.
func SetupRoutes(r *mux.Router, handler *SetupHandler) {
	r.Use(Headers)
	r.Use(Context)
	r.Use(Cors)
	r.Use(Logging)
	r.Use(Compress)

	handler.Register(r)

	dist, err := web.Dist()
	if err != nil {
		log.Printf("[W] 前端资源不可用，仅提供初始化接口: %v", err)
		return
	}
	r.PathPrefix("/").Handler(web.NewHandler(dist))
}

// Register mounts the setup endpoints on the router.
func (h *SetupHandler) Register(r *mux.Router) {
	r.HandleFunc("/api/setup/status", h.status)
	r.HandleFunc("/api/setup", h.setup)
}

// statusResponse tells the frontend whether the wizard has to run.
type statusResponse struct {
	Configured bool   `json:"configured"`
	ConfigFile string `json:"configFile"`
	Driver     string `json:"driver"`
	Addr       string `json:"addr"`
}

// status reports the setup state. Secrets are never included.
func (h *SetupHandler) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeSetupError(w, http.StatusMethodNotAllowed, "只支持 GET 请求")
		return
	}
	response := statusResponse{
		Configured: config.Configured(),
		// 还没有配置文件时，这里是向导将要写入的位置。
		ConfigFile: config.SavePath,
	}
	if cfg := config.Config(); cfg != nil && cfg.Database != nil {
		response.Driver = cfg.Database.Driver
		response.Addr = cfg.Database.Addr
	}
	writeSetupJSON(w, http.StatusOK, response)
}

// setup runs the wizard: it writes the config file, migrates the database and
// creates the initial administrator.
func (h *SetupHandler) setup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeSetupError(w, http.StatusMethodNotAllowed, "只支持 POST 请求")
		return
	}

	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		writeSetupError(w, http.StatusConflict, "初始化正在进行中")
		return
	}
	h.running = true
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		h.running = false
		h.mu.Unlock()
	}()

	var request bootstrap.Request
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := decoder.Decode(&request); err != nil {
		writeSetupError(w, http.StatusBadRequest, "请求格式错误: "+err.Error())
		return
	}

	application, err := bootstrap.Run(r.Context(), request)
	if err != nil {
		log.Errorf("[E] setup failed: %v", err)
		writeSetupError(w, http.StatusBadRequest, err.Error())
		return
	}

	if h.OnDone != nil {
		if err := h.OnDone(application); err != nil {
			_ = application.Close()
			writeSetupError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeSetupJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// writeSetupJSON writes a JSON response.
func writeSetupJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set(HeaderContentType, "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Errorf("[E] write setup response failed: %v", err)
	}
}

// writeSetupError writes a JSON error the frontend can show verbatim.
func writeSetupError(w http.ResponseWriter, status int, message string) {
	writeSetupJSON(w, status, map[string]any{"error": strings.TrimSpace(message)})
}
