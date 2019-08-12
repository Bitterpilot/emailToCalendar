package app

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/bitterpilot/emailToCalendar/models"
	"golang.org/x/net/html"
)

// ReadBody tokenizes the email's html and pulls out the year and the table
func ReadBody(e models.Email) ([]string, []models.RowContent, error) {
	tokenizer := html.NewTokenizer(bytes.NewReader(e.Body))
	year, err := readYear(tokenizer)
	if err != nil {
		return nil, nil, err
	}
	table := processTable(tokenizer)

	return year, table, nil
}

func readYear(t *html.Tokenizer) ([]string, error) {
	d := readTag(t, "p", "p")
	parts := strings.Split(d[0], " ")
	years := []string{}

	// select Year values from first line
	for _, val := range parts {
		match, err := regexp.MatchString("([0-9]{4})", val)
		if err != nil {
			return nil, err
		}
		if match == true {
			// TODO: fix this before 2100
			val = strings.TrimPrefix(val, "20")
			years = append(years, val)
		}
	}
	return years, nil
}

func processTable(t *html.Tokenizer) []models.RowContent {
	return processRows(readTag(t, "td", "table"))
}

// processRows takes a slice of strings and applies them to the models.RowContent struct
func processRows(content []string) []models.RowContent {
	nContent := []models.RowContent{}
	// check where the slice needs to be divided
	var positions []int
	for key, val := range content {
		switch {
		case val == "Day", val == "Mon", val == "Tue", val == "Wed", val == "Thu", val == "Fri", val == "Sat", val == "Sun":
			positions = append(positions, key)
		}
	}
	positions = append(positions, len(content))

	for i := 0; i < len(positions)-1; i++ {
		start := positions[i]
		end := positions[i+1]

		rowLoc := content[start:end]
		rowStruct := models.RowContent{}
		switch {
		case len(rowLoc) == 9:
			rowStruct = models.RowContent{
				Day:        rowLoc[0],
				Date:       rowLoc[1],
				StartWork:  rowLoc[2],
				EndWork:    rowLoc[3],
				TotalHours: rowLoc[4],
				Breaks:     rowLoc[5],
				Pay:        rowLoc[6],
				OrgLevel:   rowLoc[7],
				Comments:   rowLoc[8],
			}
		case len(rowLoc) == 8:
			rowStruct = models.RowContent{
				Day:        rowLoc[0],
				Date:       rowLoc[1],
				StartWork:  rowLoc[2],
				EndWork:    rowLoc[3],
				TotalHours: rowLoc[4],
				Breaks:     rowLoc[5],
				Pay:        rowLoc[6],
				OrgLevel:   rowLoc[7],
			}
		default:
			continue
		}
		nContent = append(nContent, rowStruct)
	}

	return nContent
}

func readTag(tokenizer *html.Tokenizer, readTag string, endTag string) (table []string) {
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
