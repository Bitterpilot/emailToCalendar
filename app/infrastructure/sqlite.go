package infrastructure

import (
	"database/sql"

	"github.com/bitterpilot/emailToCalendar/app/infrastructure/sqlite"
	"github.com/bitterpilot/emailToCalendar/models"
)

// Store
type Store struct {
	db *sql.DB
}

// NewDB
func NewDB(db *sql.DB) *Store {
	return &Store{db: db}
}

// ListByMsgID
func (s Store) ListByMsgID(msgID string) (string, error) {
	return sqlite.SelectByMsgID(msgID, s.db)
}

// InsertEmail
func (s Store) InsertEmail(e models.Email) (int, error) {
	return sqlite.Insert(e, s.db)
}
