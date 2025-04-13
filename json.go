package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, logger *slog.Logger, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		logger.Error("Failed to encode the response", "error", err)
		fmt.Fprint(w, `{"error":"Something went wrong"}`)
	}
}

func bind(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&dst); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	return nil
}
