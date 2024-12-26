package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/config"
	"snippetbox/internal/models"
	"strconv"
)

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "GO")

		snippets, err := app.Snippets.Latest()
		if err != nil {
			serverError(app, r, err)
			clientError(http.StatusInternalServerError, w)
		}

		data := templateData{
			Snippets: snippets,
		}

		render(app, w, r, http.StatusOK, "home.html", data)
	}
}

func GetSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				serverError(app, r, err)
			}
			return
		}

		data := templateData{
			Snippet: snippet,
		}

		render(app, w, r, http.StatusOK, "view.html", data)
	}
}

func GetSnippetCreationForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Display a form for creating a new snippet"))
	}
}

func CreateSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := "0 snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			serverError(app, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}
