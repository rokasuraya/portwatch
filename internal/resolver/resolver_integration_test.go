package resolver_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/resolver"
)

// TestConcurrentLookupAndRegister verifies that simultaneous reads and writes
// do not cause data races (run with -race).
func TestConcurrentLookupAndRegister(t *testing.T) {
	r := resolver.New(nil)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(port uint16) {
			defer wg.Done()
			r.Lookup(port, "tcp")
		}(uint16(i + 1))
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(port uint16) {
			defer wg.Done()
			r.Register(port, "tcp", "svc")
		}(uint16(i + 1024))
	}
	wg.Wait()
}

// TestRoundTrip_RegisterThenLookup ensures a registered service survives a
// subsequent lookup on the same instance.
func TestRoundTrip_RegisterThenLookup(t *testing.T) {
	r := resolver.New(nil)
	services := map[uint16]string{
		9300: "elasticsearch-transport",
		5601: "kibana",
		4222: "nats",
	}
	for port, svc := range services {
		r.Register(port, "tcp", svc)
	}
	for port, want := range services {
		got := r.Lookup(port, "tcp")
		if got != want {
			t.Errorf("port %d: expected %q, got %q", port, want, got)
		}
	}
}
