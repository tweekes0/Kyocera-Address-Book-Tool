package db

import (
	"bytes"
	"testing"
)

// Test Entries
var (
	e1 = Entry{
		ID:       1,
		Name:     "Test 1",
		Username: "username1",
		Email:    "test1@test.com"}

	e2 = Entry{
		ID:       2,
		Name:     "Test 2",
		Username: "username2",
		Email:    "test2@test.com"}

	e3 = Entry{
		ID:       3,
		Name:     "Test 3",
		Username: "username3",
		Email:    "test3@test.com"}
)

func TestDisplay(t *testing.T) {
	t.Parallel()

	tt := []struct {
		description string
		got         Entry
		expect      string
	}{
		{
			description: "Entry 1",
			got:         e1,
			expect:      "ID: 1\nName: Test 1\nUsername: username1\nEmail: test1@test.com\n",
		},
		{
			description: "Entry 2",
			got:         e2,
			expect:      "ID: 2\nName: Test 2\nUsername: username2\nEmail: test2@test.com\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			var buf bytes.Buffer
			tc.got.Display(&buf)

			if buf.String() != tc.expect {
				t.Errorf("got:\n%s\nexpected:\n%s", buf.String(), tc.expect)
			}
		})
	}
}

