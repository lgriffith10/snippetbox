package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", getSnippet)
	mux.HandleFunc("GET /snippet/create", getSnippetCreationForm)
	mux.HandleFunc("POST /snippet/create", createSnippet)

	log.Println("Listening on :8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
