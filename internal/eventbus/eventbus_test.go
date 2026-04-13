package eventbus_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"portwatch/internal/eventbus"
)

func TestNew_ReturnsEmptyBus(t *testing.T) {
	b := eventbus.New[string]()
	if b.Len() != 0 {
		t.Fatalf("expected 0 handlers, got %d", b.Len())
	}
}

func TestSubscribe_IncreasesLen(t *testing.T) {
	b := eventbus.New[int]()
	b.Subscribe("a", func(int) {})
	b.Subscribe("b", func(int) {})
	if b.Len() != 2 {
		t.Fatalf("expected 2 handlers, got %d", b.Len())
	}
}

func TestSubscribe_ReplacesExistingName(t *testing.T) {
	b := eventbus.New[int]()
	b.Subscribe("a", func(int) {})
	b.Subscribe("a", func(int) {})
	if b.Len() != 1 {
		t.Fatalf("expected 1 handler after duplicate subscribe, got %d", b.Len())
	}
}

func TestUnsubscribe_RemovesHandler(t *testing.T) {
	b := eventbus.New[string]()
	b.Subscribe("x", func(string) {})
	b.Unsubscribe("x")
	if b.Len() != 0 {
		t.Fatalf("expected 0 handlers after unsubscribe, got %d", b.Len())
	}
}

func TestUnsubscribe_NoopForUnknownName(t *testing.T) {
	b := eventbus.New[string]()
	b.Unsubscribe("nonexistent") // must not panic
}

func TestPublish_DeliveresToAllHandlers(t *testing.T) {
	b := eventbus.New[int]()
	var got []int
	var mu sync.Mutex
	collect := func(v int) {
		mu.Lock()
		got = append(got, v)
		mu.Unlock()
	}
	b.Subscribe("h1", collect)
	b.Subscribe("h2", collect)
	b.Publish(42)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 deliveries, got %d", len(got))
	}
	for _, v := range got {
		if v != 42 {
			t.Errorf("expected 42, got %d", v)
		}
	}
}

func TestPublish_ConcurrentSafe(t *testing.T) {
	b := eventbus.New[int]()
	var count atomic.Int64
	b.Subscribe("counter", func(int) { count.Add(1) })
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Publish(1)
		}()
	}
	wg.Wait()
	if count.Load() != 100 {
		t.Fatalf("expected 100 deliveries, got %d", count.Load())
	}
}
