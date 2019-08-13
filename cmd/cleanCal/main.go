package main

import (
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/bitterpilot/emailToCalendar/models"
	"github.com/bitterpilot/emailToCalendar/pkg/infrastructure/calendargetter"
)

func main() {
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

	// Lock to the test calendar
	cg := calendargetter.NewCalendarProvider("d9sun7p76ae8keqjcjbnf3mb6k@group.calendar.google.com")

	// Clean up
	list, err := cg.List()
	if err != nil {
		log.Error(err)
	}
	total := len(list)
	var count int
	for _, e := range list {
		err := cg.Delete(e)
		if err != nil {
			log.Errorf("%s\t%+v\n", pretty.Sprint(e), err)
		} else {
			count++
		}
	}
	log.Infof("Deleted %d of %d\n", count, total)
}
