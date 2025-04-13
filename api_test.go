package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// inconsistent results when run with go test -race -count 10
func TestHandleSSE(t *testing.T) {
	p := NewProvider()

	ctx, cancel := context.WithCancel(context.Background())
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	go func() {
		p.Publish(Event{
			EventType: EventTypeMessage,
			Data:      1,
		})
		p.Publish(Event{
			EventType: EventTypeMessage,
			Data:      2,
		})

		// make sure the handler send the events before cancelling the request
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	handleSSE(logger, p).ServeHTTP(w, r)

	b, err := io.ReadAll(w.Result().Body)
	if err != nil && err != context.Canceled {
		t.Fatal(err)
	}

	wantContentType := "text/event-stream"
	gotContentType := w.Header().Get("Content-Type")
	if gotContentType != wantContentType {
		t.Errorf("want %q, got %q", wantContentType, gotContentType)
	}

	wantEvents := "event: connect\ndata: null\n\nevent: message\ndata: 1\n\nevent: message\ndata: 2\n\n"
	gotEvents := string(b)
	if gotEvents != wantEvents {
		t.Errorf("want %q, got %q", wantEvents, gotEvents)
	}
}

func TestHandleSendMessage(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		wantStatus   int
		wantResponse any
	}{
		{
			name:       "ok",
			body:       `{"message":"test"}`,
			wantStatus: http.StatusNoContent,
		},
		{
			name:         "empty message",
			body:         `{"message":""}`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: map[string]any{"error": "message is empty"},
		},
		{
			name:         "failed to parse the body",
			body:         `{"message":}`,
			wantStatus:   http.StatusBadRequest,
			wantResponse: map[string]any{"error": "Failed to parse the body"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			buf.WriteString(tc.body)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", &buf)
			p := NewProvider()
			handleSendMessage(logger, p).ServeHTTP(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("want %d, got %d", tc.wantStatus, w.Code)
			}

			var wantResponse string
			if tc.wantResponse != nil {
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(tc.wantResponse); err != nil {
					t.Fatal(err)
				}
				wantResponse = buf.String()
			}

			gotResponse := w.Body.String()
			if gotResponse != wantResponse {
				t.Errorf("want %q, got %q", wantResponse, w.Body.String())
			}
		})
	}
}
