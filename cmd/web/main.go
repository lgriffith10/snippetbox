package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"snippetbox/env"
	"snippetbox/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

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

	app := config.NewApplication(db)

	app.Logger.Info("Starting server", "addr", port)

	err = http.ListenAndServe(port, routes(app))
	app.Logger.Error(err.Error())
	os.Exit(1)
}
