package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/jhillyerd/enmime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	cal "github.com/bitterpilot/emailToCalendar/calendar"
	db "github.com/bitterpilot/emailToCalendar/db"
	g "github.com/bitterpilot/emailToCalendar/gmail"
	processor "github.com/bitterpilot/emailToCalendar/shift"
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

// readYear originally returned years as 2 digit ie,[18 19] so it could be used
// more easily in the date format RFC822Z. ie("02 Jan 06 15:04 -0700").
// it now returns a map with time.Time to allow for easier manipulation later.
func readYear(eml string) (map[string]time.Time, error) {
	re := regexp.MustCompile(`([\d]{1,2} [\w]{3} [\d]{4})`)
	match := re.FindAllString(eml, -1) // FindAllString(t,n) n = max number of matches
	// verify that there is only two dates
	if len(match) != 2 {
		msg := fmt.Sprintf("Cannot parse body text to get roster date range\n Expected 2 got %d", len(match))
		return nil, errors.New(msg)
	}

	var dates map[string]time.Time
	loc, _ := time.LoadLocation("Australia/Perth")
	for i, val := range match {
		layout := "2 Jan 2006"
		ti, err := time.ParseInLocation(layout, val, loc)
		if err != nil {
			fmt.Println(err)
		}
		switch {
		case i == 0:
			dates["Start"] = ti
		case i == 1:
			dates["End"] = ti
		}
	}
	return dates, nil
}

func processTable(eml string) []processor.RowContents {
	table := readTag(eml, "td", "table")
	// days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	// 8 - 15
	// 16 - 2
	rows := processor.ProcessRows(table)
	return rows
}

func publishShifts(shifts []processor.Shift) {
	calendarID := viper.GetString("calendarID")
	// Start of calandar stuff
	for _, shift := range shifts {
		msgID := shift.MsgID
		summary := shift.Summary
		// description needs html formating
		processTime := time.Now().Format(time.RFC822) // more format options https://golang.org/pkg/time/#pkg-constants
		description := fmt.Sprintf(`Automatically created by emailToCalendar at %s<br><a href="https://mail.google.com/mail/#inbox/%s">Source</a>`, processTime, msgID)
		timezone := "Australia/Perth"
		dateTimeStart := shift.EventDateStart
		dateTimeEnd := shift.EventDateEnd
		eventID := cal.AddEvent(calendarID, summary, msgID, description, timezone, dateTimeStart, dateTimeEnd)
		db.InsertShift(summary, description, timezone, dateTimeStart, dateTimeEnd, "1", time.Now().String(), eventID, msgID)
	}
}

// func error()  {
// Func err
// Write fail to db
// Create event at current day for notifications
// Write details to log
// }
var reader *bufio.Reader

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	user := "me"
	// load user info
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("Config file not found: %v", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Config file was found with error: %v", err)
		}
	}

	label := viper.GetString("label")
	sender := viper.GetString("sender")
	subject := viper.GetString("subject")

	// use this to look for new messages
	listMessages := g.ListMessages(label, sender, subject)
	for _, val := range listMessages {
		if val.Id != db.ListByMsgID(val.Id) {
			_, date, _ := g.GetMessage(user, val.Id)
			db.InsertEmail(val.Id, val.ThreadId, date)
		}
	}
	unprocessed := db.ListUnprocssed()
	fmt.Println(unprocessed)

	shifts := []processor.Shift{}
	for _, val := range unprocessed {
		if val.MsgID == "***REMOVED***" {
			break
		}
		if len(db.ListByThdID(val.ThdID)) > 1 {
			fmt.Println(val.MsgID)
		}
		// look up specific message
		msgID := val.MsgID
		_, _, body := g.GetMessage(user, msgID)
		// fmt.Println("*** Specific Message ***")
		// fmt.Printf("msgID:%s thread:%s \nrecieved(unix timestamp):%d\nbody:\n%s\n", msgID, threadID, date, body)
		// fmt.Printf("%s", body)
		year, rows := readEmail(body)
		// range over all rows except the hearder row (row 0)
		for _, row := range rows[1:] {
			shift := processor.ProcessShift(year, row, msgID)
			shifts = append(shifts, shift)
		}
		db.MarkEmailCompleate(val.ID)
	}
	publishShifts(shifts)

}
