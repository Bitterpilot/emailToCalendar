package sqlite

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/bitterpilot/emailToCalendar/models"
)

// FindByMsgID searches for records by Message ID.
func (s EmailStore) FindByMsgID(msgID string) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Commit()

	var msg string

	stmt, err := tx.Query("SELECT msgID FROM emails WHERE msgID = ?", msgID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	for stmt.Next() {
		// for each row, scan the result into our composite object
		err = stmt.Scan(&msg)
		if err != nil {
			return "", err
		}
	}
	return msg, nil
}

// ListUnprocessed
func (s EmailStore) ListUnprocessed(e models.Email) (models.Email, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return models.Email{}, err
	}
	defer tx.Commit()

	row := tx.QueryRow("SELECT id, proccessed FROM emails WHERE msgID=?", e.MsgID)

	switch err := row.Scan(&e.ID, &e.Processed); err {
	case sql.ErrNoRows:
		// fmt.Println("No rows were returned!")
		return e, nil
	case nil:
		return e, nil
	default:
		return models.Email{}, err
	}
}

// InsertEmail into database and returns it's row ID.
// This doesn't check for duplicates
func (s EmailStore) InsertEmail(e models.Email) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare("INSERT INTO emails(msgID, thdID, timeRecieved) values(?,?,?)")
	if err != nil {
		return -1, err
	}

	res, err := stmt.Exec(e.MsgID, e.ThdID, e.TimeReceived)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	log.WithFields(log.Fields{"Database ID": id}).Info("Inserted to DB.")
	return int(id), nil
}

// MarkAsProcessed Marks an email as complete in the database.
func (s EmailStore) MarkAsProcessed(e models.Email) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare("UPDATE emails SET proccessed=? WHERE id=? ")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(e.Processed, e.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("%d Rows affected, expected 1", rows)
	}
	return nil
}

// ListByThdID
func (s EmailStore) ListByThdID(thdID string) ([]models.Email, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	var thdList []models.Email

	stmt, err := tx.Query("SELECT ID, msgID, thdID FROM emails WHERE thdID = ?", thdID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	for stmt.Next() {
		var row models.Email
		// for each row, scan the result into our composite object
		err = stmt.Scan(&row.ID, &row.MsgID, &row.ThdID)
		if err != nil {
			return nil, err
		}
		thdList = append(thdList, row)
	}
	return thdList, nil
}
