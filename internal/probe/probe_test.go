package probe_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"portwatch/internal/probe"
)

// startTCPListener binds an ephemeral TCP port and returns its port number
// along with a closer function.
func startTCPListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	// net.TCPAddr.Port is already an int — use type assertion directly
	return ln.Addr().(*net.TCPAddr).Port, func() { _ = ln.Close() }
}

func TestNew_DefaultTimeout(t *testing.T) {
	p := probe.New(0)
	if p == nil {
		t.Fatal("expected non-nil Prober")
	}
}

func TestCheck_OpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	p := probe.New(2 * time.Second)
	res := p.Check(context.Background(), "127.0.0.1", port, "tcp")

	if !res.Open {
		t.Errorf("expected port %d to be open, got err: %v", port, res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
	if res.Err != nil {
		t.Errorf("unexpected error: %v", res.Err)
	}
}

func TestCheck_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly not open in a test environment.
	p := probe.New(500 * time.Millisecond)
	res := p.Check(context.Background(), "127.0.0.1", 1, "tcp")

	if res.Open {
		t.Skip("port 1 unexpectedly open; skipping")
	}
	if res.Err == nil {
		t.Error("expected an error for closed port")
	}
}

func TestCheck_ContextDeadlineRespected(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	p := probe.New(10 * time.Second)
	// 192.0.2.1 is TEST-NET; connections will time out.
	res := p.Check(ctx, "192.0.2.1", 80, "tcp")

	if res.Open {
		t.Error("expected port to be reported closed on timeout")
	}
}

func TestCheck_ResultFields(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	p := probe.New(2 * time.Second)
	res := p.Check(context.Background(), "127.0.0.1", port, "tcp")

	if res.Host != "127.0.0.1" {
		t.Errorf("host: got %q, want %q", res.Host, "127.0.0.1")
	}
	if res.Port != port {
		t.Errorf("port: got %d, want %d", res.Port, port)
	}
	if res.Protocol != "tcp" {
		t.Errorf("protocol: got %q, want \"tcp\"", res.Protocol)
	}
}
