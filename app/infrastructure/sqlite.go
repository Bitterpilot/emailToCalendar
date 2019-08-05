package infrastructure

import (
	"database/sql"

	"github.com/bitterpilot/emailToCalendar/app/infrastructure/sqlite"
	"github.com/bitterpilot/emailToCalendar/models"
)

// Store holds a database.
type Store struct {
	db *sql.DB
}

// NewSqliteDB creates a database.
func NewSqliteDB(db *sql.DB) *Store {
	return &Store{db: db}
}

// FindByMsgID searches for records by Message ID.
func (s Store) FindByMsgID(msgID string) (string, error) {
	return sqlite.SelectByMsgID(msgID, s.db)
}

// InsertEmail into database and returns it's row ID.
func (s Store) InsertEmail(e models.Email) (int, error) {
	return sqlite.Insert(e, s.db)
}
