package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/debounce"
)

func TestNew_ReturnsDebouncerWithWait(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	if d == nil {
		t.Fatal("expected non-nil Debouncer")
	}
}

func TestTrigger_FiresAfterWait(t *testing.T) {
	var called int32
	d := debounce.New(30 * time.Millisecond)

	d.Trigger(func() { atomic.StoreInt32(&called, 1) })

	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Error("expected fn to be called after wait period")
	}
}

func TestTrigger_ResetsTimerOnRapidCalls(t *testing.T) {
	var count int32
	d := debounce.New(50 * time.Millisecond)

	for i := 0; i < 5; i++ {
		d.Trigger(func() { atomic.AddInt32(&count, 1) })
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)
	if c := atomic.LoadInt32(&count); c != 1 {
		t.Errorf("expected fn called exactly once, got %d", c)
	}
}

func TestFlush_CancelsPendingTimer(t *testing.T) {
	var called int32
	d := debounce.New(200 * time.Millisecond)

	d.Trigger(func() { atomic.StoreInt32(&called, 1) })
	flushed := d.Flush()

	if !flushed {
		t.Error("expected Flush to return true when timer was pending")
	}
	time.Sleep(250 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Error("expected fn NOT to be called after Flush")
	}
}

func TestFlush_ReturnsFalseWhenNoPending(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	if d.Flush() {
		t.Error("expected Flush to return false when no timer is pending")
	}
}

func TestReset_CancelsPendingTimer(t *testing.T) {
	var called int32
	d := debounce.New(100 * time.Millisecond)

	d.Trigger(func() { atomic.StoreInt32(&called, 1) })
	d.Reset()

	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&called) != 0 {
		t.Error("expected fn NOT to be called after Reset")
	}
}
