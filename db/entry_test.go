package db

import (
	"bytes"
	"testing"
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
			got:         *e1,
			expect:      "ID: 1\nName: Test One\nUsername: username1\nEmail: test1@test.com\n",
		},
		{
			description: "Entry 2",
			got:         *e2,
			expect:      "ID: 2\nName: Test Two\nUsername: username2\nEmail: test2@test.com\n",
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

func TestValidateEntry(t *testing.T) {
	_, err1 := NewEntry(1, "", "username", "example@example.com")
	_, err2 := NewEntry(2, "name", "123", "example@example.com")
	_, err3 := NewEntry(3, "name", "username", "invalid_email.com")

	tt := []struct {
		description string
		got         error
		want        error
	}{
		{"test invalid name", err1, ErrInvalidName},
		{"test invalid username", err2, ErrInvalidUsername},
		{"test invalid email", err3, ErrInvalidEmail},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			assertError(t, tc.got, tc.want)
		})
	}
}
