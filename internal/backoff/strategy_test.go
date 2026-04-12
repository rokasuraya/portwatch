package backoff

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errFake = errors.New("fake error")

func TestRetry_SucceedsOnFirstAttempt(t *testing.T) {
	b := New(time.Millisecond, time.Second, 2.0)
	calls := 0
	err := Retry(context.Background(), b, 3, func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetry_RetriesOnFailure(t *testing.T) {
	b := New(time.Millisecond, time.Second, 2.0)
	calls := 0
	err := Retry(context.Background(), b, 3, func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errFake
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after 3 attempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRetry_ExhaustsMaxAttempts(t *testing.T) {
	b := New(time.Millisecond, time.Second, 2.0)
	calls := 0
	err := Retry(context.Background(), b, 4, func(_ context.Context) error {
		calls++
		return errFake
	})
	if !errors.Is(err, errFake) {
		t.Fatalf("expected errFake, got %v", err)
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestRetry_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	b := New(time.Millisecond, time.Second, 2.0)
	err := Retry(ctx, b, 5, func(_ context.Context) error {
		return errFake
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
