package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestNotify_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{Opened: []int{8080, 9090}}
	alerts := n.Notify(diff)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
	for _, a := range alerts {
		if a.Level != alert.LevelAlert {
			t.Errorf("expected level ALERT for opened port, got %s", a.Level)
		}
	}

	out := buf.String()
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("output missing expected ports: %s", out)
	}
}

func TestNotify_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{Closed: []int{22, 443}}
	alerts := n.Notify(diff)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
	for _, a := range alerts {
		if a.Level != alert.LevelWarn {
			t.Errorf("expected level WARN for closed port, got %s", a.Level)
		}
	}

	out := buf.String()
	if !strings.Contains(out, "22") || !strings.Contains(out, "443") {
		t.Errorf("output missing expected ports: %s", out)
	}
}

func TestNotify_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	alerts := n.Notify(state.Diff{})

	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts for empty diff, got %d", len(alerts))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff")
	}
}

func TestNotify_MixedDiff(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{
		Opened: []int{3000},
		Closed: []int{8080},
	}
	alerts := n.Notify(diff)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}

	// Verify that opened ports produce ALERT level and closed ports produce WARN level.
	levels := make(map[string]int)
	for _, a := range alerts {
		levels[a.Level]++
	}
	if levels[alert.LevelAlert] != 1 {
		t.Errorf("expected 1 ALERT level alert, got %d", levels[alert.LevelAlert])
	}
	if levels[alert.LevelWarn] != 1 {
		t.Errorf("expected 1 WARN level alert, got %d", levels[alert.LevelWarn])
	}
}
