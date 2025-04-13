package main

import (
	"html/template"
	"log/slog"
	"net/http"
)

var indexTmpl = template.Must(template.ParseFiles("./index.html"))

func homePage(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := indexTmpl.Execute(w, nil); err != nil {
				logger.Error("Failed to execute indexTmpl", "error", err)
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
		})
}
