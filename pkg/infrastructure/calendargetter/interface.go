package calendargetter

import (
	"google.golang.org/api/calendar/v3"
)

// CalendarProvider
type CalendarProvider struct {
	service *calendar.Service
	calID   string
}

// NewCalendarProvider
func NewCalendarProvider(ID string) *CalendarProvider {
	return &CalendarProvider{
		service: newService(),
		calID:   ID,
	}
}
