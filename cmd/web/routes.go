package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.Home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.GetSnippet))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.GetSnippetCreationForm))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.CreateSnippet))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.GetUserSignUpForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.RegisterUser))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.GetUserLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.LoginUser))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.LogoutUser))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
