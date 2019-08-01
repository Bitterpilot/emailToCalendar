package infrastructure

import (
	gmailgetter "github.com/bitterpilot/emailToCalendar/app/infrastructure/gmail"
	"github.com/bitterpilot/emailToCalendar/models"
	"google.golang.org/api/gmail/v1"
)

// Provider holds objects and functions for gmail API calls
type Provider struct {
	service *gmail.Service
	user    string
}


// NewGmailProvider creates a new Provider
func NewGmailProvider(u string) *Provider {
	s := gmailgetter.NewService()
	return &Provider{
		service: s,
		user:    u,
	}
}

// List converts a gmail message type to a model.Email type
func (p Provider) List(labelIDs, sender, subject string) []models.Email {
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
func (p Provider) Get(e models.Email) models.Email {
	InternalDate, emailBody := gmailgetter.GetMessage(p.user, e.MsgID)
	e.TimeReceived = InternalDate
	e.Body = emailBody

	return e
}
