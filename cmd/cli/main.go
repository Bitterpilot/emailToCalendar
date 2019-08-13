package main

import (
	"database/sql"
	"flag"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	app "github.com/bitterpilot/emailToCalendar/internal"
	"github.com/bitterpilot/emailToCalendar/internal/cmd"
	"github.com/bitterpilot/emailToCalendar/models"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/calendargetter"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/gmailgetter"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/sqlite"
)

func main() {
	// Set up log level
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Sets the logger to debug level. Defaults to false")
	flag.Parse()
	if debug {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
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

	// Run App
	cmd.Run(c, emlreg, calreg)
}
