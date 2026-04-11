package coalesce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/coalesce"
	"github.com/user/portwatch/internal/snapshot"
)

func entry(port int, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Protocol: proto}
}

func TestNew_ReturnsCoalescer(t *testing.T) {
	c := coalesce.New(50*time.Millisecond, func(_, _ []snapshot.Entry) {})
	if c == nil {
		t.Fatal("expected non-nil coalescer")
	}
}

func TestAdd_FlushesAfterWait(t *testing.T) {
	var mu sync.Mutex
	var gotOpened []snapshot.Entry

	c := coalesce.New(50*time.Millisecond, func(opened, _ []snapshot.Entry) {
		mu.Lock()
		gotOpened = append(gotOpened, opened...)
		mu.Unlock()
	})

	c.Add([]snapshot.Entry{entry(8080, "tcp")}, nil)
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(gotOpened) != 1 || gotOpened[0].Port != 8080 {
		t.Fatalf("expected port 8080 opened, got %v", gotOpened)
	}
}

func TestAdd_CoalescesBurst(t *testing.T) {
	var mu sync.Mutex
	var calls int
	var totalOpened []snapshot.Entry

	c := coalesce.New(60*time.Millisecond, func(opened, _ []snapshot.Entry) {
		mu.Lock()
		calls++
		totalOpened = append(totalOpened, opened...)
		mu.Unlock()
	})

	c.Add([]snapshot.Entry{entry(80, "tcp")}, nil)
	time.Sleep(20 * time.Millisecond)
	c.Add([]snapshot.Entry{entry(443, "tcp")}, nil)
	time.Sleep(20 * time.Millisecond)
	c.Add([]snapshot.Entry{entry(9090, "tcp")}, nil)
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected 1 flush call, got %d", calls)
	}
	if len(totalOpened) != 3 {
		t.Fatalf("expected 3 opened entries, got %d", len(totalOpened))
	}
}

func TestFlush_EmitsImmediately(t *testing.T) {
	var mu sync.Mutex
	var gotClosed []snapshot.Entry

	c := coalesce.New(500*time.Millisecond, func(_, closed []snapshot.Entry) {
		mu.Lock()
		gotClosed = append(gotClosed, closed...)
		mu.Unlock()
	})

	c.Add(nil, []snapshot.Entry{entry(22, "tcp")})
	flushed := c.Flush()

	if !flushed {
		t.Fatal("expected Flush to return true")
	}
	mu.Lock()
	defer mu.Unlock()
	if len(gotClosed) != 1 || gotClosed[0].Port != 22 {
		t.Fatalf("expected port 22 closed, got %v", gotClosed)
	}
}

func TestFlush_ReturnsFalseWhenEmpty(t *testing.T) {
	c := coalesce.New(50*time.Millisecond, func(_, _ []snapshot.Entry) {})
	if c.Flush() {
		t.Fatal("expected Flush to return false when no pending entries")
	}
}
