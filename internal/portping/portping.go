// Package portping provides a lightweight round-trip latency sampler for
// individual ports. It dials a TCP connection and measures the time to
// establish the connection, storing the result as a duration.
package portping

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single ping attempt.
type Result struct {
	Port     int
	Protocol string
	Latency  time.Duration
	Err      error
}

// Pinger measures TCP dial latency for a given host and port.
type Pinger struct {
	timeout time.Duration
	host    string
}

// New returns a Pinger that connects to host with the given timeout.
// If timeout is zero it defaults to 2 seconds.
func New(host string, timeout time.Duration) *Pinger {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Pinger{host: host, timeout: timeout}
}

// Ping dials the given TCP port and returns the round-trip latency.
func (p *Pinger) Ping(ctx context.Context, port int) Result {
	addr := fmt.Sprintf("%s:%d", p.host, port)

	dialCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	start := time.Now()
	conn, err := (&net.Dialer{}).DialContext(dialCtx, "tcp", addr)
	latency := time.Since(start)

	if err != nil {
		return Result{Port: port, Protocol: "tcp", Err: err}
	}
	_ = conn.Close()

	return Result{Port: port, Protocol: "tcp", Latency: latency}
}

// PingAll pings each port in the provided slice and returns one Result per
// port. Ports are probed sequentially; the caller may parallelise if needed.
func (p *Pinger) PingAll(ctx context.Context, ports []int) []Result {
	out := make([]Result, 0, len(ports))
	for _, port := range ports {
		out = append(out, p.Ping(ctx, port))
	}
	return out
}
