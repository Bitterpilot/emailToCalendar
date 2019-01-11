package db

import "time"

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
