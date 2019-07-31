package gmailgetter

import (
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

// GetMessage ...
func GetMessage(user, msgID string) (int64, []byte) {
	msg, err := srv.Users.Messages.Get(user, msgID).Format("full").Do()
	if err != nil {
		// FIXME: Make sure a fatal is appropriate
		// 		  fatal will exit to OS
		// guide https://stackoverflow.com/a/33890104
		log.Fatalf("Error: gmail.go/GetMessage/msg returned %v", err)
	}

	// Check Mime types
	// start by setting a high part number so if a part of the desired mime type is not found
	// we can fail gracefully.
	// parts is an array so a 0 default part num would result in selecting the first element in the array
	// a multi-part email is unlikely to have more than 9000 parts (a foolish assumption?)
	var partNum = 9999
	for i := 0; i < len(msg.Payload.Parts); i++ {
		if msg.Payload.Parts[i].MimeType == "text/html" {
			partNum = i
		}
	}
	if partNum == 9999 {
		// FIXME: Make sure a fatal is appropriate
		// 		  fatal will exit to OS
		// guide https://stackoverflow.com/a/33890104
		log.Fatalf("Error: gmail.go/GetMessage text/html Mime Type not found in msgID: %s", msgID)
	}
	// Confirm Encoding type is base64 before proceding
	partHeaders := msg.Payload.Parts[partNum].Headers
	for key := range partHeaders {
		if partHeaders[key].Name == "Content-Transfer-Encoding" {
			if partHeaders[key].Value != "base64" {
				// FIXME: Make sure a fatal is appropriate
				// 		  fatal will exit to OS
				// guide https://stackoverflow.com/a/33890104
				log.Fatalf("Unexpected Content-Transfer-Encoding type: %v in in msgID: %s", partHeaders[key].Value, msgID)
			}
		}
	}

	// Decode and print
	data := msg.Payload.Parts[partNum].Body.Data
	emailBody, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		// FIXME: Make sure a fatal is appropriate
		// 		  fatal will exit to OS
		// guide https://stackoverflow.com/a/33890104
		log.Fatalf("Unable to decode email: %v", err)
	}

	return msg.InternalDate, emailBody
}

// ListMessages ...
func ListMessages(user, labelIDs, sender, subject string) []*gmail.Message {
	msgList, err := srv.Users.Messages.List(user).LabelIds(labelIDs).Q(fmt.Sprintf("from:%s subject:%s", sender, subject)).Do()
	if err != nil {
		fmt.Println(err)
	}
	return msgList.Messages
}
