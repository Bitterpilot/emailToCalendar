package infrastructure

import (
	gmailgetter "github.com/bitterpilot/emailToCalendar/app/infrastructure/gmail"
	"github.com/bitterpilot/emailToCalendar/models"
	"google.golang.org/api/gmail/v1"
)

// Provider
type Provider struct {
	service *gmail.Service
	// user    string
}

var user string

// NewGmailProvider
func NewGmailProvider(u string) *Provider {
	s := gmailgetter.NewService()
	user = u
	return &Provider{
		service: s,
		// user:    u,
	}
}

// List
func (p Provider) List(labelIDs, sender, subject string) []models.Email {
	list := gmailgetter.ListMessages(user, labelIDs, sender, subject)
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

func (p Provider) Get(e models.Email) models.Email {
	InternalDate, emailBody := gmailgetter.GetMessage(user, e.MsgID)
	e.TimeReceived = InternalDate
	e.Body = emailBody

	return e
}
