package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/jhillyerd/enmime"

	googlecal "github.com/bitterpilot/emailtocal/calendar"
	g "github.com/bitterpilot/emailtocal/gmail"
	processor "github.com/bitterpilot/emailtocal/shift"
)

func readEmail(v []byte) (year []string, table []processor.RowContents) {
	eml := string(v[:])
	year = readYear(eml)
	table = processTable(eml)

	return year, table
}

func decode(msg string) (body string, err error) {
	// Open a sample message file.
	r, err := os.Open(msg)
	if err != nil {
		return "", err
	}

	// Parse message body with enmime.
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return "", err
	}

	// fmt.Println(reflect.TypeOf(env.HTML)) = string

	return env.HTML, nil
}

// readTag takes the html body and reads the contents of the readTag
// insome cases you will want this process to end at a differnt endTag
// eg read all table rows <td> until you reach </table>
func readTag(body, readTag, endTag string) (table []string) {
	tokenizer := html.NewTokenizer(strings.NewReader(body))
	content := []string{}

	// While have not hit the </endTag> tag
	for tokenizer.Token().Data != endTag {
		tocNext := tokenizer.Next()
		if tocNext == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == readTag {
				inner := tokenizer.Next()
				if inner == html.TextToken {
					text := strings.TrimSpace((string)(tokenizer.Text()))
					content = append(content, text)
				}
			}
		}
	}
	return content
}

func readYear(eml string) []string {
	t := readTag(eml, "p", "p")
	parts := strings.Split(t[0], " ")
	years := []string{}

	// select Year values from first line
	for _, val := range parts {
		match, _ := regexp.MatchString("([0-9]{4})", val)
		if match == true {
			// TODO: fix this before 2100
			val = strings.TrimPrefix(val, "20")
			years = append(years, val)
		}
	}
	return years
}

func processTable(eml string) []processor.RowContents {
	table := readTag(eml, "td", "table")
	// days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	// 8 - 15
	// 16 - 2
	rows := processor.ProcessRows(table)
	return rows
}

func main() {
	// use this to look for new messages
	/*
		listMessages := g.ListMessages("Label_24", "CCCCCCCC@riteq.com.au", "Schedule for DDDDDDDD")
		for key, val := range listMessages {
			fmt.Printf("------\nitem: %d\nmsgID: %s\nthdID: %s\n", key, val.Id, val.ThreadId)
}
	*/
	// look up specific message
	msgID := "***REMOVED***"
	user := "me"
	_, _, body := g.GetMessage(user, msgID)
	// fmt.Println("*** Specific Message ***")
	// fmt.Printf("msgID:%s thread:%s \nrecieved(unix timestamp):%d\nbody:\n%s\n", msgID, threadID, date, body)
	// fmt.Printf("%s", body)
	// TODO: stream html table in to db
	// TODO: read db where processed = false and insert into cal
	year, rows := readEmail(body)
	shifts := []processor.Shift{}
	// range over all rows except the hearder row (row 0)
	for _, row := range rows[1:] {
		shift := processor.ProcessShift(year, row)
		shifts = append(shifts, shift)
	}

	// Start of calandar stuff
	for _, shift := range shifts {
		calendarID := "***REMOVED***"
		summary := shift.Summary
		// description needs html formating
		processTime := time.Now().Format(time.RFC822) // more format options https://golang.org/pkg/time/#pkg-constants
		description := fmt.Sprintf(`Automatically created by emailToCal at %s<br><a href="https://mail.google.com/mail/#inbox/%s">Source</a>`, processTime, msgID)
		timezone := "Australia/Perth"
		dateTimeStart := shift.EventDateStart
		dateTimeEnd := shift.EventDateEnd

		googlecal.AddEvent(calendarID, summary, msgID, description, timezone, dateTimeStart, dateTimeEnd)
	}

}
