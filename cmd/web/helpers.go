package main

import (
	"net/http"
	"runtime/debug"
	"snippetbox/internal/config"
)

func serverError(app *config.Application, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
}

func clientError(status int, w http.ResponseWriter) {
	http.Error(w, http.StatusText(status), status)
}
