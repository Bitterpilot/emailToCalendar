package usecase

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

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
