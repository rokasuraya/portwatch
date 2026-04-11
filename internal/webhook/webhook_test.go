package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/webhook"
)

func TestNew_DefaultTimeout(t *testing.T) {
	c := webhook.New("http://localhost", 0)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSend_PostsJSON(t *testing.T) {
	var received webhook.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := webhook.New(ts.URL, time.Second)
	p := webhook.Payload{
		Timestamp: "2024-01-01T00:00:00Z",
		Opened:    []string{"tcp:9090"},
		Closed:    []string{"tcp:8080"},
	}
	if err := client.Send(context.Background(), p); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if len(received.Opened) != 1 || received.Opened[0] != "tcp:9090" {
		t.Errorf("unexpected opened: %v", received.Opened)
	}
	if len(received.Closed) != 1 || received.Closed[0] != "tcp:8080" {
		t.Errorf("unexpected closed: %v", received.Closed)
	}
}

func TestSend_NonOKStatusReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := webhook.New(ts.URL, time.Second)
	err := client.Send(context.Background(), webhook.Payload{})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSend_InvalidURLReturnsError(t *testing.T) {
	client := webhook.New("http://127.0.0.1:0", 100*time.Millisecond)
	err := client.Send(context.Background(), webhook.Payload{})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestSend_CancelledContextReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := webhook.New(ts.URL, time.Second)
	if err := client.Send(ctx, webhook.Payload{}); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
