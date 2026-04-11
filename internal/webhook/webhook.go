// Package webhook provides HTTP webhook dispatch for port change events.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultTimeout = 5 * time.Second Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp string   `json:"timestamp"`
	Opened    []string `json:"opened"`
	Closed    []string `json:"closed"`
}

// Client sends webhook notifications to a configured URL.
type Client struct {
	url     string
	hc      *http.Client
	timeout time.Duration
}

// New returns a new Client targeting the given URL.
func New(url string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Client{
		url:     url,
		hc:      &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

// Send marshals p and POSTs it to the configured URL.
func (c *Client) Send(ctx context.Context, p Payload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
