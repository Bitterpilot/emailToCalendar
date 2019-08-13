package calendargetter

import (
	"google.golang.org/api/calendar/v3"
)

// CalendarProvider connects the google calendar service to our program.
type CalendarProvider struct {
	service *calendar.Service
	calID   string
}

// NewCalendarProvider creates a new instance of CalendarProvider.
func NewCalendarProvider(ID string) *CalendarProvider {
	return &CalendarProvider{
		service: newService(),
		calID:   ID,
	}
}
