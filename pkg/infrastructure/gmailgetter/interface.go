package gmailgetter

import (
	"google.golang.org/api/gmail/v1"
)

// EmailProvider holds objects and functions for gmail API calls
type EmailProvider struct {
	Srv  *gmail.Service
	User string
}

// NewEmailProvider creates a new Provider
func NewEmailProvider(user string) *EmailProvider {
	s := newService()
	return &EmailProvider{
		Srv:  s,
		User: user,
	}
}
