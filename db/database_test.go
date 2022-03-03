package db

import (
	"testing"
)

func TestInsert(t *testing.T) {

	t.Run("single insert", func(t *testing.T) {
		repo, teardown := setup(t)
		defer teardown()

		got, err := repo.Insert(*e1)
		assertError(t, err, nil)

		expected, _ := NewEntry(1, "Test One", "username1", "test1@test.com")
		assertError(t, err, nil)

		assertEntry(t, got, expected)
	})

	t.Run("multiple inserts", func(t *testing.T) {
		repo, teardown := setup(t)
		defer teardown()

		_e1, err1 := repo.Insert(*e1)
		_e2, err2 := repo.Insert(*e2)
		_e3, err3 := repo.Insert(*e3)

		tt := []struct {
			got      entryInfo
			expected entryInfo
		}{
			{
				got: entryInfo{
					entry: _e1,
					err:   err1,
				},
				expected: entryInfo{
					entry: e1,
					err:   nil,
				},
			},
			{
				got: entryInfo{
					entry: _e2,
					err:   err2,
				},
				expected: entryInfo{
					entry: e2,
					err:   nil,
				},
			},
			{
				got: entryInfo{
					entry: _e3,
					err:   err3,
				},
				expected: entryInfo{
					entry: e3,
					err:   nil,
				},
			},
		}

		for _, tc := range tt {
			assertEntryInfo(t, tc.got, tc.expected)
		}
	})

	t.Run("insert non-unique entry", func(t *testing.T) {
		repo, teardown := setup(t)
		defer teardown()

		_, err := repo.Insert(*e1)
		assertError(t, err, nil)

		_, err = repo.Insert(*e1)
		assertError(t, err, ErrDuplicate)
	})
}

func TestAll(t *testing.T) {
	repo, teardown := setupWithInserts(t)
	defer teardown()

	entries := []Entry{*e1, *e2, *e3}

	all, err := repo.All()
	assertError(t, err, nil)

	for i, entry := range all {
		assertEntry(t, &entries[i], &entry)
	}
}

func TestGetByUsername(t *testing.T) {
	repo, teardown := setupWithInserts(t)
	defer teardown()

	found, foundErr := repo.GetByUsername("username1")
	notFound, notFoundErr := repo.GetByUsername("unknown user")

	tt := []entryTest{
		{
			description: "search known user",
			got: entryInfo{
				entry: found,
				err:   foundErr,
			},
			expected: entryInfo{
				entry: e1,
				err:   nil,
			},
		},
		{
			description: "search unknown user",
			got: entryInfo{
				entry: notFound,
				err:   notFoundErr,
			},
			expected: entryInfo{
				entry: nil,
				err:   ErrNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			assertEntryInfo(t, tc.got, tc.expected)
		})
	}
}

func TestUpdate(t *testing.T) {
	repo, teardown := setupWithInserts(t)
	defer teardown()

	updated, err := NewEntry(1, "new name", "newUsername", "newemail@test.com")
	assertError(t, err, nil)

	found, foundErr := repo.Update(1, *updated)
	notFound, notFoundErr := repo.Update(999999, *updated)

	tt := []entryTest{
		{
			description: "update existing user",
			got: entryInfo{
				entry: found,
				err:   foundErr,
			},
			expected: entryInfo{
				entry: updated,
				err:   nil,
			},
		},
		{
			description: "update non-existent user",
			got: entryInfo{
				entry: notFound,
				err:   notFoundErr,
			},
			expected: entryInfo{
				entry: nil,
				err:   ErrUpdateFailed,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			assertEntryInfo(t, tc.got, tc.expected)
		})
	}
}

func TestDelete(t *testing.T) {
	repo, teardown := setupWithInserts(t)
	defer teardown()

	foundErr := repo.Delete(1)
	notFoundErr := repo.Delete(99999)

	e, err := NewEntry(4, "Test One", "username1", "test1@test.com")
	assertError(t, err, nil)

	inserted, err := repo.Insert(*e)
	assertEntry(t, inserted, e)
	assertError(t, err, nil)

	newFoundErr := repo.Delete(4)
	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "delete known user",
			got:         foundErr,
			expected:    nil,
		},
		{
			description: "delete unknown user",
			got:         notFoundErr,
			expected:    ErrDeleteFailed,
		},
		{
			description: "delete user after insert",
			got:         newFoundErr,
			expected:    nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			assertError(t, tc.got, tc.expected)
		})
	}
}

func TestNewTable(t *testing.T) {
	repo, teardown := setup(t)
	defer teardown()

	err1 := repo.NewTable("valid_table_name")
	err2 := repo.NewTable("____invalid_table_name")
	err3 := repo.NewTable("default_table")

	tt := []struct {
		description string
		got         error
		expected    error
	}{
		{
			description: "create table with valid name",
			got:         err1,
			expected:    nil,
		},
		{
			description: "create table with invalid name",
			got:         err2,
			expected:    ErrInvalidTableName,
		},
		{
			description: "create table with a existing name",
			got:         err3,
			expected:    ErrTableExists,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			assertError(t, tc.got, tc.expected)
		})
	}
}
