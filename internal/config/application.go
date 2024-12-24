package config

import (
	"database/sql"
	"log/slog"
	"snippetbox/internal"
	"snippetbox/internal/models"
)

type Application struct {
	Logger   *slog.Logger
	Snippets *models.SnippetModel
}

func NewApplication(db *sql.DB) *Application {
	return &Application{
		Logger:   internal.NewLogger(),
		Snippets: &models.SnippetModel{DB: db},
	}
}
