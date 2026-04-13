package fanout_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/fanout"
)

func TestNew_ReturnsFanout(t *testing.T) {
	f := fanout.New[string]()
	if f == nil {
		t.Fatal("expected non-nil Fanout")
	}
}

func TestSubscribe_IncreasesLen(t *testing.T) {
	f := fanout.New[int]()
	f.Subscribe("a", func(v int) {})
	if f.Len() != 1 {
		t.Fatalf("expected 1 subscriber, got %d", f.Len())
	}
}

func TestSubscribe_ReplacesExistingName(t *testing.T) {
	f := fanout.New[int]()
	f.Subscribe("a", func(v int) {})
	f.Subscribe("a", func(v int) {})
	if f.Len() != 1 {
		t.Fatalf("expected 1 subscriber after replace, got %d", f.Len())
	}
}

func TestUnsubscribe_RemovesHandler(t *testing.T) {
	f := fanout.New[int]()
	f.Subscribe("a", func(v int) {})
	f.Unsubscribe("a")
	if f.Len() != 0 {
		t.Fatalf("expected 0 subscribers, got %d", f.Len())
	}
}

func TestUnsubscribe_NoopForUnknownName(t *testing.T) {
	f := fanout.New[int]()
	f.Unsubscribe("missing") // should not panic
}

func TestPublish_DeliverstToAllSubscribers(t *testing.T) {
	f := fanout.New[string]()
	var mu sync.Mutex
	got := map[string]string{}

	f.Subscribe("a", func(v string) {
		mu.Lock(); defer mu.Unlock()
		got["a"] = v
	})
	f.Subscribe("b", func(v string) {
		mu.Lock(); defer mu.Unlock()
		got["b"] = v
	})

	f.Publish("hello")

	// allow goroutines to complete
	time.Sleep(20 * time.Millisecond)

	mu.Lock(); defer mu.Unlock()
	if got["a"] != "hello" || got["b"] != "hello" {
		t.Fatalf("expected both subscribers to receive 'hello', got %v", got)
	}
}

func TestPublish_NoSubscribers_NoError(t *testing.T) {
	f := fanout.New[int]()
	f.Publish(42) // should not panic
}
