package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tweekes0/kyocera-ab-tool/db"
)

var (
	ErrInvalidHeader       = errors.New("invalid header")
	ErrInvalidHeaderLength = errors.New("invalid header length")
	ErrInvalidRowLength    = errors.New("invalid row length")
	ErrNoRowsInFile        = errors.New("there are no rows in this file")
)

/*
	Checks the first line of a CSV file to ensure it has the following header:
	name,username,email.
*/

func checkCSVHeader(header []string) error {
	if len(header) != 3 {
		return ErrInvalidHeaderLength
	}

	valid := []string{"name", "username", "email"}
	for i := range header {
		if strings.ToLower(header[i]) != valid[i] {
			return ErrInvalidHeader
		}
	}

	return nil
}

/*
	Convert a string slice from a CSV file into an  Entry.
*/

func csvToEntry(row []string) (*db.Entry, error) {
	if len(row) != 3 {
		return nil, ErrInvalidRowLength
	}

	e, err := db.NewEntry(row[0], row[1], row[2])

	if err != nil {
		return nil, err
	}

	return e, nil
}

/*
	Reads csv lines from an io.Reader and returns a slice of entries if they are
	all valid, if not returns an error.
*/

func ImportCSV(rd io.Reader) ([]*db.Entry, error) {
	csvReader := csv.NewReader(rd)
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	err = checkCSVHeader(header)
	if err != nil {
		return nil, err
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, ErrNoRowsInFile
	}

	var entries []*db.Entry
	for i, row := range rows {
		e, err := csvToEntry(row)
		if err != nil {
			s := fmt.Sprintf("%v on line %d", err, i+2)
			return nil, errors.New(s)
		}

		entries = append(entries, e)
	}

	return entries, nil
}
