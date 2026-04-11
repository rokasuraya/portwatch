// Package webhook implements a lightweight HTTP webhook client used by
// portwatch to POST JSON payloads describing port-change events to an
// operator-configured endpoint.
//
// Usage:
//
//	client := webhook.New("https://example.com/hook", 5*time.Second)
//	err := client.Send(ctx, webhook.Payload{
//		Timestamp: time.Now().UTC().Format(time.RFC3339),
//		Opened:    []string{"tcp:8080"},
//		Closed:    []string{},
//	})
package webhook
