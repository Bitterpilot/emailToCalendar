package sqlite

import "database/sql"

// EmailStore holds a database.
type EmailStore struct {
	db *sql.DB
}

// CalendarStore holds a database.
type CalendarStore struct {
	db *sql.DB
}

// NewSqliteDB creates a database.
func NewSqliteDB(db *sql.DB) (*EmailStore, *CalendarStore) {
	return &EmailStore{db: db}, &CalendarStore{db: db}
}
