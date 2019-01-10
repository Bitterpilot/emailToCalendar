package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

// func init() {
// // checks if the db file exists
// if _, err := os.Stat("./foo.db"); os.IsNotExist(err) {
// 	panic("db file doesn't exist")
// }
// }
var db *sql.DB

func init() {
	// os.Remove("./db/foo.db")
	var err error
	db, err = sql.Open("sqlite3", "./db/foo.db")
	errorHandler(err)

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS emails (
		ID INTEGER NOT NULL PRIMARY KEY,
		msgID text,
		thdID text,
		timeRecieved INTEGER,
		proccessed BOOLEAN DEFAULT 0,
		error BOOLEAN DEFAULT 0
		);
		`
	_, err = db.Exec(sqlStmt)
	errorHandler(err)
	sqlStmt = `
		CREATE TABLE IF NOT EXISTS shifts (
			ID INTEGER NOT NULL PRIMARY KEY,
			Summery text,
			description text,
			TimeZone text,
			EventDateStart text,
			EventDateEnd text,
			Processed BOOLEAN DEFAULT FALSE,
			proccessTime text,
			deleted BOOLEAN DEFAULT FALSE,
			deletedTime text,
			eventID text,
			msgID int NOT NULL,
			FOREIGN KEY (msgID) REFERENCES emails ("msgID")
			);`
	_, err = db.Exec(sqlStmt)
	errorHandler(err)
	// sqlStmt = `
	// 		UPDATE emails SET proccessed=0 WHERE thdID="***REMOVED***";
	// 		UPDATE emails SET proccessed=0 WHERE thdID="***REMOVED***";`
	// _, err = db.Exec(sqlStmt)
	// errorHandler(err)

}

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

func InsertShift(Summery, description, TimeZone, EventDateStart, EventDateEnd, Processed, proccessTime, eventID, msgID string) {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	stmt, err := tx.Prepare(`
		INSERT INTO "shifts" ("Summery", "description", "TimeZone", "EventDateStart", "EventDateEnd", "Processed", "proccessTime","eventID", "msgID")
    	VALUES ( ?, ?, ?, ?, ?, ?, ?,?,?);
		`)
	errorHandler(err)
	defer stmt.Close()
	_, err = stmt.Exec(Summery, description, TimeZone, EventDateStart, EventDateEnd, Processed, proccessTime, eventID, msgID)
	errorHandler(err)
}

func ListEventIDByEmailID(msgID string) []string {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	stmt, err := tx.Query("SELECT eventID FROM shifts WHERE msgID = ? AND deleted = 0", msgID)
	errorHandler(err)
	defer stmt.Close()
	var eventList []string
	for stmt.Next() {
		var row string
		err = stmt.Scan(&row)
		errorHandler(err)
		eventList = append(eventList, row)
	}
	return eventList
}

func MarkShiftAsDeleted(eventID string) {
	tx, err := db.Begin()
	errorHandler(err)
	defer tx.Commit()

	time := time.Now().String()
	_, err = tx.Query("UPDATE shifts SET deleted=1, deletedTime=? WHERE id=?;", time, eventID)
	errorHandler(err)
}

func errorHandler(err error) {
	var tx *sql.Tx
	if err != nil {
		log.Printf("%v\n\n\n", errors.Cause(err))
		tx.Rollback()
	}
}
