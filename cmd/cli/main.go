package main

import (
	"database/sql"

	"github.com/kr/pretty"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	app "github.com/bitterpilot/emailToCalendar/internal"
	"github.com/bitterpilot/emailToCalendar/models"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/calendargetter"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/gmailgetter"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/sqlite"
)

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	// load user info
	viper.SetConfigName("config")
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found: %v", err)
		} else {
			log.Fatalf("Config file was found with error: %v", err)
		}
	}
	// Load config
	var c *models.Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("Failed to unmarshal file into config\n")
	}

	// Setup database
	db, err := sql.Open("sqlite3", "../../db/foo.db")
	if err != nil {
		log.Fatalf("Failed to open db: %v\n", err)
	}
	// set up Providers
	e := gmailgetter.NewEmailProvider(c.User)
	cal := calendargetter.NewCalendarProvider(c.CalendarID)
	// set up app packages
	ed, cald := sqlite.NewSqliteDB(db)
	emlreg := app.NewEmailRegistar(e, ed)
	calreg := app.NewCalendarRegistar(cal, cald, c)

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
