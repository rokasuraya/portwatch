package portping_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"portwatch/internal/portping"
)

// startListener binds an ephemeral TCP port and returns its port number and a
// cancel function that closes the listener.
func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	return port, func() { _ = ln.Close() }
}

func TestNew_DefaultTimeout(t *testing.T) {
	p := portping.New("127.0.0.1", 0)
	if p == nil {
		t.Fatal("expected non-nil Pinger")
	}
}

func TestPing_OpenPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := portping.New("127.0.0.1", time.Second)
	res := p.Ping(context.Background(), port)

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", res.Latency)
	}
	if res.Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", res.Protocol)
	}
}

func TestPing_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly unreachable / refused on loopback.
	p := portping.New("127.0.0.1", 200*time.Millisecond)
	res := p.Ping(context.Background(), 1)

	if res.Err == nil {
		t.Fatal("expected error for closed port")
	}
}

func TestPing_CancelledContext(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before dialling

	p := portping.New("127.0.0.1", time.Second)
	res := p.Ping(ctx, port)

	if res.Err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestPingAll_ReturnsOneResultPerPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	p := portping.New("127.0.0.1", time.Second)
	results := p.PingAll(context.Background(), []int{port, port})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Port != port {
			t.Errorf("unexpected port %d", r.Port)
		}
	}
}
