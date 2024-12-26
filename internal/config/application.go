package config

import (
	"database/sql"
	"html/template"
	"log/slog"
	"snippetbox/internal"
	"snippetbox/internal/models"
)

type Application struct {
	Logger        *slog.Logger
	Snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
}

func NewApplication(db *sql.DB, templateCache map[string]*template.Template) *Application {
	return &Application{
		Logger:        internal.NewLogger(),
		Snippets:      &models.SnippetModel{DB: db},
		TemplateCache: templateCache,
	}
}
