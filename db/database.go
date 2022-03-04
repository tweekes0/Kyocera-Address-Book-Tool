package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
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
	query := fmt.Sprintf(createTable, r.currentTable)

	_, err := r.db.Exec(query)

	if err != nil {
		log.Fatalf("cannot create table: %q", err)
	}

	return err
}

func (r *SQLiteRepository) Insert(e Entry) (*Entry, error) {
	err := ValidateEntry(&e)

	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(insert, r.currentTable)
	res, err := r.db.Exec(query, e.Name, e.Username, e.Email)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
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
	query := fmt.Sprintf(selectAll, r.currentTable)
	rows, err := r.db.Query(query)

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
	query := fmt.Sprintf(selectByUername, r.currentTable)
	row := r.db.QueryRow(query, username)

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

	query := fmt.Sprintf(update, r.currentTable)
	res, err := r.db.Exec(query, updated.Name, updated.Username, updated.Email, id)

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
	if id <= 0 {
		return ErrInvalidID
	}

	query := fmt.Sprintf(delete, r.currentTable)
	res, err := r.db.Exec(query, id)
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

func (r *SQLiteRepository) NewTable(tableName string) error {
	exists, err := r.TableExists(tableName)

	if err != nil && !errors.Is(err, ErrTableDoesNotExist) {
		return err
	}

	if exists {
		return ErrTableExists
	}

	query := fmt.Sprintf(createTable, tableName)

	_, err = r.db.Exec(query, tableName)

	if err != nil {
		log.Fatalf("cannot create table: %q", err)
	}

	r.currentTable = tableName
	return nil
}

func (r *SQLiteRepository) SwitchTable(tableName string) error {
	exists, err := r.TableExists(tableName)

	if err != nil {
		return err
	}

	if !exists {
		return ErrTableDoesNotExist
	}

	r.currentTable = tableName
	return nil
}

func (r *SQLiteRepository) ClearTable() error {

	query := fmt.Sprintf(clearTable, r.currentTable)
	_, err := r.db.Exec(query)

	if err != nil {
		log.Fatalf("could not drop table: %q", err)
	}

	return nil
}

func (r *SQLiteRepository) TableExists(tableName string) (bool, error) {
	err := validateTableName(tableName)

	if err != nil {
		return false, ErrInvalidTableName
	}

	rows, err := r.db.Query(selectTables, tableName)

	if err != nil {
		log.Fatalf("cannot query database: %q", err)
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, ErrTableDoesNotExist
}
