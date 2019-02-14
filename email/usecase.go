package email

import (
	"github.com/bitterpilot/emailToCalendar/models"
)

type Usecase interface {
	ListEmails(user string) []*models.Email
}

type External interface {
	ListEmails(user string) []*models.Email
}
