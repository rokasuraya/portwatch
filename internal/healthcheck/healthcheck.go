// Package healthcheck exposes a simple HTTP endpoint that reports the
// current liveness of the portwatch daemon.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the data returned by the health endpoint.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	LastScan  time.Time `json:"last_scan,omitempty"`
	ScanCount int64     `json:"scan_count"`
}

// Server is a lightweight HTTP server that serves health information.
type Server struct {
	addr      string
	start     time.Time
	lastScan  atomic.Value // stores time.Time
	scanCount atomic.Int64
	server    *http.Server
}

// New creates a new healthcheck Server bound to addr (e.g. ":9090").
func New(addr string) *Server {
	s := &Server{
		addr:  addr,
		start: time.Now(),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{Addr: addr, Handler: mux}
	return s
}

// RecordScan updates the last-scan timestamp and increments the scan counter.
func (s *Server) RecordScan() {
	s.lastScan.Store(time.Now())
	s.scanCount.Add(1)
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown() error {
	return s.server.Close()
}

// CurrentStatus returns a snapshot of the current health status without
// requiring an HTTP request. This is useful for logging or internal checks.
func (s *Server) CurrentStatus() Status {
	status := Status{
		OK:        true,
		Uptime:    time.Since(s.start).Truncate(time.Second).String(),
		ScanCount: s.scanCount.Load(),
	}
	if t, ok := s.lastScan.Load().(time.Time); ok && !t.IsZero() {
		status.LastScan = t
	}
	return status
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	status := s.CurrentStatus()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
