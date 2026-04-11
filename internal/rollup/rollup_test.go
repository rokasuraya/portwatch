package rollup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/snapshot"
)

func entry(proto, addr string, port int) snapshot.Entry {
	return snapshot.Entry{Proto: proto, Addr: addr, Port: port}
}

func TestNew_ReturnsRollup(t *testing.T) {
	r := rollup.New(10*time.Millisecond, func(rollup.Summary) {})
	if r == nil {
		t.Fatal("expected non-nil rollup")
	}
}

func TestAdd_ImmediateFlushOnZeroWait(t *testing.T) {
	var got rollup.Summary
	r := rollup.New(0, func(s rollup.Summary) { got = s })

	r.Add([]snapshot.Entry{entry("tcp", "0.0.0.0", 8080)}, nil)

	if len(got.Opened) != 1 || got.Opened[0].Port != 8080 {
		t.Fatalf("expected port 8080 opened, got %+v", got)
	}
}

func TestAdd_BatchesWithinWindow(t *testing.T) {
	var mu sync.Mutex
	var calls int
	r := rollup.New(50*time.Millisecond, func(s rollup.Summary) {
		mu.Lock()
		calls++
		mu.Unlock()
	})

	r.Add([]snapshot.Entry{entry("tcp", "0.0.0.0", 80)}, nil)
	r.Add([]snapshot.Entry{entry("tcp", "0.0.0.0", 443)}, nil)

	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected 1 flush, got %d", calls)
	}
}

func TestAdd_NetCancelsOpenedThenClosed(t *testing.T) {
	var got rollup.Summary
	r := rollup.New(0, func(s rollup.Summary) { got = s })

	e := entry("tcp", "0.0.0.0", 9000)
	r.Add([]snapshot.Entry{e}, nil)
	got = rollup.Summary{}
	r.Add(nil, []snapshot.Entry{e})

	if len(got.Opened) != 0 {
		t.Fatalf("expected port removed from opened set, got %+v", got.Opened)
	}
	if len(got.Closed) != 1 {
		t.Fatalf("expected 1 closed entry, got %+v", got.Closed)
	}
}

func TestFlush_EmitsImmediately(t *testing.T) {
	var got rollup.Summary
	r := rollup.New(5*time.Second, func(s rollup.Summary) { got = s })

	r.Add(nil, []snapshot.Entry{entry("udp", "0.0.0.0", 53)})
	r.Flush()

	if len(got.Closed) != 1 || got.Closed[0].Port != 53 {
		t.Fatalf("expected port 53 closed after flush, got %+v", got)
	}
}

func TestFlush_NoOpWhenEmpty(t *testing.T) {
	called := false
	r := rollup.New(0, func(rollup.Summary) { called = true })
	r.Flush()
	if called {
		t.Fatal("expected no flush callback on empty buffer")
	}
}

func TestSummary_WindowEndSet(t *testing.T) {
	var got rollup.Summary
	r := rollup.New(0, func(s rollup.Summary) { got = s })
	before := time.Now()
	r.Add([]snapshot.Entry{entry("tcp", "127.0.0.1", 22)}, nil)
	if got.WindowEnd.Before(before) {
		t.Fatal("expected WindowEnd to be set after Add")
	}
}
