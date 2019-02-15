package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/bitterpilot/emailToCalendar/email"
	"github.com/bitterpilot/emailToCalendar/email/external"
)

func main() {
	// load config
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	user := viper.GetString(`user`)

	usecase := external.NewGmailSrv(user)
	listEmails := email.Usecase.ListEmails(usecase, user)
	// for _, eml := range listEmails {
	// 	fmt.Println(eml.MsgID, eml.ThdID)
	// }

	msg := listEmails[len(listEmails)-1]
	// msg := listEmails[2]
	body := email.External.GetEmail(usecase, user, msg).Body
	fmt.Printf("%s\n", body)
}
