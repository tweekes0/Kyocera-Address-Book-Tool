package db

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)
 
func assertError(t testing.TB, got, expected error) {
	if got != expected {
		t.Fatalf("got: %q, expected: %q", got, expected)
	}
}

func assertEntry(t testing.TB, got, expected Entry) {
	if (!reflect.DeepEqual(got, expected)) {
		t.Errorf("got: %v, expected: %v\n", got, expected)
	}
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

		got, err := repo.Insert(e1)
		assertError(t, err, nil)

		expected, _ := NewEntry(1, "Test 1" , "username1","test1@test.com")
		assertError(t, err, nil)

		assertEntry(t, *got, *expected)
	})

	t.Run("multiple inserts", func(t *testing.T) {
		repo, teardown := setup(t)
		defer teardown()

		_e1, err := repo.Insert(e1)
		assertError(t, err, nil)

		_e2, err := repo.Insert(e2)
		assertError(t, err, nil)

		_e3, err := repo.Insert(e3)
		assertError(t, err, nil)

		tt := []struct {
			got      Entry
			expected Entry
		}{
			{got: *_e1, expected: e1},
			{got: *_e2, expected: e2},
			{got: *_e3, expected: e3},
		}

		for _, tc := range tt {
			assertEntry(t, tc.got, tc.expected)
		}
	})

	t.Run("insert non-unique entry", func(t *testing.T) {
		repo, teardown := setup(t)
		defer teardown()

		_, err := repo.Insert(e1)
		assertError(t, err, nil)

		_, err = repo.Insert(e1)
		assertError(t, err, ErrDuplicate)
	})
}
