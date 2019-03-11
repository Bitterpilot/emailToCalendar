package gmail

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
	"google.golang.org/api/gmail/v1"
)

type GmailHandlers struct {
	gml  *gmail.Service
	user *email.User

	logger *log.Logger
	db     *sql.DB
}

func NewGmailHandler(user *email.User, nlog *log.Logger, db *sql.DB) *GmailHandlers {
	srv := NewGmailSrv(user.Name, nlog)

	return &GmailHandlers{gml: srv, user: user, logger: nlog, db: db}
}

func (h *GmailHandlers) ListEmails() []*email.Msg {
	listMsg, err := h.gml.Users.Messages.List(h.user.Name).LabelIds(h.user.LabelID).Q(h.user.Query).Do()
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

func (h *GmailHandlers) getbody(getMsg *gmail.Message) (string, error) {
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
		h.logger.Fatalf(
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

func (h *GmailHandlers) GetEmail(msg *email.Msg) *email.Msg {
	getMsg, err := h.gml.Users.Messages.Get(h.user.Name, msg.ExternalID).Do()
	if err != nil {
		h.logger.Printf("Could not retrieve email: %v", err)
	}

	// body
	msg.Body, err = h.getbody(getMsg)
	if err != nil {
		h.logger.Println(err)
	}
	// convert gmail.message struct to email.Msg
	msg.ReceivedTime = getMsg.InternalDate
	return msg
}
