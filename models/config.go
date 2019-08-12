package models

// Config holds settings
type Config struct {
	User       string
	Label      string
	Sender     string
	Subject    string
	CalendarID string
	TimeZone   string
	Locations  map[string]string
}
