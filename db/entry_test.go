package db

import (
	"bytes"
	"testing"
)

func TestDisplay(t *testing.T) {
	tt := []struct {
		description string
		got         Entry
		expected    string
	}{
		{
			description: "Entry 1",
			got:         *e1,
			expected:    "ID: 1\nName: Test One\nUsername: username1\nEmail: test1@test.com\n",
		},
		{
			description: "Entry 2",
			got:         *e2,
			expected:    "ID: 2\nName: Test Two\nUsername: username2\nEmail: test2@test.com\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			tc.got.Display(&buf)

			if buf.String() != tc.expected {
				t.Errorf("got:\n%s\nexpected:\n%s", buf.String(), tc.expected)
			}
		})
	}
}

func TestNewEntry(t *testing.T) {
	_e1, _ := NewEntry("valid name", "validusername", "valid@email.com")
	_e2, _ := NewEntry("invalid_name", "invalid user", "email.com")
	tt := []struct {
		description string
		got         *Entry
		expected    *Entry
	}{
		{
			description: "new valid entry",
			got:         _e1,
			expected: &Entry{
				ID:       0,
				Name:     "valid name",
				Username: "validusername",
				Email:    "valid@email.com",
			},
		},
		{
			description: "new invalid entry",
			got:         _e2,
			expected:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			assertEntry(t, tc.got, tc.expected)
		})
	}
}

func TestNewTestEntry(t *testing.T) {
	_e1, _ := newTestEntry(1, "valid name", "username", "example@example.com")
	_e2, _ := newTestEntry(2, "name", "123", "example@example.com")

	tt := []struct {
		description string
		got         *Entry
		expected    *Entry
	}{
		{
			description: "new valid test entry",
			got:         _e1,
			expected: &Entry{
				ID:       1,
				Name:     "valid name",
				Username: "username",
				Email:    "example@example.com",
			},
		},
		{
			description: "new invalid test entry",
			got:         _e2,
			expected:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			assertEntry(t, tc.got, tc.expected)
		})
	}
}

func TestValidateEntry(t *testing.T) {
	_, err1 := newTestEntry(1, "", "username", "example@example.com")
	_, err2 := newTestEntry(2, "name", "123", "example@example.com")
	_, err3 := newTestEntry(3, "name", "username", "invalid_email.com")

	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "test invalid name",
			got:         err1,
			expected:    ErrInvalidName,
		},
		{
			description: "test invalid username",
			got:         err2,
			expected:    ErrInvalidUsername,
		},
		{
			description: "test invalid email",
			got:         err3,
			expected:    ErrInvalidEmail,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			assertError(t, tc.got, tc.expected)
		})
	}
}

func TestValidateField(t *testing.T) {
	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "validate invalid name",
			got:         validateField("invalid_name", namePattern, ErrInvalidName),
			expected:    ErrInvalidName,
		},
		{
			description: "validate valid name",
			got:         validateField("valid name", namePattern, ErrInvalidName),
			expected:    nil,
		},
		{
			description: "validate invalid username",
			got:         validateField("invalid username", usernamePattern, ErrInvalidUsername),
			expected:    ErrInvalidUsername,
		},
		{
			description: "validate valid username",
			got:         validateField("validusername", usernamePattern, ErrInvalidUsername),
			expected:    nil,
		},
		{
			description: "validate invalid email",
			got:         validateField("bad@email@email.com", emailPattern, ErrInvalidEmail),
			expected:    ErrInvalidEmail,
		},
		{
			description: "validate valid email",
			got:         validateField("valid@email.com", emailPattern, ErrInvalidEmail),
			expected:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			assertError(t, tc.got, tc.expected)
		})
	}
}
