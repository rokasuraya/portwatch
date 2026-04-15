// Package envelope wraps a scan diff with metadata for downstream consumers.
package envelope

import (
	"time"

	"github.com/yourorg/portwatch/internal/snapshot"
)

// Envelope carries a snapshot diff alongside contextual metadata.
type Envelope struct {
	// ID is a unique identifier for this event, typically a short hex string.
	ID string

	// CreatedAt is the wall-clock time the envelope was created.
	CreatedAt time.Time

	// Opened holds entries for ports that became open since the last scan.
	Opened []snapshot.Entry

	// Closed holds entries for ports that became closed since the last scan.
	Closed []snapshot.Entry

	// ScanDuration is how long the underlying scan took.
	ScanDuration time.Duration

	// Labels carries arbitrary key/value metadata attached at creation time.
	Labels map[string]string
}

// New creates an Envelope from a diff and scan duration.
// labels may be nil; a non-nil copy is always stored.
func New(id string, opened, closed []snapshot.Entry, dur time.Duration, labels map[string]string) *Envelope {
	lbl := make(map[string]string, len(labels))
	for k, v := range labels {
		lbl[k] = v
	}
	return &Envelope{
		ID:           id,
		CreatedAt:    time.Now().UTC(),
		Opened:       copyEntries(opened),
		Closed:       copyEntries(closed),
		ScanDuration: dur,
		Labels:       lbl,
	}
}

// IsEmpty reports whether the envelope carries no diff.
func (e *Envelope) IsEmpty() bool {
	return len(e.Opened) == 0 && len(e.Closed) == 0
}

// AddLabel attaches a key/value pair to the envelope.
func (e *Envelope) AddLabel(key, value string) {
	if e.Labels == nil {
		e.Labels = make(map[string]string)
	}
	e.Labels[key] = value
}

func copyEntries(src []snapshot.Entry) []snapshot.Entry {
	if len(src) == 0 {
		return []snapshot.Entry{}
	}
	dst := make([]snapshot.Entry, len(src))
	copy(dst, src)
	return dst
}
