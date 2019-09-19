package app

import (
	"time"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"github.com/bitterpilot/emailToCalendar/models"
)

func findMonday(t time.Time) time.Time {
	// sunday = 0
	// monday = 1
	diff := 1 - int(t.Weekday())
	return time.Date(t.Year(), t.Month(), t.Day()+diff, 0, 0, 0, 0, t.Location())
}

// Check6in7 checks for 6 working days in the calendar week
// and logs a warning.
func Check6in7(cr *CalendarRegistar, m *models.Email) {
	// get first day of fornight
	s, err := time.Parse(time.RFC3339, m.List[0].Start)
	if err != nil {
		log.Errorf("Check6in7: Parse Time: %v", err)
	}

	// divide into weeks
	start := findMonday(s)
	mid := time.Date(start.Year(), start.Month(), start.Day()+7, 0, 0, 0, 0, start.Location())
	end := time.Date(start.Year(), start.Month(), start.Day()+14, 0, 0, 0, 0, start.Location())

	// get evens from weeks
	week1, err := cr.CalendarStore.ListEventsByDateRange(start.Format(time.RFC3339), mid.Format(time.RFC3339))
	if err != nil {
		log.Errorf("Check6in7: week1 ListEventsByDateRange: %v\n", err)
	}
	week2, err := cr.CalendarStore.ListEventsByDateRange(mid.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		log.Errorf("Check6in7: week2 ListEventsByDateRange: %v\n", err)
	}

	check(start, week1)
	check(mid, week2)
}

func check(start time.Time, week []models.Event) {
	if len(week) >= 6 {
		log.Warnf("Week of the %s has %d shifts!", start.Format(time.RFC822), len(week))
		c := count(week)
		log.Warnf(`
		+---+---+---+---+---+---+---+
		| M | T | W | T | F | S | S |
		| %d | %d | %d | %d | %d | %d | %d |
		+---+---+---+---+---+---+---+`, c["m"], c["tu"], c["w"], c["te"], c["f"], c["sa"], c["su"])

		log.Tracef("%v\n", pretty.Formatter(week))
	}
}

func count(es []models.Event) map[string]int {
	ret := map[string]int{"m": 0, "tu": 0, "w": 0, "te": 0, "f": 0, "sa": 0, "su": 0}
	for _, e := range es {
		t, err := time.Parse(time.RFC3339, e.Start)
		if err != nil {
			log.Errorln(err)
		}
		switch t.Weekday() {
		case time.Sunday:
			ret["su"]++
		case time.Monday:
			ret["m"]++
		case time.Tuesday:
			ret["tu"]++
		case time.Wednesday:
			ret["w"]++
		case time.Thursday:
			ret["te"]++
		case time.Friday:
			ret["f"]++
		case time.Saturday:
			ret["sa"]++
		}
	}
	return ret
}
