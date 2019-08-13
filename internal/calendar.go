package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bitterpilot/emailToCalendar/models"
	log "github.com/sirupsen/logrus"
)

// Application logic related to retrieving and processing Events

// CalendarGetter are functions a calendar provider( Google Calendar, Outlook ) package will implement.
type CalendarGetter interface {
	Create(models.Event) (models.Event, error)
	List() ([]models.Event, error)
	Get(models.Event) (models.Event, error)
	Update(models.Event) (models.Event, error)
	Delete(models.Event) error
	Validate(e models.Event) error
}

// CalendarStore are functions a database package will implement.
type CalendarStore interface {
	InsertShift(e models.Event) error
	ListEventIDByEmailID(msgID string) ([]string, error)
	MarkShiftAsDeleted(eventID string) error
}

// CalendarRegistar contains dependencies for this package. Such as an calendar provider, database and logger.
type CalendarRegistar struct {
	// calendar getter dependencies
	CalendarGetter CalendarGetter
	CalendarStore  CalendarStore
	Config         *models.Config
}

// NewCalendarRegistar loads the dependencies declared in main() into this package for use.
func NewCalendarRegistar(e CalendarGetter, db CalendarStore, c *models.Config) *CalendarRegistar {
	return &CalendarRegistar{
		CalendarGetter: e,
		CalendarStore:  db,
		Config:         c,
	}
}

// Publish
func (r *CalendarRegistar) Publish(event models.Event) (models.Event, error) {
	processTime := time.Now()

	// Description
	// This uses the process time variable witch we want as close to uploading the event as possible
	event.Description = fmt.Sprintf(`Automatically created by emailToCalendar at %s<br><a href="https://mail.google.com/mail/#inbox/%s">Source</a>`,
		processTime.Format(time.RFC822), event.MsgID)

	// location
	for k, v := range r.Config.Locations {
		if strings.Contains(event.Summary, strings.ToUpper(k)) {
			event.Location = v
			break
		}
	}

	if err := r.CalendarGetter.Validate(event); err != nil {
		return event, err
	}

	event, err := r.CalendarGetter.Create(event)
	if err != nil {
		return models.Event{}, err
	}

	event.Processed = true
	event.ProcessedTime = processTime.Unix()

	if err := r.CalendarStore.InsertShift(event); err != nil {
		return models.Event{}, err
	}

	return event, nil
}

// BuildEvent takes a slice with years and RowContent and creates an Event object.
// The 0th item is the table header.
func (r *CalendarRegistar) BuildEvent(years []string, row models.RowContent, msgID string) (models.Event, error) {
	// date
	date := chooseYear(years, row)

	start, err := buildDate(date, row.StartWork)
	if err != nil {
		return models.Event{}, err
	}
	end, err := buildDate(date, row.EndWork)
	if err != nil {
		return models.Event{}, err
	}

	eventDateStart := start.Format(time.RFC3339)
	eventDateEnd := end.Format(time.RFC3339)

	reg := regexp.MustCompile(`\\(.*?)\\`)
	summary := reg.ReplaceAllString(row.OrgLevel, " ")
	// will be the processed orgLevel (remove everything between \ inclusive) https://regexr.com/46729
	// 	eventDateStart := "" // date + startWork
	// 	eventDateEnd := ""   // date + endWork
	// 	processed := false   // true/false/nil

	// Location
	var location string

	// Build shift
	shift := models.Event{
		Summary:  summary,
		Start:    eventDateStart,
		End:      eventDateEnd,
		Timezone: r.Config.TimeZone,
		Location: location,
		MsgID:    msgID,
	}
	return shift, nil
}

func chooseYear(years []string, row models.RowContent) (year string) {
	if years[0] == years[1] {
		year = row.Date + " " + years[0]
	} else {
		r := regexp.MustCompile(`\d{2}\s`)
		month := r.ReplaceAllString(row.Date, "")
		switch {
		case month == "Dec":
			year = row.Date + " " + years[0]
		case month == "Jan":
			year = row.Date + " " + years[1]
		}
	}

	return year
}

func buildDate(year string, day string) (time.Time, error) {
	Start := year + " " + day + " " + "+0800"
	date, err := time.Parse(time.RFC822Z, Start)
	if err != nil {
		log.Errorln(err)
	}

	return date, err
}
