package db

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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

		want := 1
		assertError(t, err, nil)

		if int(got.ID) != want {
			t.Errorf("got: %d, expected: %d\n", got.ID, want)
		}
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
			got      int
			expected int
		}{
			{got: int(_e1.ID), expected: 1},
			{got: int(_e2.ID), expected: 2},
			{got: int(_e3.ID), expected: 3},
		}

		for _, tc := range tt {
			if tc.got != tc.expected {
				t.Fatalf("got: %d, expected: %d", tc.got, tc.expected)
			}
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

func assertError(t testing.TB, got, expected error) {
	if got != expected {
		t.Fatalf("got: %q, expected: %q", got, expected)
	}
}
