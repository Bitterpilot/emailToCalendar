package main

type tableData struct {
	day        string
	date       string
	startWork  string // gcal api expects strings
	endWork    string // https://developers.google.com/calendar/create-events#metadata
	totalHours string
	breaks     string
	pay        string
	orgLevel   string
}

type shift struct {
	summary        string // will be the procesed orgLevel (remove everything between \ inclusive) ***REMOVED***
	location       string // Derived from orgLevel (if item before \ = AAAA then 303 if item before \ = BBBB)
	eventDateStart string // date + startWork
	eventDateEnd   string // date + endWork
	processed      bool   // true/false/nil
}

// https://developers.google.com/calendar/extended-properties

func main() {

}
