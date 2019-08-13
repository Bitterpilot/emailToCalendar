package cmd

import (
	"os"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	app "github.com/bitterpilot/emailToCalendar/internal"
	"github.com/bitterpilot/emailToCalendar/models"
)

// Run connects all the components from getting an email to publishing events.
func Run(c *models.Config, emailService *app.EmailRegistar, calendarService *app.CalendarRegistar) {
	// Get list of unprocessed emails
	emails, err := emailService.Unprocessed(c.Label, c.Sender, c.Subject)
	if err != nil {
		log.Fatalln(err)
	}
	if emails == nil {
		log.Info("no unprocessed emails")
		os.Exit(0)
	}

	// Decode table into events
	for i, email := range emails {
		// Decode into individual rows
		years, rows, err := app.ReadBody(email)
		if err != nil {
			log.Fatalln(err)
		}

		// Build events from rows. Skipping the header row.
		for _, row := range rows[1:] {
			event, err := calendarService.BuildEvent(years, row, email.MsgID)
			if err != nil {
				log.Errorln(err)
			}
			emails[i].List = append(emails[i].List, event)
		}
	}

	// Publish each event
	for i, email := range emails {
		for j, event := range email.List {
			emails[i].List[j], err = calendarService.Publish(event)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	// Check that all events were published
	for _, e := range emails {
		var count int
		for _, event := range e.List {
			if event.Processed {
				count++
			}
		}
		if count == len(e.List) {
			e.Processed = true
			err := emailService.MarkedAsProcessed(e)
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
