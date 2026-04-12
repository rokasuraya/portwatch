package pipeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/snapshot"
)

func makeStore(t *testing.T) *snapshot.Store {
	t.Helper()
	store, err := snapshot.NewStore("")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func goodScan(_ context.Context) ([]snapshot.Entry, error) {
	return []snapshot.Entry{{Port: 80, Proto: "tcp"}}, nil
}

func badScan(_ context.Context) ([]snapshot.Entry, error) {
	return nil, errors.New("scan failed")
}

func TestNew_ReturnsPipeline(t *testing.T) {
	p := pipeline.New(goodScan, makeStore(t))
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestTick_ReturnsDuration(t *testing.T) {
	p := pipeline.New(goodScan, makeStore(t))
	dur, err := p.Tick(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dur < 0 {
		t.Fatalf("expected non-negative duration, got %v", dur)
	}
}

func TestTick_ScanErrorPropagates(t *testing.T) {
	p := pipeline.New(badScan, makeStore(t))
	_, err := p.Tick(context.Background())
	if err == nil {
		t.Fatal("expected error from failing scan")
	}
}

func TestTick_InvokesStages(t *testing.T) {
	var called int
	stage := func(_ context.Context, opened, closed []snapshot.Entry) error {
		called++
		return nil
	}

	p := pipeline.New(goodScan, makeStore(t), stage, stage)
	_, err := p.Tick(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 2 {
		t.Fatalf("expected 2 stage calls, got %d", called)
	}
}

func TestTick_StageErrorReturned(t *testing.T) {
	sentinel := errors.New("stage error")
	stage := func(_ context.Context, _, _ []snapshot.Entry) error { return sentinel }

	p := pipeline.New(goodScan, makeStore(t), stage)
	_, err := p.Tick(context.Background())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestTick_AllStagesRunDespiteError(t *testing.T) {
	var secondCalled bool
	errStage := func(_ context.Context, _, _ []snapshot.Entry) error {
		return errors.New("first stage error")
	}
	okStage := func(_ context.Context, _, _ []snapshot.Entry) error {
		secondCalled = true
		return nil
	}

	p := pipeline.New(goodScan, makeStore(t), errStage, okStage)
	_, _ = p.Tick(context.Background())
	if !secondCalled {
		t.Fatal("expected second stage to run even after first error")
	}
}

func TestTick_DiffDetectedOnSecondTick(t *testing.T) {
	var lastOpened []snapshot.Entry
	stage := func(_ context.Context, opened, _ []snapshot.Entry) error {
		lastOpened = opened
		return nil
	}

	call := 0
	scan := func(_ context.Context) ([]snapshot.Entry, error) {
		call++
		if call == 1 {
			return []snapshot.Entry{{Port: 80, Proto: "tcp"}}, nil
		}
		return []snapshot.Entry{{Port: 80, Proto: "tcp"}, {Port: 443, Proto: "tcp"}}, nil
	}

	p := pipeline.New(scan, makeStore(t), stage)
	if _, err := p.Tick(context.Background()); err != nil {
		t.Fatalf("first tick: %v", err)
	}
	if _, err := p.Tick(context.Background()); err != nil {
		t.Fatalf("second tick: %v", err)
	}
	if len(lastOpened) != 1 || lastOpened[0].Port != 443 {
		t.Fatalf("expected port 443 opened, got %v", lastOpened)
	}
	_ = time.Second // keep time import used
}
