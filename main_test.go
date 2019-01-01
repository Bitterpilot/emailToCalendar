package main

import (
	"fmt"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	// test read email outputs the correct years in []string
	tables := []struct {
		file string
		ans  []string
	}{
		{"casesTest/01.eml", []string{"2018", "2018"}},
		{"casesTest/02.eml", []string{"2018", "2019"}},
		{"casesTest/03.eml", []string{"2018", "2018"}},
		{"casesTest/04.eml", []string{"2018", "2019"}},
		{"casesTest/05.eml", []string{"2018", "2018"}},
		{"casesTest/06.eml", []string{"2018", "2018"}},
		{"casesTest/07.eml", []string{"2018", "2018"}},
		{"casesTest/08.eml", []string{"2018", "2019"}},
		{"casesTest/09.eml", []string{"2019", "2019"}},
		{"casesTest/10.eml", []string{"2018", "2019"}},
		{"casesTest/11.eml", []string{"2018", "2018"}},
		{"casesTest/12.eml", []string{"2019", "2019"}},
		{"casesTest/13.eml", []string{"2018", "2019"}},
	}

	for _, x := range tables {
		year, _ := readEmail(x.file)

		if (year == nil) != (x.ans == nil) {
			log.Fatalf("error 1 %s", x.file)
		}

		if len(year) != len(x.ans) {
			log.Fatalf("error 2 %s", x.file)
		}

		for i := range year {
			if year[i] != x.ans[i] {
				log.Fatalf("error 3 %s", x.file)
			}
		}

		fmt.Printf("%s: pass\n", x.file)
	}

	// 
}
