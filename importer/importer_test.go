package importer

import (
	"errors"
	"testing"

	"github.com/tweekes0/kyocera-ab-tool/db"
)

func TestCheckCSVHeader(t *testing.T) {
	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "header is valid",
			got:         checkCSVHeader([]string{"name", "username", "email"}),
			expected:    nil,
		},
		{
			description: "header has email first",
			got:         checkCSVHeader([]string{"email", "name", "username"}),
			expected:    ErrInvalidHeader,
		},
		{
			description: "header has an unrecognized field",
			got: checkCSVHeader([]string{"doesn't belong", "username",
				"email"}),
			expected: ErrInvalidHeader,
		},
		{
			description: "header doesn't have corect number of fields",
			got: checkCSVHeader([]string{"email", "name", "username",
				"ssn"}),
			expected: ErrInvalidHeaderLength,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if !errors.Is(tc.got, tc.expected) {
				t.Fatalf("got: %v, expected: %v", tc.got, tc.expected)
			}
		})
	}
}

func TestCSVToEntry(t *testing.T) {
	_, err1 := csvToEntry([]string{"valid name", "valid_username",
		"email@email.com"})
	_, err2 := csvToEntry([]string{"name", "", ""})
	_, err3 := csvToEntry([]string{"name", ""})

	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "valid csv row",
			got:         err1,
			expected:    nil,
		},
		{
			description: "invalid csv row",
			got:         err2,
			expected:    db.ErrInvalidUsername,
		},
		{
			description: "invalid csv row",
			got:         err3,
			expected:    ErrInvalidRowLength,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if !errors.Is(tc.got, tc.expected) {
				t.Fatalf("got: %v, expected: %v", tc.got, tc.expected)
			}
		})
	}
}

func TestImportCSV(t *testing.T) {
	csv1 := [][]string{
		{"name", "username", "email"},
		{"Jane Doe", "janedoe", "janedoe@email.com"},
		{"John Doe", "johndoe", "johndoe@email.com"},
	}

	csv2 := [][]string{
		{"name", "username", "email", "phone number"},
		{"Jane Doe", "janedoe", "janedoe@email.com"},
		{"John Doe", "johndoe", "johndoe@email.com"},
	}

	csv3 := [][]string{
		{"invalid", "username", "email"},
		{"valid name", "janedoe", "janedoe@email.com"},
	}

	csv4 := [][]string{
		{"name", "username", "email"},
	}

	f1, td1 := SetupCSV(t, csv1)
	_, err1 := ImportCSV(f1)
	defer td1()

	f2, td2 := SetupCSV(t, csv2)
	_, err2 := ImportCSV(f2)
	defer td2()

	f3, td3 := SetupCSV(t, csv3)
	_, err3 := ImportCSV(f3)
	defer td3()

	f4, td4 := SetupCSV(t, csv4)
	_, err4 := ImportCSV(f4)
	defer td4()

	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "import valid csv",
			got:         err1,
			expected:    nil,
		},
		{
			description: "import csv with too many fields",
			got:         err2,
			expected:    ErrInvalidHeaderLength,
		},
		{
			description: "import csv with invalid header",
			got:         err3,
			expected:    ErrInvalidHeader,
		},
		{
			description: "import csv without rows",
			got:         err4,
			expected:    ErrNoRowsInFile,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {

			if !errors.Is(tc.got, tc.expected) {
				t.Fatalf("got: %v, expected: %v", tc.got, tc.expected)
			}
		})
	}
}
