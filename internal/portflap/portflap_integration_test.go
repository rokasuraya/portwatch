package portflap_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portflap"
	"github.com/user/portwatch/internal/snapshot"
)

func TestIntegration_RapidFlapTriggersThenSettles(t *testing.T) {
	var buf bytes.Buffer
	d := portflap.New(3, 5*time.Second)
	d.SetOutput(&buf)

	e := snapshot.Entry{Port: 8443, Protocol: "tcp"}

	// Rapid flap: 3 transitions — should warn.
	for i := 0; i < 3; i++ {
		d.Observe([]snapshot.Entry{e}, nil)
	}
	if !strings.Contains(buf.String(), "8443/tcp") {
		t.Fatalf("expected flap warning, got: %q", buf.String())
	}

	// Reset simulates a new observation cycle.
	d.Reset()
	buf.Reset()

	// Single transition after reset — no warning expected.
	d.Observe([]snapshot.Entry{e}, nil)
	if buf.Len() != 0 {
		t.Errorf("expected no warning after reset; got: %s", buf.String())
	}
}

func TestIntegration_MultiplePortsIndependent(t *testing.T) {
	var buf bytes.Buffer
	d := portflap.New(2, time.Minute)
	d.SetOutput(&buf)

	a := snapshot.Entry{Port: 80, Protocol: "tcp"}
	b := snapshot.Entry{Port: 443, Protocol: "tcp"}

	// Port 80 flaps twice — should warn.
	d.Observe([]snapshot.Entry{a}, nil)
	d.Observe(nil, []snapshot.Entry{a})

	// Port 443 only once — should not warn.
	d.Observe([]snapshot.Entry{b}, nil)

	out := buf.String()
	if !strings.Contains(out, "80/tcp") {
		t.Errorf("expected warning for 80/tcp; got: %s", out)
	}
	if strings.Contains(out, "443/tcp") {
		t.Errorf("unexpected warning for 443/tcp; got: %s", out)
	}
}
