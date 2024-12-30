package main

import (
	"github.com/justinas/alice"
	"net/http"
	"snippetbox/ui"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServerFS(ui.Files))
	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(app.SessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.Home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.GetSnippet))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.GetUserSignUpForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.RegisterUser))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.GetUserLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.LoginUser))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.GetSnippetCreationForm))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.CreateSnippet))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.LogoutUser))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
