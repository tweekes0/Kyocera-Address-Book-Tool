package db

import (
	"errors"
)

type databaseTest struct {
	description string
	got         entryInfo
	expected    entryInfo
}

type entryInfo struct {
	entry *Entry
	err   error
}

const (
	insert          = "INSERT INTO %v(name, username, email) values(?,?,?)"
	update          = "UPDATE %v SET name = ?, username = ?, email = ? WHERE id = ?"
	delete          = "DELETE FROM %v WHERE id = ?"
	selectAll       = "SELECT * FROM %v"
	selectTables    = "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	selectByUername = "SELECT * FROM %v WHERE username = ?"
	createTable     = `CREATE TABLE IF NOT EXISTS %v (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name text NOT NULL,
		username text UNIQUE NOT NULL, 
		email text UNIQUE NOT NULL 
		);`
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
	ErrDuplicate        = errors.New("record already exists")
	ErrNotFound         = errors.New("record does not exist")
	ErrUpdateFailed     = errors.New("record could not be updated")
	ErrDeleteFailed     = errors.New("record could not be deleted")
	ErrInvalidID        = errors.New("record ID is invalid")
	ErrInvalidEmail     = errors.New("email is not valid")
	ErrInvalidName      = errors.New("name is not valid")
	ErrInvalidUsername  = errors.New("username is not valid")
	ErrInvalidTableName = errors.New("tablename is not valid")
	ErrTableExists      = errors.New("table already exists")
)
