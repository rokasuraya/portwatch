// Package metrics tracks runtime statistics for the portwatch daemon,
// such as scan counts, alert counts, and last scan duration.
package metrics

import (
	"sync"
	"time"
)

// Metrics holds cumulative runtime statistics.
type Metrics struct {
	mu            sync.RWMutex
	ScanCount     int64
	AlertCount    int64
	LastScanAt    time.Time
	LastScanDur   time.Duration
	OpenedTotal   int64
	ClosedTotal   int64
}

// New returns a zero-value Metrics instance.
func New() *Metrics {
	return &Metrics{}
}

// RecordScan updates scan-related counters after each tick.
func (m *Metrics) RecordScan(dur time.Duration, opened, closed int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ScanCount++
	m.LastScanAt = time.Now()
	m.LastScanDur = dur
	m.OpenedTotal += int64(opened)
	m.ClosedTotal += int64(closed)

	if opened > 0 || closed > 0 {
		m.AlertCount++
	}
}

// Snapshot returns a point-in-time copy of the current metrics.
func (m *Metrics) Snapshot() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Metrics{
		ScanCount:   m.ScanCount,
		AlertCount:  m.AlertCount,
		LastScanAt:  m.LastScanAt,
		LastScanDur: m.LastScanDur,
		OpenedTotal: m.OpenedTotal,
		ClosedTotal: m.ClosedTotal,
	}
}

// Reset zeroes all counters. Useful for testing.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	*m = Metrics{}
}
