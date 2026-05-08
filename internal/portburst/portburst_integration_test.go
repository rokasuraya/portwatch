package portburst_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portburst"
	"github.com/user/portwatch/internal/snapshot"
)

func entries(ports ...int) []snapshot.Entry {
	out := make([]snapshot.Entry, len(ports))
	for i, p := range ports {
		out[i] = snapshot.Entry{Port: p, Protocol: "tcp"}
	}
	return out
}

// TestIntegration_BurstAcrossMultipleObservations verifies that events
// accumulate across successive Observe calls within the window.
func TestIntegration_BurstAcrossMultipleObservations(t *testing.T) {
	var buf bytes.Buffer
	d := portburst.New(4, time.Minute, &buf)

	d.Observe(snapshot.Diff{Opened: entries(80, 443)})
	if strings.Contains(buf.String(), "burst") {
		t.Fatal("should not warn yet")
	}
	d.Observe(snapshot.Diff{Opened: entries(8080, 9090)})
	if !strings.Contains(buf.String(), "burst detected") {
		t.Fatalf("expected burst warning after accumulation, got: %s", buf.String())
	}
}

// TestIntegration_ResetAllowsNewBurstCycle ensures that after a reset the
// detector starts fresh and does not carry over previous events.
func TestIntegration_ResetAllowsNewBurstCycle(t *testing.T) {
	var buf bytes.Buffer
	d := portburst.New(3, time.Minute, &buf)

	d.Observe(snapshot.Diff{Opened: entries(80, 443, 22)})
	buf.Reset()
	d.Reset()

	// Two more openings should not trigger the threshold of 3.
	d.Observe(snapshot.Diff{Opened: entries(8080, 9090)})
	if strings.Contains(buf.String(), "burst") {
		t.Fatalf("unexpected burst after reset: %s", buf.String())
	}
}
