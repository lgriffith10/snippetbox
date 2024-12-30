package models

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func newTestDb(t *testing.T) *sql.DB {
	credentials := struct {
		username string
		password string
		port     string
		database string
	}{
		username: "test",
		password: "test",
		port:     "3307",
		database: "snippetboxtest",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(:%s)/%s?parseTime=true", credentials.username, credentials.password, credentials.port, credentials.database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	err = executeSQLFile(db, "./fixtures/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		err := executeSQLFile(db, "./fixtures/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}

func executeSQLFile(db *sql.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("could not read SQL file: %v", err)
	}

	sqlContent := string(sqlBytes)

	sqlContent = strings.ReplaceAll(sqlContent, "\r\n", "\n")
	sqlContent = strings.ReplaceAll(sqlContent, "\r", "\n")

	statements := strings.Split(sqlContent, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			_, err := db.Exec(stmt + ";")
			if err != nil {
				return fmt.Errorf("could not execute statement: %v", err)
			}
		}
	}

	return nil
}
