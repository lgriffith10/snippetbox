package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox/env"
	"snippetbox/internal"
	"snippetbox/internal/models"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	Logger         *slog.Logger
	Snippets       models.SnippetModelInterface
	Users          models.UserModelInterface
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}

func newApplication(logger *slog.Logger, db *sql.DB, templateCache map[string]*template.Template) *Application {
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	return &Application{
		Logger:         logger,
		Snippets:       &models.SnippetModel{DB: db},
		Users:          &models.UserModel{DB: db},
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}
}

func main() {
	mode := flag.String("mode", "dev", "run mode")
	flag.Parse()

	env.SetEnvVariables(*mode)

	port := os.Getenv("GO_PORT")
	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		panic(err)
	}

	logger := internal.NewLogger()

	app := newApplication(logger, db, templateCache)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	srv := &http.Server{
		Addr:         port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Info("Starting server", "addr", port)

	err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	app.Logger.Error(err.Error())
	os.Exit(1)
}
