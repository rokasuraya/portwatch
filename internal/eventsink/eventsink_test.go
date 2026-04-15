package eventsink_test

import (
	"sync/atomic"
	"testing"

	"github.com/user/portwatch/internal/eventsink"
	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(port uint16, proto string) snapshot.Entry {
	return snapshot.Entry{Port: port, Proto: proto}
}

func TestNew_ReturnsEmptySink(t *testing.T) {
	s := eventsink.New()
	if s.Len() != 0 {
		t.Fatalf("expected 0 handlers, got %d", s.Len())
	}
}

func TestRegister_IncreasesLen(t *testing.T) {
	s := eventsink.New()
	s.Register("a", func(_, _ []snapshot.Entry) {})
	if s.Len() != 1 {
		t.Fatalf("expected 1 handler, got %d", s.Len())
	}
}

func TestRegister_ReplacesExistingName(t *testing.T) {
	s := eventsink.New()
	s.Register("a", func(_, _ []snapshot.Entry) {})
	s.Register("a", func(_, _ []snapshot.Entry) {})
	if s.Len() != 1 {
		t.Fatalf("expected 1 handler after replace, got %d", s.Len())
	}
}

func TestRegister_NilRemovesHandler(t *testing.T) {
	s := eventsink.New()
	s.Register("a", func(_, _ []snapshot.Entry) {})
	s.Register("a", nil)
	if s.Len() != 0 {
		t.Fatalf("expected 0 handlers after nil register, got %d", s.Len())
	}
}

func TestUnregister_RemovesHandler(t *testing.T) {
	s := eventsink.New()
	s.Register("a", func(_, _ []snapshot.Entry) {})
	s.Unregister("a")
	if s.Len() != 0 {
		t.Fatalf("expected 0 handlers after unregister, got %d", s.Len())
	}
}

func TestUnregister_NoopForUnknownName(t *testing.T) {
	s := eventsink.New()
	s.Unregister("nonexistent") // must not panic
}

func TestEmit_CallsAllHandlers(t *testing.T) {
	s := eventsink.New()
	var countA, countB atomic.Int32
	s.Register("a", func(_, _ []snapshot.Entry) { countA.Add(1) })
	s.Register("b", func(_, _ []snapshot.Entry) { countB.Add(1) })

	opened := []snapshot.Entry{makeEntry(80, "tcp")}
	closed := []snapshot.Entry{makeEntry(443, "tcp")}
	s.Emit(opened, closed)

	if countA.Load() != 1 {
		t.Errorf("handler a: expected 1 call, got %d", countA.Load())
	}
	if countB.Load() != 1 {
		t.Errorf("handler b: expected 1 call, got %d", countB.Load())
	}
}

func TestEmit_PassesDiffToHandler(t *testing.T) {
	s := eventsink.New()
	var gotOpened, gotClosed []snapshot.Entry
	s.Register("recv", func(o, c []snapshot.Entry) {
		gotOpened = o
		gotClosed = c
	})

	opened := []snapshot.Entry{makeEntry(22, "tcp")}
	closed := []snapshot.Entry{makeEntry(8080, "tcp")}
	s.Emit(opened, closed)

	if len(gotOpened) != 1 || gotOpened[0].Port != 22 {
		t.Errorf("unexpected opened: %v", gotOpened)
	}
	if len(gotClosed) != 1 || gotClosed[0].Port != 8080 {
		t.Errorf("unexpected closed: %v", gotClosed)
	}
}

func TestEmit_EmptyDiffNoError(t *testing.T) {
	s := eventsink.New()
	s.Register("noop", func(_, _ []snapshot.Entry) {})
	s.Emit(nil, nil) // must not panic
}
