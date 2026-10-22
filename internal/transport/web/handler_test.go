package web

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func testFS() fs.FS {
	return fstest.MapFS{
		"index.html":     {Data: []byte("<html>app</html>")},
		"favicon.svg":    {Data: []byte("<svg/>")},
		"assets/app.js":  {Data: []byte("console.log(1)")},
		"assets/app.css": {Data: []byte("body{}")},
	}
}

func TestHandlerServesIndexAtRoot(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(testFS()).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "<html>app</html>", recorder.Body.String())
	require.Equal(t, "no-cache", recorder.Header().Get("Cache-Control"))
	require.Contains(t, recorder.Header().Get("Content-Type"), "text/html")
}

func TestHandlerFallsBackToIndexForClientRoute(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(testFS()).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/tickets/abc", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "<html>app</html>", recorder.Body.String())
}

func TestHandlerServesAssetsWithImmutableCache(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(testFS()).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/assets/app.js", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "console.log(1)", recorder.Body.String())
	require.Equal(t, "public, max-age=31536000, immutable", recorder.Header().Get("Cache-Control"))
}

func TestHandlerServesRootFiles(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(testFS()).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/favicon.svg", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "<svg/>", recorder.Body.String())
	require.Empty(t, recorder.Header().Get("Cache-Control"))
}

func TestHandlerReportsMissingBuild(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(fstest.MapFS{}).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Contains(t, recorder.Body.String(), "make build")
}

func TestHandlerRejectsWrites(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(testFS()).ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", nil))

	require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	require.Equal(t, "GET, HEAD", recorder.Header().Get("Allow"))
}

func TestHandlerHeadHasNoBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewHandler(testFS()).ServeHTTP(recorder, httptest.NewRequest(http.MethodHead, "/", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Empty(t, recorder.Body.String())
}
