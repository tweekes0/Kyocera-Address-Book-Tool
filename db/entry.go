package db

import (
	"fmt"
	"io"
	"regexp"
)

/*
	Entry struct models how data will be inserted into database.

	ID: a unique id of the Entry's record in the database,
	Name: the owner of the entry
	Username: unique identifier for the entry
	Email: email address of the Entry's owner
*/

type Entry struct {
	ID       int64
	Name     string
	Username string
	Email    string
}

/*
	A receiver function for Entry to be written to an io.Writer

	writer: an output stream where the properties of the entry can be written
*/

func (e *Entry) Display(writer io.Writer) {
	fmt.Fprintf(writer, "ID: %d\nName: %v\nUsername: %v\nEmail: %v\n",
		e.ID, e.Name, e.Username, e.Email)
}

/*
	Entry struct constructor.

	Given the proper parameters,  a reference to an Entry will be returned.
	If there is an issue with one of the fields an error will be returned.
*/

func NewEntry(name, username, email string) (*Entry, error) {
	p := new(Entry)
	p.Name = name
	p.Username = username
	p.Email = email

	err := validateEntry(p)
	if err != nil {
		return nil, err
	}

	return p, err
}

func newTestEntry(id int64, name, username, email string) (*Entry, error) {
	p := new(Entry)
	p.ID = id
	p.Name = name
	p.Username = username
	p.Email = email

	err := validateEntry(p)
	if err != nil {
		return nil, err
	}

	return p, err
}

/*
	Function to ensure that the fields in the Entry struct conform to a specific
	pattern.

	field: a string value from the Entry struct
	pattern: a pattern that the field must meet
	err: The error to be return if the field fails to conform to the pattern
*/

func validateField(field, pattern string, err error) error {
	r := regexp.MustCompile(pattern)

	if !r.MatchString(field) {
		return err
	}

	return nil
}

/*
	Function that checks all the string fields of a given Entry. If a field
	fails to conform to its pattern a corresponding error is returned.

	e: Pointer for an Entry that is to be checked.
*/

func validateEntry(e *Entry) error {

	err := validateField(e.Name, namePattern, ErrInvalidName)
	if err != nil {
		return err
	}

	err = validateField(e.Username, usernamePattern, ErrInvalidUsername)
	if err != nil {
		return err
	}

	err = validateField(e.Email, emailPattern, ErrInvalidEmail)
	if err != nil {
		return err
	}

	return nil
}
