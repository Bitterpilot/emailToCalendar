package app

import (
	"regexp"
	"time"

	"github.com/bitterpilot/emailToCalendar/models"
	log "github.com/sirupsen/logrus"
)

// Application logic related to retrieving and processing Events

// CalendarGetter are functions a calendar provider( Gmail, Outlook ) package will implement.
type CalendarGetter interface {
	Create(models.Event) (models.Event, error)
	List() ([]models.Event, error)
	Get(models.Event) (models.Event, error)
	Update(models.Event) (models.Event, error)
	Delete(models.Event) error
}

// CalendarStore are functions a database package will implement.
type CalendarStore interface {
	FindByMsgID(string) (string, error)
	FindByEventID(string) (string, error)
	InsertEvent(models.Event) (int, error)
	AlterEvent(models.Event) error
}

// CalendarRegistar contains dependencies for this package. Such as an calendar provider, database and logger.
type CalendarRegistar struct {
	// calendar getter dependencies
	calendarGetter CalendarGetter
	calendarStore  CalendarStore
}

// NewCalendarRegistar loads the dependencies declared in main() into this package for use.
func NewCalendarRegistar(e CalendarGetter, db CalendarStore) *CalendarRegistar {
	return &CalendarRegistar{
		calendarGetter: e,
		calendarStore:  db,
	}
}

// ProcessShift takes a slice with years and RowContent and creates an Event object.
// The 0th item is the table header.
func processShift(year []string, row models.RowContent, msgID string) models.Event {
	date := ""
	if year[0] == year[1] {
		date = row.Date + " " + year[0]
	} else {
		r := regexp.MustCompile(`\d{2}\s`)
		month := r.ReplaceAllString(row.Date, "")
		switch {
		case month == "Dec":
			date = row.Date + " " + year[0]
		case month == "Jan":
			date = row.Date + " " + year[1]
		}
	}
	Start := date + " " + row.StartWork + " " + "+0800"
	log.Debugln(Start)
	dateStart, err := time.Parse(
		time.RFC822Z,
		Start)
	if err != nil {
		log.Errorln(err)
	}
	End := date + " " + row.EndWork + " " + "+0800"
	log.Debugln(End)
	dateEnd, err := time.Parse(
		time.RFC822Z,
		End)
	if err != nil {
		log.Errorln(err)
	}
	eventDateStart := dateStart.Format(time.RFC3339)
	eventDateEnd := dateEnd.Format(time.RFC3339)

	r := regexp.MustCompile(`\\(.*?)\\`)
	summary := r.ReplaceAllString(row.OrgLevel, " ")
	// will be the processed orgLevel (remove everything between \ inclusive) https://regexr.com/46729
	// 	eventDateStart := "" // date + startWork
	// 	eventDateEnd := ""   // date + endWork
	// 	processed := false   // true/false/nil

	shift := models.Event{Summary: summary, Start: eventDateStart, End: eventDateEnd, MsgID: msgID, Processed: false}
	return shift
}
