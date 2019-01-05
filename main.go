package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/jhillyerd/enmime"

	g "github.com/bitterpilot/emailtocal/gmail"
)

func readEmail(v string) (year, table []string) {
	eml, _ := decode(v)
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
	z := html.NewTokenizer(strings.NewReader(body))
	content := []string{}

	// While have not hit the </endTag> tag
	for z.Token().Data != endTag {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == readTag {
				inner := z.Next()
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := strings.TrimSpace(text)
					content = append(content, t)
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
	for _, v := range parts {
		match, _ := regexp.MatchString("([0-9]{4})", v)
		if match == true {
			years = append(years, v)
		}
	}
	return years
}

func processTable(eml string) []string {
	table := readTag(eml, "td", "table")
	// days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	// 8 - 15
	// 16 - 2
	table2 := []string{}
	for key := range table {
		modlo := key % 8
		if modlo == 0 {
			table2 = table[key:(key + 8)]
		}
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
	fmt.Printf("%s", body)
	// TODO: stream html table in to db
	// TODO: read db where processed = false and insert into cal
	// TODO: FIXME: msgID with an extra column(comments) "***REMOVED***"
	// struct or map for each table row all kept in an array
	// []map{tablerow, tablerow.....}
	// tablerow is described in shift/shift.go

	// Start of calandar stuff
	/*
		calendarID := "***REMOVED***"
		summary := "AAAA CSO"
		messageID := "***REMOVED***"
		// description needs html formating
		processTime := time.Now().Format(time.RFC822) // more format options https://golang.org/pkg/time/#pkg-constants
		description := fmt.Sprintf(`Automatically created by emailToCal at %s<br><a href="https://mail.google.com/mail/#inbox/%s">Source</a>`, processTime, messageID)
		timezone := "Australia/Perth"
		dateTimeStart := "2019-01-01T09:00:00+08:00"
		dateTimeEnd := "2019-01-01T17:00:00+08:00"

		googlecal.AddEvent(calendarID, summary, messageID, description, timezone, dateTimeStart, dateTimeEnd)
	*/
}
