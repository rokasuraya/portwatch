package scanner

import (
	"net"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	timeout := 500 * time.Millisecond
	s := New(timeout)

	if s == nil {
		t.Fatal("expected scanner to be created")
	}

	if s.timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, s.timeout)
	}
}

func TestScanTCPPort(t *testing.T) {
	s := New(100 * time.Millisecond)

	// Start a test server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	// Test open port
	p, err := s.ScanTCPPort(port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.State != "open" {
		t.Errorf("expected port %d to be open, got %s", port, p.State)
	}

	if p.Number != port {
		t.Errorf("expected port number %d, got %d", port, p.Number)
	}

	if p.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", p.Protocol)
	}

	// Test closed port
	p, err = s.ScanTCPPort(65534)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.State != "closed" {
		t.Errorf("expected port 65534 to be closed, got %s", p.State)
	}
}

func TestScanPortRange(t *testing.T) {
	s := New(100 * time.Millisecond)

	// Test invalid ranges
	_, err := s.ScanPortRange(0, 100)
	if err == nil {
		t.Error("expected error for invalid start port")
	}

	_, err = s.ScanPortRange(1, 65536)
	if err == nil {
		t.Error("expected error for invalid end port")
	}

	_, err = s.ScanPortRange(100, 50)
	if err == nil {
		t.Error("expected error for reversed range")
	}

	// Test valid range
	ports, err := s.ScanPortRange(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ports == nil {
		t.Error("expected non-nil result")
	}
}

func TestScanPortRangeResultCount(t *testing.T) {
	s := New(100 * time.Millisecond)

	const start, end = 1, 10
	ports, err := s.ScanPortRange(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := end - start + 1
	if len(ports) != expected {
		t.Errorf("expected %d port results, got %d", expected, len(ports))
	}
}
