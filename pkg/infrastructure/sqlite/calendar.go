package sqlite

import (
	"time"

	"github.com/bitterpilot/emailToCalendar/models"
)

// InsertShift
func (s CalendarStore) InsertShift(e models.Event) error {
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

	res, err := stmt.Exec(e.Summary, e.Description, e.Timezone, e.Start, e.End, e.Processed, e.ProcessedTime, e.EventID, e.MsgID)
	if err != nil {
		return err
	}
	res.LastInsertId()

	return nil
}

// ListEventIDByEmailID
func (s CalendarStore) ListEventIDByEmailID(msgID string) ([]string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	stmt, err := tx.Query("SELECT eventID FROM shifts WHERE msgID = ? AND deleted = 0", msgID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var eventList []string
	for stmt.Next() {
		var row string
		err = stmt.Scan(&row)
		if err != nil {
			return nil, err
		}
		eventList = append(eventList, row)
	}
	return eventList, nil
}

// MarkShiftAsDeleted
func (s CalendarStore) MarkShiftAsDeleted(eventID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	time := time.Now().Unix()
	_, err = tx.Query("UPDATE shifts SET deleted=1, deletedTime=? WHERE id=?;", time, eventID)
	if err != nil {
		return err
	}
	return nil
}
