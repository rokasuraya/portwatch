package gwatcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"portwatch/internal/gwatcher"
	"portwatch/internal/scanner"
	"portwatch/internal/snapshot"
)

func TestRunner_CancelsCleanly(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, nil)
	r := gwatcher.NewRunner(w, func() (*snapshot.Snapshot, error) {
		return snapshot.New(nil), nil
	}, func([]gwatcher.Event) {}, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := r.Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestRunner_InvokesOnEvent(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, nil)

	called := make(chan struct{}, 1)
	snapshotFn := func() (*snapshot.Snapshot, error) {
		return snapshot.New([]scanner.Entry{{Port: 80, Proto: "tcp", Open: true}}), nil
	}
	onEvent := func(evs []gwatcher.Event) {
		if len(evs) > 0 {
			select {
			case called <- struct{}{}:
			default:
			}
		}
	}

	r := gwatcher.NewRunner(w, snapshotFn, onEvent, 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	go r.Run(ctx) //nolint:errcheck

	select {
	case <-called:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("onEvent was never called")
	}
}

func TestRunner_SkipsOnSnapshotError(t *testing.T) {
	m := buildMatcher(t)
	w := gwatcher.New(m, nil)

	snapshotFn := func() (*snapshot.Snapshot, error) {
		return nil, errors.New("scan failed")
	}
	events := 0
	r := gwatcher.NewRunner(w, snapshotFn, func([]gwatcher.Event) { events++ }, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	r.Run(ctx) //nolint:errcheck

	if events != 0 {
		t.Fatalf("expected 0 events on error, got %d", events)
	}
}
