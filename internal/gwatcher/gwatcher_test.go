package gwatcher_test

import (
	"bytes"
	"testing"

	"portwatch/internal/gwatcher"
	"portwatch/internal/portgroup"
	"portwatch/internal/scanner"
	"portwatch/internal/snapshot"
)

func buildMatcher(t *testing.T) *portgroup.Matcher {
	t.Helper()
	reg := portgroup.New()
	_ = reg.Define("web", []portgroup.Entry{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	})
	_ = reg.Define("db", []portgroup.Entry{
		{Port: 5432, Proto: "tcp"},
	})
	return portgroup.NewMatcher(reg)
}

func snap(entries []scanner.Entry) *snapshot.Snapshot {
	return snapshot.New(entries)
}

func TestNew_DefaultsToStdout(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, nil)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestObserve_NilSnapshot(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, &bytes.Buffer{})
	events := w.Observe(nil)
	if len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
}

func TestObserve_JoinsGroupOnFirstSeen(t *testing.T) {
	m := buildMatcher(t)
	buf := &bytes.Buffer{}
	w := gwatcher.New(m, buf)

	s := snap([]scanner.Entry{{Port: 80, Proto: "tcp", Open: true}})
	events := w.Observe(s)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if !events[0].Joined || events[0].Group != "web" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestObserve_NoChangeNoEvents(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, &bytes.Buffer{})
	s := snap([]scanner.Entry{{Port: 80, Proto: "tcp", Open: true}})
	w.Observe(s)
	events := w.Observe(s)
	if len(events) != 0 {
		t.Fatalf("expected 0 events on stable snapshot, got %d", len(events))
	}
}

func TestObserve_LeavesGroupWhenPortRemoved(t *testing.T) {
	m := buildMatcher(t)
	buf := &bytes.Buffer{}
	w := gwatcher.New(m, buf)

	w.Observe(snap([]scanner.Entry{{Port: 5432, Proto: "tcp", Open: true}}))
	events := w.Observe(snap([]scanner.Entry{}))

	if len(events) != 1 {
		t.Fatalf("expected 1 leave event, got %d", len(events))
	}
	if events[0].Joined || events[0].Group != "db" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestObserve_UnknownPortNoEvents(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, &bytes.Buffer{})
	s := snap([]scanner.Entry{{Port: 9999, Proto: "tcp", Open: true}})
	events := w.Observe(s)
	if len(events) != 0 {
		t.Fatalf("expected 0 events for ungrouped port, got %d", len(events))
	}
}
