package main

import (
	"fmt"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	// test read email outputs the correct years in []string
	tables := []struct {
		input string
		want  []string
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
		got, _ := readEmail(x.input)

		if (got == nil) != (x.want == nil) {
			log.Fatalf("Got:%s Want:%s", x.input, x.want)
		}

		if len(got) != len(x.want) {
			log.Fatalf("error 2 %s", x.input)
		}

		for i := range got {
			if got[i] != x.want[i] {
				log.Fatalf("error 3 %s", x.input)
			}
		}

		fmt.Printf("%s: pass\n", x.input)
	}

	//
}
