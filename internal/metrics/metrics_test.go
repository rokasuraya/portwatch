package metrics_test

import (
	"testing"
	"time"

	"portwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	m := metrics.New()
	s := m.Snapshot()

	if s.ScanCount != 0 || s.AlertCount != 0 || s.OpenedTotal != 0 || s.ClosedTotal != 0 {
		t.Errorf("expected all zero values, got %+v", s)
	}
}

func TestRecordScan_IncrementsScanCount(t *testing.T) {
	m := metrics.New()
	m.RecordScan(10*time.Millisecond, 0, 0)
	m.RecordScan(20*time.Millisecond, 0, 0)

	if got := m.Snapshot().ScanCount; got != 2 {
		t.Errorf("expected ScanCount=2, got %d", got)
	}
}

func TestRecordScan_TracksOpenedClosed(t *testing.T) {
	m := metrics.New()
	m.RecordScan(5*time.Millisecond, 3, 1)

	s := m.Snapshot()
	if s.OpenedTotal != 3 {
		t.Errorf("expected OpenedTotal=3, got %d", s.OpenedTotal)
	}
	if s.ClosedTotal != 1 {
		t.Errorf("expected ClosedTotal=1, got %d", s.ClosedTotal)
	}
}

func TestRecordScan_IncrementsAlertCountOnDiff(t *testing.T) {
	m := metrics.New()
	m.RecordScan(5*time.Millisecond, 0, 0) // no diff — no alert
	m.RecordScan(5*time.Millisecond, 1, 0) // diff — alert
	m.RecordScan(5*time.Millisecond, 0, 2) // diff — alert

	if got := m.Snapshot().AlertCount; got != 2 {
		t.Errorf("expected AlertCount=2, got %d", got)
	}
}

func TestRecordScan_StoresLastScanDuration(t *testing.T) {
	m := metrics.New()
	want := 42 * time.Millisecond
	m.RecordScan(want, 0, 0)

	if got := m.Snapshot().LastScanDur; got != want {
		t.Errorf("expected LastScanDur=%v, got %v", want, got)
	}
}

func TestSnapshot_IsIndependent(t *testing.T) {
	m := metrics.New()
	m.RecordScan(1*time.Millisecond, 1, 0)
	s1 := m.Snapshot()

	m.RecordScan(1*time.Millisecond, 1, 0)
	s2 := m.Snapshot()

	if s1.ScanCount == s2.ScanCount {
		t.Error("expected snapshots to be independent copies")
	}
}

func TestReset_ZeroesCounters(t *testing.T) {
	m := metrics.New()
	m.RecordScan(5*time.Millisecond, 2, 3)
	m.Reset()

	s := m.Snapshot()
	if s.ScanCount != 0 || s.OpenedTotal != 0 || s.ClosedTotal != 0 {
		t.Errorf("expected zeroed metrics after Reset, got %+v", s)
	}
}
