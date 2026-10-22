package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/platform/config"
)

func setupRouter(t *testing.T) *mux.Router {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"listen": ":8443"}`), 0o600))
	config.ParseConfig(path)

	router := mux.NewRouter()
	SetupRoutes(router, NewSetupHandler(nil))
	return router
}

func TestSetupStatusReportsUnconfigured(t *testing.T) {
	router := setupRouter(t)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/setup/status", nil))
	require.Equal(t, http.StatusOK, recorder.Code)

	var status struct {
		Configured bool   `json:"configured"`
		ConfigFile string `json:"configFile"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &status))
	require.False(t, status.Configured)
	require.NotEmpty(t, status.ConfigFile)
}

func TestSetupRejectsBrokenPayload(t *testing.T) {
	router := setupRouter(t)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/setup", strings.NewReader("{"))
	router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSetupReportsDatabaseFailure(t *testing.T) {
	router := setupRouter(t)

	body := `{"database":{"driver":"mysql","addr":"root@tcp(127.0.0.1:1)/halo"},
		"secret":{"secret_key":"0123456789abcdef0123456789abcdef","crypto":"0123456789abcdef"},
		"admin":{"email":"admin@example.com","password":"halo-admin"}}`
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/setup", strings.NewReader(body)))
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "元数据库")
}
