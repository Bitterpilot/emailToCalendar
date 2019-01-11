package db

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

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
}

func errorHandler(err error) {
	var tx *sql.Tx
	if err != nil {
		log.Printf("%v\n\n\n", errors.Cause(err))
		tx.Rollback()
	}
}
