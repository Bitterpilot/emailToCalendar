package app

import (
	log "github.com/sirupsen/logrus"

	"github.com/bitterpilot/emailToCalendar/models"
)

// Application logic related to retrieving and processing emails.

// EmailGetter are functions a email provider( Gmail, Outlook ) package will implement.
type EmailGetter interface {
	ListMessages(labelIDs, sender, subject string) ([]models.Email, error)
	GetMessage(models.Email) (models.Email, error)
}

// EmailStore are functions a database package will implement.
type EmailStore interface {
	FindByMsgID(string) (string, error)
	CheckUnProcessed(e models.Email) (models.Email, error)
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
	list, err := eR.emailGetter.ListMessages(labelIDs, sender, subject)
	if err != nil {
		return nil, err
	}
	var unprocessed []models.Email
	// compare emails to db records
	for _, email := range list {
		// Set sensible default for dbid
		email.ID = -1
		// Look for email in DB.
		email, err := eR.emailStore.CheckUnProcessed(email)
		if err != nil {
			return nil, err
		}
		// If not in db then get time received and body
		if !email.Processed {
			log.WithFields(log.Fields{"EmailID": email.MsgID, "DatabaseID": email.ID}).Info("Got message.")
			// Get body and time received
			email, err = eR.emailGetter.GetMessage(email)
			if err != nil {
				return nil, err
			}
			// If not in db insert to db
			if email.ID == -1 {
				email.ID, err = eR.emailStore.InsertEmail(email)
				if err != nil {
					return nil, err
				}
			}
			unprocessed = append(unprocessed, email)
		}
	}
	return unprocessed, nil
}
