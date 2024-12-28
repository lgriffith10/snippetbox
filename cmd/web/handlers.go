package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/config"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func Home(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.Snippets.Latest()
		if err != nil {
			serverError(app, r, err)
			clientError(http.StatusInternalServerError, w)
		}

		data := newTemplateData(r)
		data.Snippets = snippets

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

		data := newTemplateData(r)
		data.Snippet = snippet

		render(app, w, r, http.StatusOK, "view.html", data)
	}
}

func GetSnippetCreationForm(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := newTemplateData(r)
		data.Form = CreateSnippetForm{
			Expires: 365,
		}

		render(app, w, r, http.StatusOK, "create.html", data)
	}
}

type CreateSnippetForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string

	validator.Validator
}

func CreateSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			clientError(http.StatusBadRequest, w)
			return
		}

		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			clientError(http.StatusBadRequest, w)
			return
		}

		form := CreateSnippetForm{
			Title:       r.PostForm.Get("title"),
			Content:     r.PostForm.Get("content"),
			Expires:     expires,
			FieldErrors: map[string]string{},
		}

		permittedYears := []int{
			1, 7, 365,
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedValue(form.Expires, permittedYears), "expires", "This field cannot be blank")

		if !form.Valid() {
			data := newTemplateData(r)
			data.Form = form
			fmt.Println(data.Form)
			render(app, w, r, http.StatusUnprocessableEntity, "create.html", data)
			return
		}

		id, err := app.Snippets.Insert(form.Title, form.Content, expires)
		if err != nil {
			serverError(app, r, err)
			clientError(http.StatusInternalServerError, w)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}
