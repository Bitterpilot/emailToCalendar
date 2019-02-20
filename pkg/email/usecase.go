package email

import (
	"github.com/bitterpilot/emailToCalendar/pkg/models"
)

type Usecase interface {
	ListEmails(user string) []*models.Email
	GetEmail(user string, msg *models.Email) *models.Email
}

type External interface {
	ListEmails(user string) []*models.Email
	GetEmail(user string, msg *models.Email) *models.Email
}
