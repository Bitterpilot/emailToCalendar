package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/bitterpilot/emailToCalendar/pkg/email/db"
	"github.com/bitterpilot/emailToCalendar/pkg/email/external"	
	"github.com/bitterpilot/emailToCalendar/pkg/email/external/gmail"
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

	g := gmail.NewGmailSrv(user)
	listEmails := g.ListEmails(query, label)

	msg := listEmails[len(listEmails)-1]
	msg.Body = g.GetEmail(msg).Body

	x, _ := msg.ReadDate()
	fmt.Println(x)

	dbHandler := db.NewSqliteHandler(viper.GetString(`db`))
	handlers := make(map[string]external.DbHandler)
	handlers["DbEmailRepo"] = dbHandler

	newE := external.NewDbEmailRepo(handlers)

	nm := newE.FindByExternalThreadID("167b9fa15f6230de")
	for _, val := range nm {
		fmt.Printf("Id: %d, Ext:%s, Thd: %s\n", val.ID, val.ExternalID, val.ThreadID)
	}
}
