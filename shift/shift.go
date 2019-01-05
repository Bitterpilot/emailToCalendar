package processor

// RowContents is what is expected from a riteq schedule email
type RowContents struct {
	day        string
	date       string
	startWork  string // gcal api expects strings
	endWork    string // https://developers.google.com/calendar/create-events#metadata
	totalHours string
	breaks     string
	pay        string
	orgLevel   string
	comments   string
}

// https://developers.google.com/calendar/extended-properties
type shift struct {
	summary        string // will be the procesed orgLevel (remove everything between \ inclusive) ***REMOVED***
	location       string // Derived from orgLevel (if item before \ = AAAA then 303 if item before \ = BBBB)
	eventDateStart string // date + startWork
	eventDateEnd   string // date + endWork
	processed      bool   // true/false/nil
}

const bodyC = `<html><head></head><body><p>Your schedule for 26 Nov 2018 through to 9 Dec 2018 is shown below</p></body></html><p><html><head></head><body><table style="width:80%;" border="1" cellspacing="0"><TD bgcolor="CornflowerBlue" align="center">Day</TD><TD bgcolor="CornflowerBlue" align="center">Date</TD><TD bgcolor="CornflowerBlue" align="center">Start Work</TD><TD bgcolor="CornflowerBlue" align="center">End Work</TD><TD bgcolor="CornflowerBlue" align="center"> Total Hours </TD><TD bgcolor="CornflowerBlue" align="center"> Breaks </TD><TD bgcolor="CornflowerBlue" align="center">Pay </TD><TD bgcolor="CornflowerBlue" align="left"> Org Level </TD><TD bgcolor="CornflowerBlue" align="center"> Comments</TD><TR><TD align="center" bgcolor="White">Mon</TD><TD align="center" bgcolor="White">26 Nov</TD><TD align="center" bgcolor="White">13:45</TD><TD bgcolor="White" align="center" rowspan="1">20:00</TD><TD bgcolor="White" align="center"> 06:15 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">05:45</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TD align="center" bgcolor="White">&nbsp;</TD><TR><TD align="center"bgcolor="LightBlue">Wed</TD><TD align="center" bgcolor="LightBlue">28 Nov</TD><TD align="center" bgcolor="LightBlue">13:45</TD><TD bgcolor="LightBlue" align="center" rowspan="1">21:45</TD><TD bgcolor="LightBlue" align="center"> 08:00 </TD><TD align="center" bgcolor="LightBlue">00:30</TD><TD align="center" bgcolor="LightBlue">07:30</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD><TD align="center" bgcolor="LightBlue">&nbsp;</TD><TR><TD align="center" bgcolor="White">Fri</TD><TD align="center" bgcolor="White">30 Nov</TD><TD align="center" bgcolor="White">13:45</TD><TD bgcolor="White" align="center" rowspan="1">20:00</TD><TD bgcolor="White" align="center"> 06:15 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">05:45</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TD align="center" bgcolor="White">New Shift</TD><TR><TD align="center" bgcolor="LightBlue">Sat</TD><TD align="center" bgcolor="LightBlue">01 Dec</TD><TD align="center" bgcolor="LightBlue">08:15</TD><TD bgcolor="LightBlue" align="center" rowspan="1">13:15</TD><TD bgcolor="LightBlue" align="center"> 05:00 </TD><TD align="center" bgcolor="LightBlue">00:00</TD><TD align="center" bgcolor="LightBlue">05:00</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD><TD align="center" bgcolor="LightBlue">&nbsp;</TD><TR><TD align="center" bgcolor="White">Mon</TD><TD align="center" bgcolor="White">03 Dec</TD><TD align="center" bgcolor="White">05:30</TD><TD bgcolor="White" align="center" rowspan="1">14:00</TD><TD bgcolor="White" align="center"> 08:30 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">08:00</TD><TD align="left" bgcolor="White">AAAA\DryOperations\Snr CSO</TD><TD align="center" bgcolor="White">&nbsp;</TD><TR><TD align="center" bgcolor="LightBlue">Wed</TD><TD align="center" bgcolor="LightBlue">05 Dec</TD><TD align="center" bgcolor="LightBlue">13:45</TD><TD bgcolor="LightBlue" align="center" rowspan="1">21:45</TD><TD bgcolor="LightBlue" align="center"> 08:00 </TD><TD align="center" bgcolor="LightBlue">00:30</TD><TD align="center" bgcolor="LightBlue">07:30</TD><TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD><TD align="center" bgcolor="LightBlue">&nbsp;</TD><TR><TD align="center" bgcolor="White">Fri</TD><TD align="center" bgcolor="White">07 Dec</TD><TD align="center" bgcolor="White">07:30</TD><TD bgcolor="White" align="center" rowspan="1">14:00</TD><TD bgcolor="White" align="center"> 06:30 </TD><TD align="center" bgcolor="White">00:30</TD><TD align="center" bgcolor="White">06:00</TD><TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD><TD align="center" bgcolor="White">&nbsp;</TD></table></body></html></p><p></p><html><head></head><body><p></p><p>Please find following your schedule should you have any concerns with the outlined dates and times please contact your supervisor.</p></p></body></html>`
const bodyN = `<html><head></head><body> <p>Your schedule for 10 Dec 2018 through to 23 Dec 2018 is shown below</p></body></html><p> <html> <head></head> <body> <table style="width:80%;" border="1" cellspacing="0"> <TD bgcolor="CornflowerBlue" align="center">Day</TD> <TD bgcolor="CornflowerBlue" align="center">Date</TD> <TD bgcolor="CornflowerBlue" align="center">Start Work</TD> <TD bgcolor="CornflowerBlue" align="center">End Work</TD> <TD bgcolor="CornflowerBlue" align="center"> Total Hours </TD> <TD bgcolor="CornflowerBlue" align="center"> Breaks </TD> <TD bgcolor="CornflowerBlue" align="center">Pay </TD> <TD bgcolor="CornflowerBlue" align="left"> Org Level </TD> <TR> <TD align="center" bgcolor="White">Mon</TD> <TD align="center" bgcolor="White">10 Dec</TD> <TD align="center" bgcolor="White">13:45</TD> <TD bgcolor="White" align="center" rowspan="1">20:00</TD> <TD bgcolor="White" align="center"> 06:15 </TD> <TD align="center" bgcolor="White">00:30</TD> <TD align="center" bgcolor="White">05:45</TD> <TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD> <TR> <TD align="center" bgcolor="LightBlue">Wed</TD> <TD align="center" bgcolor="LightBlue">12 Dec</TD> <TD align="center" bgcolor="LightBlue">09:00</TD> <TD bgcolor="LightBlue" align="center" rowspan="1">12:30</TD> <TD bgcolor="LightBlue" align="center"> 03:30 </TD> <TD align="center" bgcolor="LightBlue">00:00</TD> <TD align="center" bgcolor="LightBlue">03:30</TD> <TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD> <TR> <TD align="center" bgcolor="White">Fri</TD> <TD align="center" bgcolor="White">14 Dec</TD> <TD align="center" bgcolor="White">13:45</TD> <TD bgcolor="White" align="center" rowspan="1">20:00</TD> <TD bgcolor="White" align="center"> 06:15 </TD> <TD align="center" bgcolor="White">00:30</TD> <TD align="center" bgcolor="White">05:45</TD> <TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD> <TR> <TD align="center" bgcolor="LightBlue">Sat</TD> <TD align="center" bgcolor="LightBlue">15 Dec</TD> <TD align="center" bgcolor="LightBlue">12:00</TD> <TD bgcolor="LightBlue" align="center" rowspan="1">18:15</TD> <TD bgcolor="LightBlue" align="center"> 06:15 </TD> <TD align="center" bgcolor="LightBlue">00:00</TD> <TD align="center" bgcolor="LightBlue">06:15</TD> <TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Dry Ops Officer</TD> <TR> <TD align="center" bgcolor="White">Sun</TD> <TD align="center" bgcolor="White">16 Dec</TD> <TD align="center" bgcolor="White">13:00</TD> <TD bgcolor="White" align="center" rowspan="1">18:15</TD> <TD bgcolor="White" align="center"> 05:15 </TD> <TD align="center" bgcolor="White">00:00</TD> <TD align="center" bgcolor="White">05:15</TD> <TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD> <TR> <TD align="center" bgcolor="LightBlue">Tue</TD> <TD align="center" bgcolor="LightBlue">18 Dec</TD> <TD align="center" bgcolor="LightBlue">13:45</TD> <TD bgcolor="LightBlue" align="center" rowspan="1">21:15</TD> <TD bgcolor="LightBlue" align="center"> 07:30 </TD> <TD align="center" bgcolor="LightBlue">00:30</TD> <TD align="center" bgcolor="LightBlue">07:00</TD> <TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD> <TR> <TD align="center" bgcolor="White">Thu</TD> <TD align="center" bgcolor="White">20 Dec</TD> <TD align="center" bgcolor="White">07:30</TD> <TD bgcolor="White" align="center" rowspan="1">14:00</TD> <TD bgcolor="White" align="center"> 06:30 </TD> <TD align="center" bgcolor="White">00:30</TD> <TD align="center" bgcolor="White">06:00</TD> <TD align="left" bgcolor="White">AAAA\Dry Operations\Snr CSO</TD> <TR> <TD align="center" bgcolor="LightBlue">Fri</TD> <TD align="center" bgcolor="LightBlue">21 Dec</TD> <TD align="center" bgcolor="LightBlue">07:30</TD> <TD bgcolor="LightBlue" align="center" rowspan="1">14:00</TD> <TD bgcolor="LightBlue" align="center"> 06:30 </TD> <TD align="center" bgcolor="LightBlue">00:30</TD> <TD align="center" bgcolor="LightBlue">06:00</TD> <TD align="left" bgcolor="LightBlue">AAAA\Dry Operations\Snr CSO</TD> </table> </body> </html></p><p></p><html><head></head><body> <p></p><p>Please find following your schedule should you have any concerns with the outlined dates and times please contact your supervisor.</p></p></body></html>`

// ProcessRows takes a slice of strings and applies them to the RowContents struct
func ProcessRows(content []string) []RowContents {
	nContent := []RowContents{}
	// check where the slice needs to be devided
	position := 0
	for key, val := range content {
		if val == "Comments" {
			position = key + 1
		} else if val == "Org Level" {
			position = key + 1
		}
	}
	// fmt.Println("*********")
	start := 0
	end := position
	for i := 1; i <= (len(content) / position); i++ {
		// fmt.Println(content[start:end])
		rowLoc := content[start:end]
		rowStruct := RowContents{}
		if len(rowLoc) == 9 {
			rowStruct = RowContents{rowLoc[0], rowLoc[1], rowLoc[2], rowLoc[3], rowLoc[4], rowLoc[5], rowLoc[6], rowLoc[7], rowLoc[8]}
		} else {
			rowStruct = RowContents{rowLoc[0], rowLoc[1], rowLoc[2], rowLoc[3], rowLoc[4], rowLoc[5], rowLoc[6], rowLoc[7], ""}
		}
		nContent = append(nContent, rowStruct)
		start = end
		end = end + position
	}
	return nContent
}
