package email

import (
	"google.golang.org/api/gmail/v1"
	"log"
	"github.com/bitterpilot/emailToCalendar/pkg/email/providers"
)

type ServiceHandler struct {
	user    string
	Service interface{}
	logger  *log.Logger
}

func NewService(user, provider string, logger *log.Logger) (*ServiceHandler, error) {
	switch {
	case provider == "gmail":
		s, err := gmailHandler(user)
		if err != nil {
			return nil, err
		}

		sh := &ServiceHandler{
			user:    user,
			Service: s,
			logger:  logger,
		}
		return s, nil
	}
}

func gmailHandler(user string) (*gmail.Service, error) {
	s := providers.NewGmailSrv(user)
    return s,nil
}
