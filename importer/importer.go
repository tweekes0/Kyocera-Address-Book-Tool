package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/kitar0s/kyocera-ab-tool/db"
)

var (
	ErrInvalidHeader       = errors.New("invalid header")
	ErrInvalidHeaderLength = errors.New("invalid header length")
)

/*
	Checks the first line of a CSV file to ensure it has the following header:
	name,username,email.
*/

func checkCSVHeader(header []string) error {
	if len(header) != 3 {
		return ErrInvalidHeaderLength
	}

	for i := range header {
		switch i {
		case 0:
			if strings.ToLower(header[i]) != "name" {
				return ErrInvalidHeader
			}
		case 1:
			if strings.ToLower(header[i]) != "username" {
				return ErrInvalidHeader
			}
		case 2:
			if strings.ToLower(header[i]) != "email" {
				return ErrInvalidHeader
			}
		}
	}

	return nil
}

/*
	Convert a string slice from a CSV file into an  Entry.
*/

func csvToEntry(row []string) (*db.Entry, error) {
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
