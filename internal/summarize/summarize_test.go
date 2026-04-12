package summarize_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"portwatch/internal/summarize"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	s := summarize.New(time.Minute, nil)
	if s == nil {
		t.Fatal("expected non-nil Summarizer")
	}
}

func TestRecord_AccumulatesCounts(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.New(time.Minute, &buf)

	s.Record(3, 1, 1)
	s.Record(2, 0, 0)

	// Force a flush via a very short interval and a cancelled context.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Run with a long ticker so only the ctx.Done flush fires.
	s2 := summarize.New(time.Hour, &buf)
	s2.Record(3, 1, 1)
	s2.Record(2, 0, 0)
	s2.Run(ctx)

	out := buf.String()
	if !strings.Contains(out, "scans=2") {
		t.Errorf("expected scans=2 in output, got: %s", out)
	}
	if !strings.Contains(out, "opened=3") {
		t.Errorf("expected opened=3 in output, got: %s", out)
	}
	if !strings.Contains(out, "closed=1") {
		t.Errorf("expected closed=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "alerts=1") {
		t.Errorf("expected alerts=1 in output, got: %s", out)
	}
}

func TestRun_FlushesOnCancel(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.New(time.Hour, &buf)
	s.Record(1, 0, 0)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after context cancellation")
	}

	if !strings.Contains(buf.String(), "scans=1") {
		t.Errorf("expected flush on cancel, got: %s", buf.String())
	}
}

func TestRun_FlushesOnTick(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.New(30*time.Millisecond, &buf)
	s.Record(0, 2, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	if !strings.Contains(buf.String(), "closed=2") {
		t.Errorf("expected closed=2 after tick flush, got: %s", buf.String())
	}
}
