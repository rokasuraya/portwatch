// Package probe performs on-demand reachability checks against
// individual host/port/protocol combinations.
//
// It is intentionally separate from the bulk scanner so that the
// daemon can re-verify a specific port before raising an alert,
// reducing false positives caused by transient scan errors.
//
// Usage:
//
//	p := probe.New(2 * time.Second)
//	result := p.Check(ctx, "127.0.0.1", 22, "tcp")
//	if result.Open {
//	    fmt.Println("port is open, latency:", result.Latency)
//	}
package probe
