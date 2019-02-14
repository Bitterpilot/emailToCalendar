package usecase

import (
	"github.com/bitterpilot/emailToCalendar/email"
	"github.com/bitterpilot/emailToCalendar/models"
)

type emailUsecase struct {
	emailRepository email.Repository
}

// NewEmailUsecase ...
func NewEmailUsecase(e email.Repository) email.Usecase {
	return &emailUsecase{
		emailRepository: e,
	}
}

func (e *emailUsecase) ListEmails(user string) []*models.Email {
	item := e.emailRepository.ListEmails(user)
	return item
}
