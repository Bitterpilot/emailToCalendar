package app

import (
	log "github.com/sirupsen/logrus"

	"github.com/bitterpilot/emailToCalendar/models"
)

// Application logic related to retrieving and processing emails.

// EmailGetter are functions a email provider( Gmail, Outlook ) package will implement.
type EmailGetter interface {
	List(labelIDs, sender, subject string) []models.Email
	Get(models.Email) models.Email
}

// EmailStore are functions a database package will implement.
type EmailStore interface {
	FindByMsgID(msgID string) (string, error)
	InsertEmail(models.Email) (int, error)
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

// Unprocessed checks all messages in a users inbox against Database records looking for emails marked as unprocessed.
// It returns a models.Email with all fields.
func (eR EmailRegistar) Unprocessed(labelIDs, sender, subject string) ([]models.Email, error) {
	// List all emails
	list := eR.emailGetter.List(labelIDs, sender, subject)
	var unprocessed []models.Email
	// compare emails to db records
	for _, email := range list {
		// Look for email in DB.
		dbRecord, err := eR.emailStore.FindByMsgID(email.MsgID)
		if err != nil {
			return nil, err
		}
		// If not in db then get time received and body, insert to db,
		if email.MsgID != dbRecord {
			log.WithFields(log.Fields{"EmailID": email.MsgID}).Info("Get message and insert into db.")
			// Get body and time received
			email = eR.emailGetter.Get(email)
			email.DBID, err = eR.emailStore.InsertEmail(email)
			if err != nil {
				return nil, err
			}

			unprocessed = append(unprocessed, email)
		}
	}
	return unprocessed, nil
}

// body
// table
