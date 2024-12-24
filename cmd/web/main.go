package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"snippetbox/env"
)

func main() {
	mode := flag.String("mode", "dev", "run mode")
	flag.Parse()

	env.SetEnvVariables(*mode)
	port := os.Getenv("GO_PORT")

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", getSnippet)
	mux.HandleFunc("GET /snippet/create", getSnippetCreationForm)
	mux.HandleFunc("POST /snippet/create", createSnippet)

	log.Printf("Listening on localhost%s", port)

	err := http.ListenAndServe(port, mux)
	log.Fatal(err)
}
