package usecase

import (
	"github.com/bitterpilot/emailToCalendar/pkg/email"
	"github.com/bitterpilot/emailToCalendar/pkg/models"
)

type emailUsecase struct {
	emailExternal email.External
}

// NewEmailUsecase ...
func NewEmailUsecase(e email.External) email.Usecase {
	return &emailUsecase{
		emailExternal: e,
	}
}

func (e *emailUsecase) GetEmail(user string, msg *models.Email) *models.Email {
	return nil
}

func (e *emailUsecase) ListEmails(user string) []*models.Email {
	item := e.emailExternal.ListEmails(user)
	return item
}
