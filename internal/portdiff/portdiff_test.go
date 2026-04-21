package portdiff_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portdiff"
	"github.com/user/portwatch/internal/snapshot"
)

type stubLabeler struct{ m map[string]string }

func (s *stubLabeler) Label(port int, proto string) string {
	return s.m[proto+"/"+fmt.Sprintf("%d", port)]
}

import "fmt"

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries, time.Now())
}

func TestCompute_EmptyBothSnapshots(t *testing.T) {
	d := portdiff.Compute(makeSnap(nil), makeSnap(nil), nil)
	if !d.IsEmpty() {
		t.Fatal("expected empty diff")
	}
}

func TestCompute_DetectsOpened(t *testing.T) {
	prev := makeSnap(nil)
	next := makeSnap([]snapshot.Entry{{Port: 80, Protocol: "tcp"}})
	d := portdiff.Compute(prev, next, nil)
	if len(d.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(d.Opened))
	}
	if d.Opened[0].Port != 80 {
		t.Errorf("expected port 80, got %d", d.Opened[0].Port)
	}
	if d.Opened[0].Op != "opened" {
		t.Errorf("expected op=opened, got %s", d.Opened[0].Op)
	}
}

func TestCompute_DetectsClosed(t *testing.T) {
	prev := makeSnap([]snapshot.Entry{{Port: 443, Protocol: "tcp"}})
	next := makeSnap(nil)
	d := portdiff.Compute(prev, next, nil)
	if len(d.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(d.Closed))
	}
	if d.Closed[0].Port != 443 {
		t.Errorf("expected port 443, got %d", d.Closed[0].Port)
	}
}

func TestCompute_AppliesLabel(t *testing.T) {
	l := &stubLabeler{m: map[string]string{"tcp/22": "ssh"}}
	prev := makeSnap(nil)
	next := makeSnap([]snapshot.Entry{{Port: 22, Protocol: "tcp"}})
	d := portdiff.Compute(prev, next, l)
	if d.Opened[0].Label != "ssh" {
		t.Errorf("expected label ssh, got %q", d.Opened[0].Label)
	}
}

func TestEntry_String_WithLabel(t *testing.T) {
	e := portdiff.Entry{Op: "opened", Port: 22, Protocol: "tcp", Label: "ssh"}
	got := e.String()
	if got != "opened tcp/22 (ssh)" {
		t.Errorf("unexpected string: %s", got)
	}
}

func TestEntry_String_NoLabel(t *testing.T) {
	e := portdiff.Entry{Op: "closed", Port: 8080, Protocol: "tcp"}
	got := e.String()
	if got != "closed tcp/8080" {
		t.Errorf("unexpected string: %s", got)
	}
}

func TestDiff_Summary_NoChanges(t *testing.T) {
	d := portdiff.Diff{}
	if d.Summary() != "no changes" {
		t.Errorf("unexpected summary: %s", d.Summary())
	}
}

func TestDiff_Summary_Mixed(t *testing.T) {
	d := portdiff.Diff{
		Opened: []portdiff.Entry{{}, {}},
		Closed: []portdiff.Entry{{}},
	}
	got := d.Summary()
	if got != "2 opened, 1 closed" {
		t.Errorf("unexpected summary: %s", got)
	}
}
