package db

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

/*
	Structs used for testing Entries and Tables respectively.
*/

type entryInfo struct {
	entry *Entry
	err   error
}

type tableInfo struct {
	tableName string
	err       error
}

/*
	Name of the default table that is always created and cannot be deleted in
	the application.
*/

const DEFAULT_TABLE = "default_table"

/*
	SQLite queries
*/

const (
	insert          = "INSERT INTO %v(name, username, email) values(?,?,?);"
	update          = "UPDATE %v SET name=?, username=?, email=? WHERE username=?;"
	delete          = "DELETE FROM %v WHERE username=?;"
	selectAll       = "SELECT * FROM %v;"
	selectTable     = "SELECT name FROM sqlite_master WHERE type='table' AND name=?;"
	selectByUername = "SELECT * FROM %v WHERE username=?;"
	createTable     = `CREATE TABLE IF NOT EXISTS %v (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name text NOT NULL,
		username text UNIQUE NOT NULL, 
		email text UNIQUE NOT NULL 
		);`
	clearTable  = "DELETE FROM %v"
	deleteTable = "DROP TABLE %v"
	listTables  = `SELECT name from sqlite_master WHERE TYPE="table" AND name 
	NOT LIKE '%sql%';`
)

/*
	Patterns for regular expressions
*/

const (
	namePattern         = `^[a-zA-Z]+([ ]?[a-zA-Z]+)*$`
	usernamePattern     = `^[a-zA-Z]+([\._-]?[a-zA-Z0-9])*$`
	emailPattern        = `^[a-zA-Z]+([\._-]?[a-zA-Z0-9])+@[a-zA-Z]+(\.[a-zA-Z]+)+$`
	tablePattern        = `^[a-zA-Z_]{1}([a-zA-Z0-9]+[_]?)*$`
	bracketTablePattern = `^[\[][a-zA-Z0-9]+([ +!?._\-a-zA-Z0-9])*[\]]$`
)

var (
	e1 = &Entry{
		ID:       1,
		Name:     "Test One",
		Username: "username1",
		Email:    "test1@test.com"}

	e2 = &Entry{
		ID:       2,
		Name:     "Test Two",
		Username: "username2",
		Email:    "test2@test.com"}

	e3 = &Entry{
		ID:       3,
		Name:     "Test Three",
		Username: "username3",
		Email:    "test3@test.com"}
)

/*
	Errors
*/

var (
	ErrDuplicate            = errors.New("record already exists")
	ErrNotFound             = errors.New("record does not exist")
	ErrUpdateFailed         = errors.New("record could not be updated")
	ErrDeleteFailed         = errors.New("record could not be deleted")
	ErrInvalidID            = errors.New("record ID is invalid")
	ErrInvalidEmail         = errors.New("email is not valid")
	ErrInvalidName          = errors.New("name is not valid")
	ErrInvalidUsername      = errors.New("username is not valid")
	ErrInvalidTableName     = errors.New("tablename is not valid")
	ErrTableExists          = errors.New("table already exists")
	ErrTableDoesNotExist    = errors.New("table does not exist")
	ErrTableCannotBeDeleted = errors.New("this table cannot be deleted")
)

func assertError(t testing.TB, got, expected error) {
	if got != expected {
		t.Fatalf("got: %v, expected: %v", got, expected)
	}
}

func assertEntry(t testing.TB, got, expected *Entry) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %v, expected: %v\n", got, expected)
	}
}

func assertEntryInfo(t testing.TB, got, expected entryInfo) {
	assertError(t, got.err, expected.err)
	assertEntry(t, got.entry, expected.entry)
}

func assertTableInfo(t testing.TB, got, expected tableInfo) {
	assertError(t, got.err, expected.err)
	if got.tableName != expected.tableName {
		t.Fatalf("got: %v, expected: %v", got.tableName, expected.tableName)
	}
}

/*
	Function to create a database reference and return the clean up function
*/

func setup(t *testing.T) (*SQLiteRepository, func()) {
	t.Parallel()

	f, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatalf("could not create file: %q", err)
	}

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		log.Fatalf("could not open sqlite db: %q", err)
	}

	entryRepo, err := NewSQLiteRepository(db)
	if err != nil {
		log.Fatalf("could not create db connection: %q", err)
	}

	err = entryRepo.Initialize()
	if err != nil {
		log.Fatalf("could not initialize sqlite db: %q", err)
	}

	teardown := func() {
		os.Remove(f.Name())
	}

	return entryRepo, teardown
}

/*
	Function to create a database and load Entries into it and return
	clean up function.
*/

func SetupWithInserts(t *testing.T) (*SQLiteRepository, func()) {
	entryRepo, teardown := setup(t)

	_, err := entryRepo.Insert(*e1)
	assertError(t, err, nil)
	_, err = entryRepo.Insert(*e2)
	assertError(t, err, nil)
	_, err = entryRepo.Insert(*e3)
	assertError(t, err, nil)

	return entryRepo, teardown
}

/*
	Function to make sure tables are named properly. and cannot contain sql
*/

func validateTableName(tableName string) error {
	pat1 := regexp.MustCompile(tablePattern)
	pat2 := regexp.MustCompile(bracketTablePattern)
	b := !strings.Contains(tableName, "sql") && tableName != "table"

	if (pat1.MatchString(tableName) || pat2.MatchString(tableName)) && b {
		return nil
	}

	return ErrInvalidTableName
}
