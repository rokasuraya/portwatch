package notifier_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"portwatch/internal/notifier"
)

func makeEvent(opened, closed []string) notifier.Event {
	return notifier.Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Opened:    opened,
		Closed:    closed,
	}
}

func TestDispatch_WritesToStdout(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf, "")

	err := n.Dispatch(makeEvent([]string{"tcp:8080"}, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "tcp:8080") {
		t.Errorf("expected output to contain port, got: %s", buf.String())
	}
}

func TestDispatch_SendsWebhook(t *testing.T) {
	var received notifier.Event
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	var buf bytes.Buffer
	n := notifier.New(&buf, ts.URL)

	evt := makeEvent([]string{"tcp:443"}, []string{"tcp:80"})
	if err := n.Dispatch(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.Opened) == 0 || received.Opened[0] != "tcp:443" {
		t.Errorf("webhook payload mismatch: %+v", received)
	}
}

func TestDispatch_WebhookError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	var buf bytes.Buffer
	n := notifier.New(&buf, ts.URL)

	err := n.Dispatch(makeEvent(nil, nil))
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDispatch_NoWebhook_NoError(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf, "")

	if err := n.Dispatch(makeEvent(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
