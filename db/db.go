package db

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./db/foo.db")
	if err != nil {
		log.Println(err)
	}

	tx, err := db.Begin()
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
	_, err = tx.Exec(sqlStmt)
	errorHandler(err, tx)
	tx.Commit()

	tx, err = db.Begin()
	sqlStmt = `
	CREATE TABLE IF NOT EXISTS user (
		ID INTEGER NOT NULL PRIMARY KEY,
		username text,
		userQLable text,
		userQEmail text,
		userQSubject text
		);
		`
	_, err = tx.Exec(sqlStmt)
	errorHandler(err, tx)
	tx.Commit()

	tx, err = db.Begin()
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
	_, err = tx.Exec(sqlStmt)
	errorHandler(err, tx)
	tx.Commit()
}

func errorHandler(err error, tx *sql.Tx) {
	if err != nil {
		log.Printf("%v\n\n\n", errors.Cause(err))
		tx.Rollback()
	}
}
