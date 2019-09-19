package sqlite

import (
	"time"

	"github.com/bitterpilot/emailToCalendar/models"
)

// InsertShift into a sqlite database.
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

// ListEventIDByEmailID list events originating from an email.
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

// ListEventsByDateRange returns a slice of events
func (s CalendarStore) ListEventsByDateRange(begin, end string) ([]models.Event, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	stmt, err := tx.Query("SELECT eventID, EventDateStart FROM shifts WHERE EventDateStart BETWEEN ? AND ? AND deleted=0", begin, end)
	if err != nil {
		return nil, err
	}

	var ret []models.Event
	for stmt.Next() {
		var row models.Event
		err = stmt.Scan(&row.EventID, &row.Start)
		if err != nil {
			return nil, err
		}
		ret = append(ret, row)
	}
	return ret, nil
}

// MarkShiftAsDeleted marks an event as deleted from a service(ie; Google Calendar).
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
