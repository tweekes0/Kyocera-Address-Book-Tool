package db

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
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

func assertError(t testing.TB, got, expected error) {
	if got != expected {
		t.Fatalf("got: %q, expected: %q", got, expected)
	}
}

func assertEntry(t testing.TB, got, expected *Entry) {
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got: %v, expected: %v\n", got, expected)
	}
}

func assertEntryInfo(t testing.TB, got, expected entryInfo) {
	assertError(t, got.err, expected.err)
	assertEntry(t, got.entry, expected.entry)
}

func setup(t *testing.T) (*SQLiteRepository, func()) {
	t.Parallel()

	f, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		log.Fatal(err)
	}

	teardown := func() {
		os.Remove(f.Name())
	}

	entryRepo := NewSQLiteRepository(db)
	entryRepo.Initialize()

	return entryRepo, teardown
}

func TestInsert(t *testing.T) {

	t.Run("single insert", func(t *testing.T) {
		repo, teardown := setup(t)
		defer teardown()

		got, err := repo.Insert(*e1)
		assertError(t, err, nil)

		expected, _ := NewEntry(1, "Test 1", "username1", "test1@test.com")
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
	repo, teardown := setup(t)
	defer teardown()

	_e1, err := repo.Insert(*e1)
	assertError(t, err, nil)

	_e2, err := repo.Insert(*e2)
	assertError(t, err, nil)

	_e3, err := repo.Insert(*e3)
	assertError(t, err, nil)

	entries := []Entry{*_e1, *_e2, *_e3}

	all, err := repo.All()
	assertError(t, err, nil)

	for i, entry := range all {
		assertEntry(t, &entries[i], &entry)
	}
}

func TestGetByUsername(t *testing.T) {
	repo, teardown := setup(t)
	defer teardown()

	_e1, err := repo.Insert(*e1)
	assertError(t, err, nil)

	found, foundErr := repo.GetByUsername("username1")
	notFound, notFoundErr := repo.GetByUsername("unknown user")

	tt := []databaseTest{
		{
			description: "test known user",
			got: entryInfo{
				entry: found,
				err:   foundErr,
			},
			expected: entryInfo{
				entry: _e1,
				err:   nil,
			},
		},
		{
			description: "test unknown user",
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
		assertEntryInfo(t, tc.got, tc.expected)
	}
}
