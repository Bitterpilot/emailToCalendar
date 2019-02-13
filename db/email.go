package db

import "fmt"

func InsertEmail(msgID, thdID string, timeReceived int64) {
	tx, err := db.Begin()
	errorHandler(err, tx)
	defer tx.Commit()

	var eMsg string
	row := tx.QueryRow(`SELECT msgID FROM emails WHERE msgID=?`, msgID)
	if err := row.Scan(&eMsg); err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("%s\n%s\n---\n", msgID, eMsg)
	if msgID != eMsg {
		stmt, err := tx.Prepare("insert into emails(msgID,thdID, timeReceived) values(?, ?, ?)")
		errorHandler(err, tx)
		defer stmt.Close()
		_, err = stmt.Exec(msgID, thdID, timeReceived)
		errorHandler(err, tx)
	}
}

type EmailMeta struct {
	ID    int
	MsgID string
	ThdID string
}

func ListUnprocessed() []EmailMeta {
	tx, err := db.Begin()
	errorHandler(err, tx)
	defer tx.Commit()

	var unprocessed []EmailMeta
	stmt, err := tx.Query("SELECT ID, msgID, thdID FROM emails WHERE proccessed=0")
	errorHandler(err, tx)
	defer stmt.Close()
	for stmt.Next() {
		var row EmailMeta
		// for each row, scan the result into our composite object
		err = stmt.Scan(&row.ID, &row.MsgID, &row.ThdID)
		errorHandler(err, tx)
		unprocessed = append(unprocessed, row)
	}
	return unprocessed
}

func ListByThdID(thdID string) []EmailMeta {
	tx, err := db.Begin()
	errorHandler(err, tx)
	defer tx.Commit()

	var thdList []EmailMeta

	stmt, err := tx.Query("SELECT ID, msgID, thdID FROM emails WHERE thdID = ?", thdID)
	errorHandler(err, tx)
	defer stmt.Close()
	for stmt.Next() {
		var row EmailMeta
		// for each row, scan the result into our composite object
		err = stmt.Scan(&row.ID, &row.MsgID, &row.ThdID)
		errorHandler(err, tx)
		thdList = append(thdList, row)
	}
	return thdList
}

func ListByMsgID(msgID string) string {
	tx, err := db.Begin()
	errorHandler(err, tx)
	defer tx.Commit()

	var msgList string

	stmt, err := tx.Query("SELECT msgID FROM emails WHERE msgID = ?", msgID)
	errorHandler(err, tx)
	defer stmt.Close()

	for stmt.Next() {
		// for each row, scan the result into our composite object
		err = stmt.Scan(&msgList)
		errorHandler(err, tx)
	}
	return msgList
}

func MarkEmailCompleat(ID int) {
	tx, err := db.Begin()
	errorHandler(err, tx)
	defer tx.Commit()

	stmt, err := tx.Prepare("UPDATE emails SET proccessed=1 WHERE id=?")
	errorHandler(err, tx)

	res, err := stmt.Exec(ID)
	errorHandler(err, tx)

	_, err = res.RowsAffected()
	errorHandler(err, tx)

	// fmt.Println(affect)
}

// can't use ID until it is entered to the db
func MarkEmailFailed(msgID, thdID string, timeRecieved int64) {
	InsertEmail(msgID, thdID, timeRecieved)

	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	stmt, err := tx.Prepare("UPDATE emails SET error=1 WHERE msgID=?")
	errorHandler(err)

	res, err := stmt.Exec(msgID)
	errorHandler(err)

	_, err = res.RowsAffected()
	errorHandler(err)

	// fmt.Println(affect)
}
