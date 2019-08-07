package gmailgetter

import (
	"encoding/base64"
	"fmt"

	"github.com/bitterpilot/emailToCalendar/models"
	"google.golang.org/api/gmail/v1"
)

// GetMessage returns the internal date and body from the Gmail API.
func (p *EmailProvider) GetMessage(e models.Email) (models.Email, error) {
	msg, err := p.Srv.Users.Messages.Get(p.User, e.MsgID).Format("full").Do()
	if err != nil {
		return models.Email{}, err
	}

	e.TimeReceived = msg.InternalDate
	e.Body, err = getBody(msg)
	if err != nil {
		return models.Email{}, err
	}
	return e, nil
}

// getBody from gmail Message
func getBody(msg *gmail.Message) ([]byte, error) {
	// Check Mime types
	// start by setting a high part number so if a part of the desired mime type is not found
	// we can fail gracefully.
	// parts is an array so a 0 default part num would result in selecting the first element in the array
	var partNum = -1
	for i := 0; i < len(msg.Payload.Parts); i++ {
		if msg.Payload.Parts[i].MimeType == "text/html" {
			partNum = i
		}
	}
	if partNum == -1 {
		return nil, fmt.Errorf("text/html Mime Type not found in msgID: %s", msg.Id)
	}
	// Confirm Encoding type is base64 before proceding
	partHeaders := msg.Payload.Parts[partNum].Headers
	for key := range partHeaders {
		if partHeaders[key].Name == "Content-Transfer-Encoding" {
			if partHeaders[key].Value != "base64" {
				return nil, fmt.Errorf("unexpected Content-Transfer-Encoding type: %v in in msgID: %s", partHeaders[key].Value, msg.Id)
			}
		}
	}

	// Decode and print
	data := msg.Payload.Parts[partNum].Body.Data
	return base64.URLEncoding.DecodeString(data)
}

// ListMessages returns a list of messages with basic info from the Gmail API.
// It only returns the messages with a declared label, sender and subject.
//
// Returned values are Gmail Message ID and Thread ID.
func (p *EmailProvider) ListMessages(labelIDs, sender, subject string) ([]models.Email, error) {
	msgList, err := p.Srv.Users.Messages.List(p.User).LabelIds(labelIDs).Q(fmt.Sprintf("from:%s subject:%s", sender, subject)).Do()
	if err != nil {
		return nil, err
	}

	ret := make([]models.Email, 0, len(msgList.Messages))
	for _, msg := range msgList.Messages {
		ret = append(ret, models.Email{
			MsgID: msg.Id,
			ThdID: msg.ThreadId,
		})
	}
	return ret, nil
}
