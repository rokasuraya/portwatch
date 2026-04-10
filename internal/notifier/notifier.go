// Package notifier provides pluggable notification backends for portwatch.
// Supported backends: stdout, webhook.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Event represents a port-change notification payload.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []string  `json:"opened"`
	Closed    []string  `json:"closed"`
}

// Notifier dispatches events to one or more backends.
type Notifier struct {
	backends []backend
}

type backend interface {
	send(e Event) error
}

// New creates a Notifier. If webhookURL is non-empty a webhook backend is
// added; a stdout backend is always included as fallback.
func New(w io.Writer, webhookURL string) *Notifier {
	n := &Notifier{}
	n.backends = append(n.backends, &stdoutBackend{w: w})
	if webhookURL != "" {
		n.backends = append(n.backends, &webhookBackend{
			url:    webhookURL,
			client: &http.Client{Timeout: 5 * time.Second},
		})
	}
	return n
}

// Dispatch sends the event to every registered backend.
// It returns the first error encountered, if any.
func (n *Notifier) Dispatch(e Event) error {
	for _, b := range n.backends {
		if err := b.send(e); err != nil {
			return err
		}
	}
	return nil
}

// stdoutBackend writes a human-readable summary to an io.Writer.
type stdoutBackend struct{ w io.Writer }

func (s *stdoutBackend) send(e Event) error {
	_, err := fmt.Fprintf(s.w, "[%s] opened=%v closed=%v\n",
		e.Timestamp.Format(time.RFC3339), e.Opened, e.Closed)
	return err
}

// webhookBackend POSTs a JSON payload to a URL.
type webhookBackend struct {
	url    string
	client *http.Client
}

func (wb *webhookBackend) send(e Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("notifier: marshal: %w", err)
	}
	resp, err := wb.client.Post(wb.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned %d", resp.StatusCode)
	}
	return nil
}
