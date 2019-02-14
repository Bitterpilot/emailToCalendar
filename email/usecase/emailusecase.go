package usecase

import (
	"github.com/bitterpilot/emailToCalendar/email"
	"github.com/bitterpilot/emailToCalendar/models"
)

type emailUsecase struct {
	emailexternal email.external
}

// NewEmailUsecase ...
func NewEmailUsecase(e email.external) email.Usecase {
	return &emailUsecase{
		emailexternal: e,
	}
}

func (e *emailUsecase) ListEmails(user string) []*models.Email {
	item := e.emailexternal.ListEmails(user)
	return item
}
