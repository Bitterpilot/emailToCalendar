package models

// Event is a basic type of an calendar event
type Event struct {
	EventID     string
	Summary     string
	Start       string
	End         string
	Timezone    string
	Location    string
	Description string
	MsgID       string
	Link        string
	Processed   bool
	ProcessedTime int64
}

// RowContent is what is expected from a riteq schedule email
// and is used as a bridging type between before being made into an event
type RowContent struct {
	Day        string
	Date       string
	StartWork  string
	EndWork    string
	TotalHours string
	Breaks     string
	Pay        string
	OrgLevel   string
	Comments   string
}
