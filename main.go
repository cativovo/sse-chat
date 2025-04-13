package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	logger := slog.Default()

	provider := NewProvider()

	mux := http.NewServeMux()

	mux.Handle("GET /", homePage(logger))

	mux.Handle("GET /sse", handleSSE(logger, provider))
	mux.Handle("POST /message", handleSendMessage(logger, provider))

	port := 4000

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	logger.Info("Server is listening and serving", "port", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
