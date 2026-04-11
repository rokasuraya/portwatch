package healthcheck_test

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

// freePort returns an available TCP port on localhost.
func freePort(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()
	return addr
}

// startServer creates and starts a healthcheck server on a free port,
// registering a cleanup function to shut it down when the test completes.
func startServer(t *testing.T) *healthcheck.Server {
	t.Helper()
	addr := freePort(t)
	srv := healthcheck.New(addr)
	go func() { _ = srv.ListenAndServe() }()
	time.Sleep(50 * time.Millisecond) // allow server to start
	t.Cleanup(func() { srv.Shutdown() })
	return srv
}

func TestIntegration_ListenAndServe(t *testing.T) {
	srv := startServer(t)

	srv.RecordScan()
	srv.RecordScan()
	srv.RecordScan()

	resp, err := http.Get("http://" + srv.Addr() + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var status healthcheck.Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !status.OK {
		t.Error("expected ok=true")
	}
	if status.ScanCount != 3 {
		t.Errorf("expected scan_count=3, got %d", status.ScanCount)
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
	if status.LastScan.IsZero() {
		t.Error("expected last_scan to be populated")
	}
}
