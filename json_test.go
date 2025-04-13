package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	testCases := []struct {
		name   string
		status int
		v      any
		want   string
	}{
		{
			name:   "200",
			status: http.StatusOK,
			v:      map[string]any{"test": "test"},
			want:   `{"test":"test"}`,
		},
		{
			name:   "400",
			status: http.StatusBadRequest,
			v:      map[string]any{"test": "test"},
			want:   `{"test":"test"}`,
		},
		{
			name:   "500",
			status: http.StatusInternalServerError,
			v:      func() {},
			want:   `{"error":"Something went wrong"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			JSON(w, logger, tc.status, tc.v)

			gotContentType := w.Header().Get("Content-Type")
			wantContentType := "application/json"
			if gotContentType != wantContentType {
				t.Errorf("want %q, got %q", wantContentType, gotContentType)
			}

			// Encode adds \n
			got := strings.TrimSpace(w.Body.String())
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}

			if w.Code != tc.status {
				t.Errorf("want %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func TestBind(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`{"test":"test"}`)
	r := httptest.NewRequest(http.MethodPost, "/", &buf)

	var dst struct {
		Test string `json:"test"`
	}
	if err := bind(r, &dst); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := "test"
	if dst.Test != want {
		t.Errorf("want %q, got %q", want, dst.Test)
	}
}
