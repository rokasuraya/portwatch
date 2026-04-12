// Package audit provides a structured audit log for port change events,
// recording who observed a change, when, and what the diff contained.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Note      string    `json:"note,omitempty"`
}

// Audit writes structured audit entries to an output destination.
type Audit struct {
	mu  sync.Mutex
	out io.Writer
}

// New returns an Audit that writes to out.
// If out is nil, os.Stderr is used.
func New(out io.Writer) *Audit {
	if out == nil {
		out = os.Stderr
	}
	return &Audit{out: out}
}

// Log writes a single audit entry as a JSON line.
func (a *Audit) Log(event, protocol string, port int, note string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Port:      port,
		Protocol:  protocol,
		Note:      note,
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(a.out, "%s\n", b)
	return err
}

// LogOpened is a convenience wrapper for an opened-port event.
func (a *Audit) LogOpened(protocol string, port int) error {
	return a.Log("opened", protocol, port, "")
}

// LogClosed is a convenience wrapper for a closed-port event.
func (a *Audit) LogClosed(protocol string, port int) error {
	return a.Log("closed", protocol, port, "")
}
