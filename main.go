package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tweekes0/kyocera-ab-tool/db"
	"github.com/tweekes0/kyocera-ab-tool/prompt"
)

const (
	DB_DRIVER    = "sqlite3"    // Database driver
	DB_FILENAME  = "sqlite.db"  // Database filename
	DATABASE_DIR = "./Database" // Database directory
)

func main() {
	// Create the Database directory if it doesn't exist
	_, err := os.Stat(DATABASE_DIR)
	if os.IsNotExist(err) {
		if err = os.Mkdir(DATABASE_DIR, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		msg := fmt.Sprintf("Creating the %v directory", DATABASE_DIR)
		prompt.OutputMessage(os.Stdout, '!', msg)
	}

	// Create path to database in the Database directory
	db_path := filepath.Join(DATABASE_DIR, DB_FILENAME)

	_, err = os.Stat(db_path)
	if os.IsNotExist(err) {
		_, err := os.Create(db_path)
		errChecker(err)

		msg := "Creating database file"
		prompt.OutputMessage(os.Stdout, '!', msg)
	}

	// Create a reference to a SQL database
	sqlite, err := sql.Open(DB_DRIVER, db_path)
	if err != nil {
		log.Fatalf("Could not open database: %q", err)
	}

	// Create sqlite repository and initialize it
	r, err := db.NewSQLiteRepository(sqlite)
	errChecker(err)

	err = r.Initialize()
	errChecker(err)

	// CLI application
	prompt.Prompt(r, os.Stdin, os.Stdout)
}

func errChecker(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
