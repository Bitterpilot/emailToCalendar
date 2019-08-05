// This code was designed to be ran in main
package gmail_test

/*

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	gmailgetter "github.com/bitterpilot/emailToCalendar/app/infrastructure/gmail"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load Configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("../../config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("Config file not found: %v", err)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Config file was found with error: %v", err)
		}
	}

	user := viper.GetString("user")
	label := viper.GetString("label")
	sender := viper.GetString("sender")
	subject := viper.GetString("subject")
	// calid := viper.GetString("calendarID")
	// Logger
	log.SetReportCaller(true)

	getter := gmailgetter.NewEmailProvider(user)

	// list emails
	list, err := getter.ListMessages(label, sender, subject)
	if err != nil {
		log.Errorf("\n%v\n", err)
	}
	for _, email := range list {
		log.Infof("\nmsgID: %s\n", email.MsgID)
	}

	msg, err := getter.GetMessage(list[1])
	if err != nil {
		log.Errorf("\n%v\n", err)
	}
	log.Infof("\n%s\n", msg.Body)
}


*/