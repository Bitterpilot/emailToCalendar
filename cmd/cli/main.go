package main

import (
	"fmt"

	"github.com/bitterpilot/emailToCalendar/email"
	"github.com/bitterpilot/emailToCalendar/email/external"
)

func main() {
	user := "me"
	euse := external.NewGmailSrv(user)
	ListEmails := email.Usecase.ListEmails(euse, user)
	for _, eml := range ListEmails {
		fmt.Println(eml.MsgID)
	}
}
