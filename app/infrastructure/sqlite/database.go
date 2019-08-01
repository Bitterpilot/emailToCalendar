package sqlite

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/bitterpilot/emailToCalendar/models"
)

// SelectByMsgID finds all messages in the database.
func SelectByMsgID(msgID string, db *sql.DB) (string, error) {
	tx, err := db.Begin()
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

// Insert inserts an email into the database.
func Insert(e models.Email, db *sql.DB) (int, error) {
	tx, err := db.Begin()
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
