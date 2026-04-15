package snapshot_test

import (
	"testing"

	"portwatch/internal/snapshot"
)

func TestNew_CopiesEntries(t *testing.T) {
	orig := []snapshot.Entry{{Protocol: "tcp", Port: 80}}
	s := snapshot.New(orig)
	orig[0].Port = 9999
	if s.Entries[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", s.Entries[0].Port)
	}
}

func TestNew_StampsTime(t *testing.T) {
	s := snapshot.New(nil)
	if s.CapturedAt.IsZero() {
		t.Fatal("expected non-zero CapturedAt")
	}
}

func TestEntry_String(t *testing.T) {
	e := snapshot.Entry{Protocol: "udp", Port: 53}
	if got := e.String(); got != "udp:53" {
		t.Fatalf("expected udp:53, got %s", got)
	}
}

func TestCompare_DetectsOpened(t *testing.T) {
	prev := snapshot.New([]snapshot.Entry{})
	curr := snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 443}})

	d := snapshot.Compare(prev, curr)
	if len(d.Opened) != 1 || d.Opened[0].Port != 443 {
		t.Fatalf("expected one opened port 443, got %+v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Fatalf("expected no closed ports, got %+v", d.Closed)
	}
}

func TestCompare_DetectsClosed(t *testing.T) {
	prev := snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 22}})
	curr := snapshot.New([]snapshot.Entry{})

	d := snapshot.Compare(prev, curr)
	if len(d.Closed) != 1 || d.Closed[0].Port != 22 {
		t.Fatalf("expected one closed port 22, got %+v", d.Closed)
	}
	if len(d.Opened) != 0 {
		t.Fatalf("expected no opened ports, got %+v", d.Opened)
	}
}

func TestCompare_NoChange(t *testing.T) {
	entries := []snapshot.Entry{{Protocol: "tcp", Port: 8080}}
	prev := snapshot.New(entries)
	curr := snapshot.New(entries)

	d := snapshot.Compare(prev, curr)
	if !d.IsEmpty() {
		t.Fatalf("expected empty diff, got %+v", d)
	}
}

func TestCompare_MultipleChanges(t *testing.T) {
	prev := snapshot.New([]snapshot.Entry{
		{Protocol: "tcp", Port: 22},
		{Protocol: "tcp", Port: 80},
	})
	curr := snapshot.New([]snapshot.Entry{
		{Protocol: "tcp", Port: 80},
		{Protocol: "tcp", Port: 443},
	})

	d := snapshot.Compare(prev, curr)
	if len(d.Opened) != 1 || d.Opened[0].Port != 443 {
		t.Fatalf("expected one opened port 443, got %+v", d.Opened)
	}
	if len(d.Closed) != 1 || d.Closed[0].Port != 22 {
		t.Fatalf("expected one closed port 22, got %+v", d.Closed)
	}
}

func TestCompare_ProtocolDistinct(t *testing.T) {
	// tcp:80 and udp:80 are distinct entries; opening udp:80 should not
	// suppress the detection of a new port just because tcp:80 existed.
	prev := snapshot.New([]snapshot.Entry{{Protocol: "tcp", Port: 80}})
	curr := snapshot.New([]snapshot.Entry{
		{Protocol: "tcp", Port: 80},
		{Protocol: "udp", Port: 80},
	})

	d := snapshot.Compare(prev, curr)
	if len(d.Opened) != 1 || d.Opened[0].Protocol != "udp" || d.Opened[0].Port != 80 {
		t.Fatalf("expected one opened udp:80, got %+v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Fatalf("expected no closed ports, got %+v", d.Closed)
	}
}

func TestDiff_IsEmpty_True(t *testing.T) {
	d := snapshot.Diff{}
	if !d.IsEmpty() {
		t.Fatal("expected IsEmpty to return true")
	}
}

func TestDiff_IsEmpty_False(t *testing.T) {
	d := snapshot.Diff{Opened: []snapshot.Entry{{Protocol: "tcp", Port: 80}}}
	if d.IsEmpty() {
		t.Fatal("expected IsEmpty to return false")
	}
}
