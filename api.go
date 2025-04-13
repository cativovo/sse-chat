package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func send(rc *http.ResponseController, w io.Writer, e Event) error {
	b, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.EventType, string(b))

	return rc.Flush()
}

func handleSSE(logger *slog.Logger, p Provider) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			c, cancel := p.Subscribe()
			defer cancel()
			rc := http.NewResponseController(w)

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")

			e := Event{EventType: EventTypeConnect}
			if err := send(rc, w, e); err != nil {
				logger.Error("Failed to send the event", "error", err, "event", e)
				JSON(w, logger, http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
				return
			}

			for {
				select {
				case <-r.Context().Done():
					logger.Info("Client has disconnected")
					return
				case e, ok := <-c:
					if !ok {
						logger.Info("Channel is closed")
						return
					}

					logger.Info("Sending event", "event", e)
					if err := send(rc, w, e); err != nil {
						logger.Error("Failed to send the event", "error", err, "event", e)
						JSON(w, logger, http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
						return
					}
				}
			}
		})
}

func handleSendMessage(logger *slog.Logger, p Provider) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var dst struct {
				Message string `json:"message"`
			}
			if err := bind(r, &dst); err != nil {
				JSON(w, logger, http.StatusBadRequest, map[string]any{"error": "Failed to parse the body"})
				return
			}

			if dst.Message == "" {
				JSON(w, logger, http.StatusBadRequest, map[string]any{"error": "message is empty"})
				return
			}

			p.Publish(Event{
				EventType: EventTypeMessage,
				Data:      dst,
			})

			w.WriteHeader(http.StatusNoContent)
		})
}
