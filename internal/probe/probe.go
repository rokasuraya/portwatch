// Package probe provides a lightweight TCP/UDP reachability check
// that can be used to verify whether a specific port is currently
// accepting connections before recording it as open.
package probe

import (
	"context"
	"fmt"
	"net"
	"time"
)

// DefaultTimeout is used when no deadline is set on the context.
const DefaultTimeout = 2 * time.Second

// Result holds the outcome of a single probe attempt.
type Result struct {
	Host     string
	Port     int
	Protocol string
	Open     bool
	Latency  time.Duration
	Err      error
}

// Prober checks whether a host:port is reachable.
type Prober struct {
	timeout time.Duration
}

// New returns a Prober with the given timeout.
// If timeout is zero, DefaultTimeout is used.
func New(timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Prober{timeout: timeout}
}

// Check attempts to connect to host:port using the given protocol.
// Supported protocols are "tcp" and "udp".
func (p *Prober) Check(ctx context.Context, host string, port int, protocol string) Result {
	addr := fmt.Sprintf("%s:%d", host, port)

	deadline := p.timeout
	if dl, ok := ctx.Deadline(); ok {
		if remaining := time.Until(dl); remaining < deadline {
			deadline = remaining
		}
	}

	start := time.Now()
	conn, err := net.DialTimeout(protocol, addr, deadline)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Host:     host,
			Port:     port,
			Protocol: protocol,
			Open:     false,
			Latency:  latency,
			Err:      err,
		}
	}
	_ = conn.Close()
	return Result{
		Host:     host,
		Port:     port,
		Protocol: protocol,
		Open:     true,
		Latency:  latency,
	}
}
