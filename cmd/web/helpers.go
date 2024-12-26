package main

import (
	"fmt"
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

func render(
	app *config.Application,
	w http.ResponseWriter,
	r *http.Request,
	status int,
	page string,
	data templateData,
) {
	ts, ok := app.TemplateCache[page]

	if !ok {
		err := fmt.Errorf("template %s does not exists", page)
		serverError(app, r, err)
		clientError(http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		serverError(app, r, err)
		clientError(http.StatusInternalServerError, w)
		return
	}
}
