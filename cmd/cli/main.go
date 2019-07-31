package main

import (
	"database/sql"
	"fmt"

	"github.com/bitterpilot/emailToCalendar/app"
	"github.com/bitterpilot/emailToCalendar/app/infrastructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../config")
	user := "me"
	// load user info
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("Config file not found: %v", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Config file was found with error: %v", err)
		}
	}

	label := viper.GetString("label")
	sender := viper.GetString("sender")
	subject := viper.GetString("subject")

	// DB
	db, err := sql.Open("sqlite3", "../../db/foo.db")
	if err != nil {
		log.Fatalf("DB Open: %v\n", err)
	}
	// insert the sqlite db connection into infrastructure and get an object with the methods that we need.
	emailStore := infrastructure.NewDB(db)

	// Email Provider
	// This just takes a user becasue I needed to hand it though somewhere.
	// It could return without anything being parsed into it to get an object with the methods that we need.
	emailGetter := infrastructure.NewGmailProvider(user)

	// Parse the objects into the main logic flow
	reg := app.NewEmailRegistar(emailGetter, emailStore)
	unproccesed, err := reg.Unprocessed(label, sender, subject)
	if err != nil {
		log.Warnf("Unprocessed: %v", err)
	}

	for _, item := range unproccesed {
		fmt.Printf("%+v\n", item)
	}
}
