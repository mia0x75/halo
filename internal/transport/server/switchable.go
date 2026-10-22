package server

import (
	"net/http"
	"sync"
)

// Switchable is a handler that can be replaced while the server runs, so the
// setup wizard can hand over to the normal routes without a restart.
type Switchable struct {
	mu sync.RWMutex
	h  http.Handler
}

// NewSwitchable builds a switchable handler serving the given handler.
func NewSwitchable(h http.Handler) *Switchable {
	return &Switchable{h: h}
}

// Set replaces the handler.
func (s *Switchable) Set(h http.Handler) {
	s.mu.Lock()
	s.h = h
	s.mu.Unlock()
}

// ServeHTTP dispatches to the current handler.
func (s *Switchable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	h := s.h
	s.mu.RUnlock()
	if h == nil {
		http.Error(w, "服务尚未就绪", http.StatusServiceUnavailable)
		return
	}
	h.ServeHTTP(w, r)
}
