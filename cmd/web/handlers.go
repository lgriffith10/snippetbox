package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.serverError(r, err)
		app.clientError(http.StatusInternalServerError, w)
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *Application) GetSnippet(w http.ResponseWriter, r *http.Request) {
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
			app.serverError(r, err)
		}
		return
	}

	flash := app.SessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.Flash = flash

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *Application) GetSnippetCreationForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = CreateSnippetForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.html", data)
}

type CreateSnippetForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`

	validator.Validator `form:"-"`
}

func (app *Application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var form CreateSnippetForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(http.StatusBadRequest, w)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(r, err)
		app.clientError(http.StatusInternalServerError, w)
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type UserSignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`

	validator.Validator `form:"-"`
}

func (app *Application) GetUserSignUpForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *Application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var form UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(http.StatusBadRequest, w)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "Name", "Name cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "Email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "Email", "This must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "Password", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "Password", "Password must be at least 8 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("Email", "Email is already taken")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
			return
		} else {
			app.serverError(r, err)
		}

		return
	}

	app.SessionManager.Put(r.Context(), "flash", "User successfully created")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) GetUserLoginForm(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Display a form for logging user")
}

func (app *Application) LoginUser(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Login user")
}

func (app *Application) LogoutUser(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "Logout user")
}
