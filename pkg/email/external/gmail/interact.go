package gmail

import (
	"errors"
	"fmt"
	"log"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
	"google.golang.org/api/gmail/v1"
)

func (gmail *gmailSrv) ListEmails(query, label string) []*email.Msg {
	listMsg, err := gmail.User.Messages.List(gmail.Username).LabelIds(label).Q(query).Do()
	if err != nil {
		log.Printf("Could not list emails: %v", err)
	}

	// convert gmail.message struct to email.Msg
	var emails []*email.Msg
	for _, msg := range listMsg.Messages {
		ms := &email.Msg{
			ExternalID: msg.Id,
			ThreadID:   msg.ThreadId,
		}
		emails = append(emails, ms)
	}

	return emails
}

func getbody(getMsg *gmail.Message) (string, error) {
	// Check Mime types
	// start by setting a high part number so if a part of the desired mime
	// type is not found we can fail gracefully.
	// parts is an array so a 0 default part num would result in selecting the
	// first element in the array a multi-part email is unlikely to have
	// more than 9000 parts ... right?
	var partNum = 9999
	for i := 0; i < len(getMsg.Payload.Parts); i++ {
		if getMsg.Payload.Parts[i].MimeType == "text/html" {
			partNum = i
		}
	}
	if partNum == 9999 {
		log.Fatalf(
			`Error: gmail.go/GetMessage text/html Mime Type not found in msgID: %s`,
			getMsg.Id)
	}

	data := getMsg.Payload.Parts[partNum].Body.Data
	// Confirm Encoding type is base64 before proceding
	partHeaders := getMsg.Payload.Parts[partNum].Headers
	for key := range partHeaders {
		if partHeaders[key].Name == "Content-Transfer-Encoding" {
			switch {
			case partHeaders[key].Value == "base64":
				return data, nil
			case partHeaders[key].Value == "quoted-printable":
				return data, nil
			default:
				errMsg := fmt.Sprintf("can not identify Content-Transfer-Encoding got: %v",
					partHeaders[key].Value)
				return "", errors.New(errMsg)
			}
		}
	}
	return "", errors.New("")
}

func (gmail *gmailSrv) GetEmail(msg *email.Msg) *email.Msg {
	getMsg, err := gmail.User.Messages.Get(gmail.Username, msg.ExternalID).Do()
	if err != nil {
		log.Printf("Could not retrieve email: %v", err)
	}

	// body
	msg.Body, err = getbody(getMsg)
	// convert gmail.message struct to email.Msg
	msg.ReceivedTime = getMsg.InternalDate
	return msg
}
