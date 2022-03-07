package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
)

/*
	SQLiteRepository struct abstracts SQLite database

	db: reference to a database, enabling db operations
	currentTable: the table that certain statements will be ran against
*/

type SQLiteRepository struct {
	db           *sql.DB
	currentTable string
}

/*
	SQLiteRepository struct constructor

	Given a proper SQL database reference. A reference to a new SQLiteRepository
	will be returned.
*/

func NewSQLiteRepository(db *sql.DB) (*SQLiteRepository, error) {
	err := validateTableName(DEFAULT_TABLE)
	if err != nil {
		return nil, err
	}

	return &SQLiteRepository{
		db:           db,
		currentTable: DEFAULT_TABLE,
	}, nil
}

func (r *SQLiteRepository) CurrentTable() string {
	return r.currentTable
}

/*
	Createas the default table. 
	
	Logs to console and terminates execution if there is an issue with SQL
*/

func (r *SQLiteRepository) Initialize() error {
	query := fmt.Sprintf(createTable, DEFAULT_TABLE)

	_, err := r.db.Exec(query)

	if err != nil {
		log.Fatalf("cannot create table: %q", err)
	}

	return err
}

/*
	Inserts en Entry into currentTable and returns the reference of the Entry 
	with an ID given to it from the database.

	Logs to console and terminates execution if there is an issue with SQLer
*/

func (r *SQLiteRepository) Insert(e Entry) (*Entry, error) {
	err := validateEntry(&e)
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
		log.Fatalf("cannot insert record into table: %q", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	e.ID = id
	return &e, nil
}

/*
	Queries currentTable return all the Entries  

	Logs to console and terminates execution if there is an issue with SQL
*/

func (r *SQLiteRepository) All() (all []*Entry, err error) {
	query := fmt.Sprintf(selectAll, r.currentTable)
	rows, err := r.db.Query(query)

	if err != nil {
		log.Fatalf("cannot query table: %q", err)
	}
	defer rows.Close()

	for rows.Next() {
		e := new(Entry)
		err := rows.Scan(&e.ID, &e.Name, &e.Username, &e.Email)
		if err != nil {
			log.Fatalf("cannot scan row: %q", err)
		}

		all = append(all, e)
	}

	return all, nil
}

/*
	Queries currentTable to return a reference to an Entry when given a valid
	username.

	Logs to console and terminates execution if there is an issue with SQL
*/

func (r *SQLiteRepository) GetByUsername(username string) (*Entry, error) {
	err := validateField(username, usernamePattern, ErrInvalidUsername)
	if err != nil {
		return nil, err
	}
	
	query := fmt.Sprintf(selectByUername, r.currentTable)
	row := r.db.QueryRow(query, username)

	e := new(Entry)
	err = row.Scan(&e.ID, &e.Name, &e.Username, &e.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
	}

	return e, nil
}

/*
	Updates an Entry in the currentTable given it's id with the newly 
	updated Entry. Returns the updated entry if there are no issues. 
	If there are no updates no Entry is returned and corresponding 
	error is returned also.

	Logs error to console and terminates execution if there is an issue with the 
	SQL.
*/

func (r *SQLiteRepository) Update(id int64, updated *Entry) (*Entry, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	err := validateEntry(updated) 
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(update, r.currentTable)

	res, err := r.db.Exec(query, updated.Name, updated.Username, updated.Email, id)
	if err != nil {
		log.Fatalf("cannot execute statement: %q", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("cannot update record: %q", err)
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return updated, nil
}

/*
	Deletes an Entry in the currentTable given a valid id.

	Logs to console and terminates execution if there is an issue with SQL
*/

func (r *SQLiteRepository) Delete(username string) error {
	err := validateField(username, usernamePattern, ErrInvalidUsername)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(delete, r.currentTable)

	res, err := r.db.Exec(query, username)
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

/*
	Creates a new table within the database. currentTable is updated if the 
	table is valid and does not already exist.

	Logs to console and terminates execution if there is an issue with SQL
*/

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

/*
	Updates the currentTable variable to to the tableName parameter given the
	table exists and has a valid name.

	Logs to console and terminates execution if there is an issue with SQL
*/

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

/*
	Wipes the current table.

	Logs to console and terminates execution if there is an issue with SQL
*/
func (r *SQLiteRepository) ClearTable() error {
	query := fmt.Sprintf(clearTable, r.currentTable)

	_, err := r.db.Exec(query)
	if err != nil {
		log.Fatalf("could not drop table: %q", err)
	}

	return nil
}

/*
	Checks for the existence of a the given tableName and will return a bool
	for the table's existence. If the table does not exist an error will be 
	returned to the caller as to why not.

	Logs to console and terminates execution if there is an issue with SQL
*/

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

/*
	Delete the table from the database that is passed. This function cannot and
	will not delete the DEFAULT_TABLE.
*/

func (r *SQLiteRepository) DeleteTable(tableName string) error {
	_, err := r.TableExists(tableName)

	if err != nil {
		return err
	}

	if tableName == DEFAULT_TABLE {
		return ErrTableCannotBeDeleted
	}

	if r.currentTable == tableName {
		r.currentTable = DEFAULT_TABLE
	}

	query := fmt.Sprintf(deleteTable, tableName)

	_, err = r.db.Exec(query)
	if err != nil {
		return ErrTableCannotBeDeleted
	}

	return nil
}
