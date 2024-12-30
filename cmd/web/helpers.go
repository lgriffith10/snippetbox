package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

}

func (app *Application) clientError(status int, w http.ResponseWriter) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) render(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	page string,
	data templateData,
) {
	ts, ok := app.TemplateCache[page]

	if !ok {
		err := fmt.Errorf("template %s does not exists", page)
		app.serverError(r, err)
		app.clientError(http.StatusInternalServerError, w)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(r, err)
		app.clientError(http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(status)

	_, _ = buf.WriteTo(w)
}

func (app *Application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
