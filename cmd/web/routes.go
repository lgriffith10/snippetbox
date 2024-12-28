package main

import (
	"net/http"
	"snippetbox/internal/config"
)

func routes(app *config.Application) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /{$}", Home(app))
	mux.HandleFunc("GET /snippet/view/{id}", GetSnippet(app))
	mux.HandleFunc("GET /snippet/create", GetSnippetCreationForm(app))
	mux.HandleFunc("POST /snippet/create", CreateSnippet(app))

	return recoverPanic(app, logRequest(app, commonHeaders(mux)))
}
