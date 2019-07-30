package app

import (
	"log"

	"github.com/bitterpilot/emailToCalendar/models"
)

// Application logic related to retrieving and processing emails.

// EmailGetter are functions a email provider( Gmail, Outlook ) package will implement.
type EmailGetter interface {
	List(labelIDs, sender, subject string) []models.Email
}

// EmailStore are functions a database package will implement.
type EmailStore interface {
	ListByMsgID(msgID string) (string, error)
	// InsertEmail(labelIDs, sender, subject string)
}

// EmailRegistar contains dependencies for this package. Such as an email provider, database and logger.
type EmailRegistar struct {
	// email getter dependencies
	emailGetter EmailGetter
	emailStore  EmailStore
}

// NewEmailRegistar loads the dependencies declared in main() into this package for use.
func NewEmailRegistar(e EmailGetter, db EmailStore) *EmailRegistar {
	return &EmailRegistar{
		emailGetter: e,
		emailStore:  db,
	}
}

// Application logic

// Unprocessed checks
func (eR EmailRegistar) Unprocessed(labelIDs, sender, subject string) error {
	list := eR.emailGetter.List(labelIDs, sender, subject)

	for _, email := range list {
		dbRecord, err := eR.emailStore.ListByMsgID(email.MsgID)
		if err != nil {
			return err
		}
		if email.MsgID != dbRecord {
			log.Println(email.MsgID)
			// eR.emailStore.InsertEmail(email.MsgID, email.ThdID, email.TimeReceived.String())
		}
	}
	return nil
}

// id
// msg
// body
// table
