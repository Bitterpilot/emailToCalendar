package email

import (
	"google.golang.org/api/gmail/v1"
	"log"

	"github.com/bitterpilot/emailToCalendar/pkg/email/providers"
)

type ServiceHandler struct {
	User   string
	Srv    *gmail.Service
	Logger *log.Logger
}

func NewService(user, provider string, logger *log.Logger) (*ServiceHandler, error) {
	switch {
	case provider == "gmail":
		s, err := gmailHandler(user)
		if err != nil {
			return nil, err
		}

		sh := &ServiceHandler{
			User:   user,
			Srv:    s,
			Logger: logger,
		}
		return sh, nil
	}
	return nil, nil
}

func gmailHandler(u string) (*gmail.Service, error) {
	s := providers.NewGmailSrv(u)
	return s, nil
}
