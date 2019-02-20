package processor

import (
	"fmt"
	"regexp"
	"time"
)

// RowContents is what is expected from a riteq schedule email
type RowContents struct {
	Day        string
	Date       string
	StartWork  string // google cal api expects strings
	EndWork    string // https://developers.google.com/calendar/create-events#metadata
	TotalHours string
	Breaks     string
	Pay        string
	OrgLevel   string
	Comments   string
}

// Shift ...
// https://developers.google.com/calendar/extended-properties
type Shift struct {
	Summary string // will be the processed orgLevel (remove everything between \ inclusive) https://regexr.com/46729
	//location string // Derived from orgLevel (if item before \ = AAAA then 303 if item before \ = BBBB)
	EventDateStart string // date + startWork
	EventDateEnd   string // date + endWork
	MsgID          string
	Processed      bool // true/false/nil
}

// ProcessRows takes a slice of strings and applies them to the RowContents struct
func ProcessRows(content []string) []RowContents {
	nContent := []RowContents{}
	// check where the slice needs to be divided
	position := 0
	for key, val := range content {
		if val == "Comments" {
			position = key + 1
		} else if val == "Org Level" {
			position = key + 1
		}
	}
	// fmt.Println("*********")
	start := 0
	end := position
	for i := 1; i <= (len(content) / position); i++ {
		// fmt.Println(content[start:end])
		rowLoc := content[start:end]
		rowStruct := RowContents{}
		if len(rowLoc) == 9 {
			rowStruct = RowContents{rowLoc[0], rowLoc[1], rowLoc[2], rowLoc[3], rowLoc[4], rowLoc[5], rowLoc[6], rowLoc[7], rowLoc[8]}
		} else {
			rowStruct = RowContents{rowLoc[0], rowLoc[1], rowLoc[2], rowLoc[3], rowLoc[4], rowLoc[5], rowLoc[6], rowLoc[7], ""}
		}
		nContent = append(nContent, rowStruct)
		start = end
		end = end + position
	}
	return nContent
}

// ProcessShift skip the 0th item when iterating
func ProcessShift(year []string, row RowContents, msgID string) Shift {
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
	// fmt.Println(Start)
	dateStart, err := time.Parse(
		time.RFC822Z,
		Start)
	if err != nil {
		fmt.Println(err)
	}
	End := date + " " + row.EndWork + " " + "+0800"
	// fmt.Println(End)
	dateEnd, err := time.Parse(
		time.RFC822Z,
		End)
	if err != nil {
		fmt.Println(err)
	}
	eventDateStart := dateStart.Format(time.RFC3339)
	eventDateEnd := dateEnd.Format(time.RFC3339)

	r := regexp.MustCompile(`\\(.*?)\\`)
	summary := r.ReplaceAllString(row.OrgLevel, " ")
	// will be the processed orgLevel (remove everything between \ inclusive) https://regexr.com/46729
	// 	eventDateStart := "" // date + startWork
	// 	eventDateEnd := ""   // date + endWork
	// 	processed := false   // true/false/nil

	shift := Shift{summary, eventDateStart, eventDateEnd, msgID, false}
	return shift
}
