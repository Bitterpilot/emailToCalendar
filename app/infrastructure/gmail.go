package infrastructure

import (
	gmailgetter "github.com/bitterpilot/emailToCalendar/app/infrastructure/gmail"
	"github.com/bitterpilot/emailToCalendar/models"
	"google.golang.org/api/gmail/v1"
)

// EmailProvider holds objects and functions for gmail API calls
type EmailProvider struct {
	service *gmail.Service
	user    string
}


// NewEmailProvider creates a new Provider
func NewEmailProvider(u string) *EmailProvider {
	s := gmailgetter.NewService()
	return &EmailProvider{
		service: s,
		user:    u,
	}
}

// List converts a gmail message type to a model.Email type
func (p EmailProvider) List(labelIDs, sender, subject string) []models.Email {
	list := gmailgetter.ListMessages(p.user, labelIDs, sender, subject)
	var ret []models.Email
	for _, gmsg := range list {
		msg := models.Email{
			MsgID: gmsg.Id,
			ThdID: gmsg.ThreadId,
		}
		ret = append(ret, msg)
	}
	return ret
}

// Get a body and date from a gmail message type and converts to a model.Email type
func (p EmailProvider) Get(e models.Email) models.Email {
	InternalDate, emailBody := gmailgetter.GetMessage(p.user, e.MsgID)
	e.TimeReceived = InternalDate
	e.Body = emailBody

	return e
}
