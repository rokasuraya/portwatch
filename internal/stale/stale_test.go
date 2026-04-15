package stale

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(entries []snapshot.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_DefaultsToStderr(t *testing.T) {
	d := New(time.Minute, nil)
	if d.out == nil {
		t.Fatal("expected non-nil writer")
	}
	if d.maxAge != time.Minute {
		t.Fatalf("expected maxAge=1m, got %s", d.maxAge)
	}
}

func TestObserve_RecordsNewPorts(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New(time.Minute, buf)

	snap := makeSnap([]snapshot.Entry{
		{Proto: "tcp", Port: 80},
		{Proto: "tcp", Port: 443},
	})
	d.Observe(snap)

	if len(d.seen) != 2 {
		t.Fatalf("expected 2 tracked ports, got %d", len(d.seen))
	}
}

func TestObserve_RemovesClosedPorts(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New(time.Minute, buf)

	d.Observe(makeSnap([]snapshot.Entry{{Proto: "tcp", Port: 80}, {Proto: "tcp", Port: 9000}}))
	d.Observe(makeSnap([]snapshot.Entry{{Proto: "tcp", Port: 80}}))

	if len(d.seen) != 1 {
		t.Fatalf("expected 1 tracked port after closure, got %d", len(d.seen))
	}
	if _, ok := d.seen[portKey("tcp", 80)]; !ok {
		t.Fatal("expected port 80 to still be tracked")
	}
}

func TestCheck_NoWarningWhenFresh(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New(time.Hour, buf)

	d.Observe(makeSnap([]snapshot.Entry{{Proto: "tcp", Port: 22}}))
	count := d.Check()

	if count != 0 {
		t.Fatalf("expected 0 stale ports, got %d", count)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}

func TestCheck_WarnsForStalePort(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New(time.Millisecond, buf)

	d.Observe(makeSnap([]snapshot.Entry{{Proto: "tcp", Port: 8080}}))

	// backdate the entry so it appears old
	k := portKey("tcp", 8080)
	e := d.seen[k]
	e.FirstAt = e.FirstAt.Add(-time.Hour)
	d.seen[k] = e

	count := d.Check()
	if count != 1 {
		t.Fatalf("expected 1 stale port, got %d", count)
	}
	if !strings.Contains(buf.String(), "8080") {
		t.Fatalf("expected port 8080 in output, got: %s", buf.String())
	}
}

func TestCheck_EmptyTrackerReturnsZero(t *testing.T) {
	buf := &bytes.Buffer{}
	d := New(time.Minute, buf)
	if d.Check() != 0 {
		t.Fatal("expected 0 for empty detector")
	}
}
