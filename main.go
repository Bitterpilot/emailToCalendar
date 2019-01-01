package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"regexp"
	"strings"

	"github.com/jhillyerd/enmime"
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
	for key, val := range table {
		modlo := key % 8
		if modlo == 0 {
			fmt.Println(key, val)
			table2 = table[key:(key + 9)]
		}
	}
	fmt.Println(table2)
	return table2
}

func main() {
	cases := []string{"casesTest/01.eml", "casesTest/02.eml", "casesTest/03.eml", "casesTest/04.eml", "casesTest/05.eml", "casesTest/06.eml", "casesTest/07.eml", "casesTest/08.eml", "casesTest/09.eml", "casesTest/10.eml", "casesTest/11.eml", "casesTest/12.eml", "casesTest/13.eml"}
	_, table := readEmail(cases[0])
	// for _, v := range table {
	fmt.Println(table)
	// }

	// for k, v := range cases {
	// year, table := readEmail(v)
	// 	fmt.Printf("---\n%d\n%s\n%s\n", k, year, table)
	// }
}
