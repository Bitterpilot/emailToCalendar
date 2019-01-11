package db

import "fmt"

func InsertEmail(msgID, thdID string, timeRecieved int64) {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	var eMsg string
	row := tx.QueryRow(`SELECT msgID FROM emails WHERE msgID=?`, msgID)
	if err := row.Scan(&eMsg); err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("%s\n%s\n---\n", msgID, eMsg)
	if msgID != eMsg {
		stmt, err := tx.Prepare("insert into emails(msgID,thdID, timeRecieved) values(?, ?, ?)")
		errorHandler(err)
		defer stmt.Close()
		_, err = stmt.Exec(msgID, thdID, timeRecieved)
		errorHandler(err)
	}
}

type EmailMeta struct {
	ID    int
	MsgID string
	ThdID string
}

func ListUnprocssed() []EmailMeta {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	var unprocessed []EmailMeta
	stmt, err := tx.Query("SELECT ID, msgID, thdID FROM emails WHERE proccessed=0")
	errorHandler(err)
	defer stmt.Close()
	for stmt.Next() {
		var row EmailMeta
		// for each row, scan the result into our composite object
		err = stmt.Scan(&row.ID, &row.MsgID, &row.ThdID)
		errorHandler(err)
		unprocessed = append(unprocessed, row)
	}
	return unprocessed
}

func ListByThdID(thdID string) []EmailMeta {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	var thdList []EmailMeta

	stmt, err := tx.Query("SELECT ID, msgID, thdID FROM emails WHERE thdID = ?", thdID)
	errorHandler(err)
	defer stmt.Close()
	for stmt.Next() {
		var row EmailMeta
		// for each row, scan the result into our composite object
		err = stmt.Scan(&row.ID, &row.MsgID, &row.ThdID)
		errorHandler(err)
		thdList = append(thdList, row)
	}
	return thdList
}

func MarkEmailCompleate(ID int) {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	stmt, err := tx.Prepare("UPDATE emails SET proccessed=1 WHERE id=?")
	errorHandler(err)

	res, err := stmt.Exec(ID)
	errorHandler(err)

	_, err = res.RowsAffected()
	errorHandler(err)

	// fmt.Println(affect)
}
