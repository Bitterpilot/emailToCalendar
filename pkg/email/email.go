package email

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// ReadDate reads the time that the shift email covers
func (msg *Msg) ReadDate() (map[string]time.Time, error) {
	decoded, err := base64.StdEncoding.DecodeString(msg.Body)
	if err != nil {
		fmt.Println(err)
	}
	sDecode := string(decoded)
	re := regexp.MustCompile(`([\d]{1,2} [\w]{3} [\d]{4})`)
	match := re.FindAllString(sDecode, -1) // FindAllString(t,n) n = max number of matches
	// verify that there is only two dates
	if len(match) != 2 {
		e := fmt.Sprintf("Cannot parse body text to get roster date range\n Expected 2 got %d\n%v", len(match), match)
		return nil, errors.New(e)
	}

	dates := make(map[string]time.Time)
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

// ReadTable reads the table in an email
func (msg *Msg) ReadTable() ([]map[string]string, error) {
	decoded, _ := base64.URLEncoding.DecodeString(msg.Body)
	sDecode := string(decoded)
	table := readTag(sDecode, "td", "table")
	rows := processRows(table)
	return rows, nil
}

func processRows(content []string) []map[string]string {
	var nContent []map[string]string
	// check where the slice needs to be devided
	var positions []int
	for key, val := range content {
		switch {
		case val == "Day",
			val == "Mon",
			val == "Tue",
			val == "Wed",
			val == "Thu",
			val == "Fri",
			val == "Sat",
			val == "Sun":
			positions = append(positions, key)
		}
	}
	positions = append(positions, len(content))

	for i := 0; i < len(positions)-1; i++ {
		start := positions[i]
		end := positions[i+1]

		rowLoc := content[start:end]
		row := make(map[string]string)
		switch {
		case len(rowLoc) == 9:
			row["Day"] = rowLoc[0]
			row["Date"] = rowLoc[1]
			row["StartWork"] = rowLoc[2]
			row["EndWork"] = rowLoc[3]
			row["TotalHours"] = rowLoc[4]
			row["Breaks"] = rowLoc[5]
			row["Pay"] = rowLoc[6]
			row["OrgLevel"] = rowLoc[7]
			row["Comments"] = rowLoc[8]
		case len(rowLoc) == 8:
			row["Day"] = rowLoc[0]
			row["Date"] = rowLoc[1]
			row["StartWork"] = rowLoc[2]
			row["EndWork"] = rowLoc[3]
			row["TotalHours"] = rowLoc[4]
			row["Breaks"] = rowLoc[5]
			row["Pay"] = rowLoc[6]
			row["OrgLevel"] = rowLoc[7]
		default:
			continue
		}
		nContent = append(nContent, row)
	}

	return nContent
}

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
