package readTable

import (
	"strings"

	"golang.org/x/net/html"
)

// var body = strings.NewReader(`<html><head></head><body><p>Your schedule for 10 Dec 2018 through to 23 Dec 2018 is shown below</p></body></html><p><html><head></head><body><table style="width:80%;" border="1" cellspacing="0"><TD bgcolor="CornflowerBlue" align="center">Day</TD><TD bgcolor="CornflowerBlue" align="center">Date</TD><TD bgcolor="CornflowerBlue" align="center">Start Work</TD><TD bgcolor="CornflowerBlue" align="center">End Work</TD><TD bgcolor="CornflowerBlue" align="center"> Total Hours </TD><TD bgcolor="CornflowerBlue" align="center"> Breaks </TD><TD bgcolor="CornflowerBlue" align="center">Pay </TD><TD bgcolor="CornflowerBlue" align="left"> Org Level </TD><TR><TD align="center" bgcolor="White">Mon</TD><TD align="center" bgcolor="White">10 Dec</TD><TD align="center" bgcolor="White">13:45</TD><TD bgcolor="White" align="center" rowspan="1">20:00</TD><TD bgcolor="White" align="center"> 06:15 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">05:45</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TR><TD align="center" bgcolor="LightBlue">Wed</TD><TD align="center" bgcolor="LightBlue">12 Dec</TD><TD align="center" bgcolor="LightBlue">09:00</TD><TD bgcolor="LightBlue" align="center" rowspan="1">12:30</TD><TD bgcolor="LightBlue" align="center"> 03:30 </TD><TD align="center" bgcolor="LightBlue">00:00</TD><TD align="center" bgcolor="LightBlue">03:30</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD><TR><TD align="center" bgcolor="White">Fri</TD><TD align="center" bgcolor="White">14 Dec</TD><TD align="center" bgcolor="White">13:45</TD><TD bgcolor="White" align="center" rowspan="1">20:00</TD><TD bgcolor="White" align="center"> 06:15 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">05:45</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TR><TD align="center" bgcolor="LightBlue">Sat</TD><TD align="center" bgcolor="LightBlue">15 Dec</TD><TD align="center" bgcolor="LightBlue">12:00</TD><TD bgcolor="LightBlue" align="center" rowspan="1">18:15</TD><TD bgcolor="LightBlue" align="center"> 06:15 </TD><TD align="center" bgcolor="LightBlue">00:00</TD><TD align="center" bgcolor="LightBlue">06:15</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Dry Ops Officer</TD><TR><TD align="center" bgcolor="White">Sun</TD><TD align="center" bgcolor="White">16 Dec</TD><TD align="center" bgcolor="White">13:00</TD><TD bgcolor="White" align="center" rowspan="1">18:15</TD><TD bgcolor="White" align="center"> 05:15 </TD><TD align="center" bgcolor="White">00:00</TD><TD align="center" bgcolor="White">05:15</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TR><TD align="center" bgcolor="LightBlue">Tue</TD><TD align="center" bgcolor="LightBlue">18 Dec</TD><TD align="center" bgcolor="LightBlue">13:45</TD><TD bgcolor="LightBlue" align="center" rowspan="1">21:15</TD><TD bgcolor="LightBlue" align="center"> 07:30 </TD><TD align="center" bgcolor="LightBlue">00:30</TD><TD align="center" bgcolor="LightBlue">07:00</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD><TR><TD align="center" bgcolor="White">Thu</TD><TD align="center" bgcolor="White">20 Dec</TD><TD align="center" bgcolor="White">07:30</TD><TD bgcolor="White" align="center" rowspan="1">14:00</TD><TD bgcolor="White" align="center"> 06:30 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">06:00</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TR><TD align="center" bgcolor="LightBlue">Fri</TD><TD align="center" bgcolor="LightBlue">21 Dec</TD><TD align="center" bgcolor="LightBlue">07:30</TD><TD bgcolor="LightBlue" align="center" rowspan="1">14:00</TD><TD bgcolor="LightBlue" align="center"> 06:30 </TD><TD align="center" bgcolor="LightBlue">00:30</TD><TD align="center" bgcolor="LightBlue">06:00</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD></table></body></html></p><p></p><html><head></head><body><p></p><p>Please find following your schedule should you have any concerns with the outlined dates and times please contact your supervisor.</p></p></body></html>`)

// ReadTable takes a string containgin a html table
func ReadTable(body string) (table []string) {
	b := strings.NewReader(body)

	z := html.NewTokenizer(b)
	content := []string{}

	// While have not hit the </table> tag
	for z.Token().Data != "table" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == "td" {
				inner := z.Next()
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := strings.TrimSpace(text)
					content = append(content, t)
				}
			}
		}
	}
	// Print to check the slice's content
	// for i, k := range content {
	// 	fmt.Println(i, k)
	// }
	return content
}
