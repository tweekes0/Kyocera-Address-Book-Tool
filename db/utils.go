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

type entryInfo struct {
	entry *Entry
	err   error
}

type tableInfo struct {
	tableName string
	err       error
}

const (
	insert          = "INSERT INTO %v(name, username, email) values(?,?,?)"
	update          = "UPDATE %v SET name = ?, username = ?, email = ? WHERE id = ?"
	// delete          = "DELETE FROM %v WHERE id = ?"
	delete          = "DELETE FROM %v WHERE username = ?"
	selectAll       = "SELECT * FROM %v"
	selectTables    = "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	selectByUername = "SELECT * FROM %v WHERE username = ?"
	createTable     = `CREATE TABLE IF NOT EXISTS %v (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name text NOT NULL,
		username text UNIQUE NOT NULL, 
		email text UNIQUE NOT NULL 
		);`
	clearTable  = "DELETE FROM %v"
	deleteTable = "DROP TABLE %v"

	DEFAULT_TABLE = "default_table"

	namePattern = "^[a-zA-Z]+([ ]?[a-zA-Z]+)*$"
	usernamePattern = "^[a-zA-Z]+([\\._-]?[a-zA-Z0-9])*$"
	emailPattern = "^[a-zA-Z]+([\\._-]?[a-zA-Z0-9])+@[a-zA-Z]+(\\.[a-zA-Z]+)+$"
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
		t.Fatalf("got: %q, expected: %q", got, expected)
	}
}

func assertEntry(t testing.TB, got, expected *Entry) {
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got: %v, expected: %v\n", got, expected)
	}
}

func assertEntryInfo(t testing.TB, got, expected entryInfo) {
	assertError(t, got.err, expected.err)
	assertEntry(t, got.entry, expected.entry)
}

func assertTableInfo(t testing.TB, got, expected tableInfo) {
	assertError(t, got.err, expected.err)
	if got.tableName != expected.tableName {
		t.Errorf("got: %v, expected: %v", got.tableName, expected.tableName)
	}
}

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

	teardown := func() {
		os.Remove(f.Name())
	}

	entryRepo, err := NewSQLiteRepository(db)
	if err != nil {
		log.Fatalf("could not create db connection: %q", err)
	}
	
	err = entryRepo.Initialize()
	if err != nil {
		log.Fatalf("could not initialize sqlite db: %q", err)
	}

	return entryRepo, teardown
}

func setupWithInserts(t *testing.T) (*SQLiteRepository, func()) {
	entryRepo, teardown := setup(t)

	_, err := entryRepo.Insert(*e1)
	assertError(t, err, nil)
	_, err = entryRepo.Insert(*e2)
	assertError(t, err, nil)
	_, err = entryRepo.Insert(*e3)
	assertError(t, err, nil)

	return entryRepo, teardown
}

func validateTableName(tableName string) error {
	const tablePattern = "^[a-zA-Z]{2,}([-_]{1}[a-zA-Z0-9]+)*$"
	re := regexp.MustCompile(tablePattern)

	if !re.MatchString(tableName) || strings.Contains(tableName, "sqlite") {
		return ErrInvalidTableName
	}

	return nil
}
