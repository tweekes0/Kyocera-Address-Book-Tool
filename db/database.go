package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate       = errors.New("record already exists")
	ErrNotFound        = errors.New("record does not exist")
	ErrUpdateFailed    = errors.New("record could not be updated")
	ErrDeleteFailed    = errors.New("record could not be deleted")
	ErrInvalidID       = errors.New("record ID is invalid")
	ErrInvalidEmail    = errors.New("email is not valid")
	ErrInvalidName     = errors.New("name is not valid")
	ErrInvalidUsername = errors.New("username is not valid")
)

type Repository interface {
	Initialize() error
	Insert(e Entry) error
	All() ([]Entry, error)
	GetByUsername(username string) (*Entry, error)
	Update(id int64, update Entry) (*Entry, error)
	Delete(id int64) error
}

type SQLiteRepository struct {
	db           *sql.DB
	currentTable string
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db:           db,
		currentTable: "default_table",
	}
}

func (r *SQLiteRepository) Initialize() error {
	query := `
		CREATE TABLE IF NOT EXISTS default_table(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name text NOT NULL,
		username text UNIQUE NOT NULL, 
		email text UNIQUE NOT NULL 
		);`

	_, err := r.db.Exec(query)

	return err
}

func (r *SQLiteRepository) Insert(e Entry) (*Entry, error) {
	err := ValidateEntry(&e) 

	if err != nil {
		return nil, err
	}

	res, err := r.db.Exec(`INSERT INTO default_table(name, username, email) 
		values(?,?,?)`, e.Name, e.Username, e.Email)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		log.Fatal("Sad hood movie")
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	e.ID = id
	return &e, nil
}

func (r *SQLiteRepository) All() (all []Entry, e error) {
	rows, err := r.db.Query("SELECT * FROM default_table")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e Entry
		err := rows.Scan(&e.ID, &e.Name, &e.Username, &e.Email)
		if err != nil {
			return nil, err
		}

		all = append(all, e)
	}

	return all, nil
}

func (r *SQLiteRepository) GetByUsername(username string) (*Entry, error) {
	row := r.db.QueryRow("SELECT * FROM default_table WHERE username = ?",
		username)

	var e Entry
	err := row.Scan(&e.ID, &e.Name, &e.Username, &e.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &e, nil
}

func (r *SQLiteRepository) Update(id int64, updated Entry) (*Entry, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	res, err := r.db.Exec(`
		UPDATE default_table SET name = ?, username = ?,
		email = ? WHERE id = ?`, updated.Name, updated.Username,
		updated.Email, id)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	if id >= 0 {
		return ErrInvalidID
	}

	res, err := r.db.Exec("DELETE FROM default_table WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
