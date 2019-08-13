package cmd

import (
	"github.com/kr/pretty"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	app "github.com/bitterpilot/emailToCalendar/internal"
	"github.com/bitterpilot/emailToCalendar/models"
)

// Run
func Run(c *models.Config, emlreg *app.EmailRegistar, calreg *app.CalendarRegistar) {
	// Logic
	list, err := emlreg.Unprocessed(c.Label, c.Sender, c.Subject)
	if err != nil {
		log.Fatalln(err)
	}
	if list == nil {
		log.Info("no new emails")
	}

	for i, e := range list {
		years, rows, err := app.ReadBody(e)
		if err != nil {
			log.Fatalln(err)
		}
		for _, row := range rows[1:] {
			event, err := calreg.BuildEvent(years, row, e.MsgID)
			if err != nil {
				log.Errorln(err)
			}
			list[i].List = append(list[i].List, event)
		}
	}

	for i, ev := range list {
		for j, evnt := range ev.List {
			list[i].List[j], err = calreg.Publish(evnt)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	for _, e := range list {
		var count int
		for _, evnt := range e.List {
			if evnt.Processed {
				count++
			}
		}
		if count == len(e.List) {
			e.Processed = true
			err := emlreg.MarkedAsProcessed(e)
			if err != nil {
				log.Fatalln(err)
			}
			log.Infof("\n*******************\nMSG: %s\nAll Processed!\n*******************\n", e.MsgID)
		} else {
			for _, ev := range e.List {
				if !ev.Processed {
					log.Infof("%+v\n", pretty.Formatter(ev))
				}
			}
		}
	}
}
