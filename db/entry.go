package db

import (
	"fmt"
	"io"
	"regexp"
)

type BasicEntry interface {
	Display()
}

type Entry struct {
	ID       int64
	Name     string
	Username string
	Email    string
}

func NewEntry(id int64, name, username, email string) (*Entry, error) {
	p := new(Entry)
	p.Name = name
	p.Username = username
	p.Email = email

	err := ValidateEntry(p)
	if err != nil {
		return nil, err
	}

	return p, err
}

func validateField(field, pattern string, err error) error {
	r := regexp.MustCompile(pattern)

	if !r.MatchString(field) {
		return err
	}

	return nil
}

func ValidateEntry(e *Entry) error {
	const namePattern = "^[a-zA-Z]+([ ]?[a-zA-Z]{1,})*"
	const usernamePattern = "^[a-zA-Z]+([\\._-]?[a-zA-Z0-9])*"
	const emailPattern = "^[a-zA-Z]+([\\._-]?[a-zA-Z0-9])+@[a-zA-Z]+(\\.[a-zA-Z]+)+"

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

func (e Entry) Display(writer io.Writer) {
	fmt.Fprintf(writer, "ID: %d\nName: %v\nUsername: %v\nEmail: %v\n",
		e.ID, e.Name, e.Username, e.Email)
}
