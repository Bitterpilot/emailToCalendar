package main

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/viper"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
	"github.com/bitterpilot/emailToCalendar/pkg/email/external"
)

func main() {
	// load config
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	user := viper.GetString(`user`)
	query := fmt.Sprintf("from:%s subject:%s",
		viper.GetString(`gmailFilter.sender`),
		viper.GetString(`gmailFilter.subject`))
	label := viper.GetString(`gmailFilter.label`)

	gmail := external.NewGmailSrv(user)
	listEmails := email.External.ListEmails(gmail, user, query, label)
	// for _, eml := range listEmails {
	// 	fmt.Println(eml.MsgID, eml.ThdID)
	// }

	msg := listEmails[len(listEmails)-1]
	// msg := listEmails[2]
	body := email.External.GetEmail(gmail, user, msg).Body
	fmt.Printf("%s\n", body)
	fmt.Println()
	d, err := base64.URLEncoding.DecodeString(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", d)
}
