package infrastructure

import (
	"database/sql"
)

// Store ...
type Store struct {
	db *sql.DB
}

// NewDB ...
func NewDB(db *sql.DB) *Store {
	return &Store{db: db}
}

// ListByMsgID ...
func (s Store) ListByMsgID(msgID string) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Commit()

	var msgList string

	stmt, err := tx.Query("SELECT msgID FROM emails WHERE msgID = ?", msgID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	for stmt.Next() {
		// for each row, scan the result into our composite object
		err = stmt.Scan(&msgList)
		if err != nil {
			return "", err
		}
	}
	return msgList, nil
}
