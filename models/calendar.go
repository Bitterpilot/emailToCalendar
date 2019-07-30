package models

import "time"

// Event is a basic type of an event
type Event struct {
	Title      string
	Start      *time.Time
	End        *time.Time
	Timezone   string
	Location   string
	Descrption string
}
