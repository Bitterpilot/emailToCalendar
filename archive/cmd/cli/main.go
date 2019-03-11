package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
	"github.com/bitterpilot/emailToCalendar/pkg/email/interfaces/gmail"

	"github.com/spf13/viper"

	"github.com/bitterpilot/emailToCalendar/pkg/email/db"
	"github.com/bitterpilot/emailToCalendar/pkg/email/interfaces/store"
)

func main() {
	// load config
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	user := &email.User{
		Name: viper.GetString(`user`),
		Query: fmt.Sprintf("from:%s subject:%s",
			viper.GetString(`gmailFilter.sender`),
			viper.GetString(`gmailFilter.subject`)),
		LabelID: viper.GetString(`gmailFilter.label`),
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	g := gmail.NewGmailHandler(user, logger, nil)

	listEmails := g.ListEmails()

	msg := listEmails[len(listEmails)-1]
	msg.Body = g.GetEmail(msg).Body

	x, err := msg.ReadDate()
	if err != nil {
		fmt.Println("%v", err)
	}
	fmt.Println(x)

	dbHandler := db.NewSqliteHandler(viper.GetString(`db`))
	handlers := make(map[string]store.DbHandler)
	handlers["DbEmailRepo"] = dbHandler

	newE := store.NewDbEmailRepo(handlers)

	nm := newE.FindByExternalThreadID("167b9fa15f6230de")
	for _, val := range nm {
		fmt.Printf("Id: %d, Ext:%s, Thd: %s\n", val.ID, val.ExternalID, val.ThreadID)
	}
}
