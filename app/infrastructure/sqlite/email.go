package sqlite

import (
	"database/sql"

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
// CheckUnProcessed
func (s EmailStore) CheckUnProcessed(e models.Email) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Commit()

	row := tx.QueryRow("SELECT msgID FROM emails WHERE proccessed=0 AND msgID=?", e.MsgID)

	var id string
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		// fmt.Println("No rows were returned!")
		return "", err
	case nil:
		return id, nil
	default:
		return "", err
	}
}

// InsertEmail into database and returns it's row ID.
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
