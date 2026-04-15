package envelope_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/envelope"
	"github.com/yourorg/portwatch/internal/snapshot"
)

func makeEntries(ports ...int) []snapshot.Entry {
	out := make([]snapshot.Entry, len(ports))
	for i, p := range ports {
		out[i] = snapshot.Entry{Port: p, Protocol: "tcp"}
	}
	return out
}

func TestNew_SetsFields(t *testing.T) {
	opened := makeEntries(80, 443)
	closed := makeEntries(8080)
	dur := 42 * time.Millisecond
	labels := map[string]string{"host": "localhost"}

	e := envelope.New("abc123", opened, closed, dur, labels)

	if e.ID != "abc123" {
		t.Fatalf("expected ID abc123, got %s", e.ID)
	}
	if e.ScanDuration != dur {
		t.Fatalf("expected duration %v, got %v", dur, e.ScanDuration)
	}
	if len(e.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(e.Opened))
	}
	if len(e.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(e.Closed))
	}
	if e.Labels["host"] != "localhost" {
		t.Fatalf("expected label host=localhost")
	}
	if e.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestNew_IsolatesSlices(t *testing.T) {
	opened := makeEntries(22)
	e := envelope.New("x", opened, nil, 0, nil)
	opened[0].Port = 9999
	if e.Opened[0].Port == 9999 {
		t.Fatal("envelope should not share backing array with caller")
	}
}

func TestNew_NilLabels(t *testing.T) {
	e := envelope.New("y", nil, nil, 0, nil)
	if e.Labels == nil {
		t.Fatal("Labels should never be nil after New")
	}
}

func TestIsEmpty_TrueWhenNoDiff(t *testing.T) {
	e := envelope.New("z", nil, nil, 0, nil)
	if !e.IsEmpty() {
		t.Fatal("expected IsEmpty true")
	}
}

func TestIsEmpty_FalseWhenOpened(t *testing.T) {
	e := envelope.New("z", makeEntries(80), nil, 0, nil)
	if e.IsEmpty() {
		t.Fatal("expected IsEmpty false")
	}
}

func TestAddLabel_SetsValue(t *testing.T) {
	e := envelope.New("id", nil, nil, 0, nil)
	e.AddLabel("env", "prod")
	if e.Labels["env"] != "prod" {
		t.Fatalf("expected label env=prod, got %q", e.Labels["env"])
	}
}

func TestAddLabel_OverwritesExisting(t *testing.T) {
	e := envelope.New("id", nil, nil, 0, map[string]string{"env": "dev"})
	e.AddLabel("env", "prod")
	if e.Labels["env"] != "prod" {
		t.Fatalf("expected overwritten label env=prod")
	}
}
