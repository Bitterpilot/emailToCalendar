package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/mattn/go-sqlite3"

	"github.com/bitterpilot/emailToCalendar/models"
)

func main() {
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

}
