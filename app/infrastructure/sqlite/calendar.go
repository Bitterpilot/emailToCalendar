package sqlite

import "time"

// InsertShift
func (s CalendarStore) InsertShift(Summery, description, TimeZone, EventDateStart, EventDateEnd, Processed, proccessTime, eventID, msgID string) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(`
		INSERT INTO "shifts" ("Summery", "description", "TimeZone", "EventDateStart", "EventDateEnd", "Processed", "proccessTime","eventID", "msgID")
    	VALUES ( ?, ?, ?, ?, ?, ?, ?,?,?);
		`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(Summery, description, TimeZone, EventDateStart, EventDateEnd, Processed, proccessTime, eventID, msgID)
	if err != nil {
		return err
	}
}

// ListEventIDByEmailID
func (s CalendarStore) ListEventIDByEmailID(msgID string) []string {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	stmt, err := tx.Query("SELECT eventID FROM shifts WHERE msgID = ? AND deleted = 0", msgID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var eventList []string
	for stmt.Next() {
		var row string
		err = stmt.Scan(&row)
		if err != nil {
			return err
		}
		eventList = append(eventList, row)
	}
	return eventList
}

// MarkShiftAsDeleted
func (s CalendarStore) MarkShiftAsDeleted(eventID string) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	time := time.Now().String()
	_, err = tx.Query("UPDATE shifts SET deleted=1, deletedTime=? WHERE id=?;", time, eventID)
	if err != nil {
		return err
	}
}
