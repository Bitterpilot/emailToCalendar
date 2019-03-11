package gmail

import (
	"fmt"

	"google.golang.org/api/gmail/v1"
)

func (h *GmailHandlers) Watch(labelIds []string) {
	if len(labelIds) == 0 {
		labelIds[0] = "inbox"
	}

	req := &gmail.WatchRequest{
		LabelFilterAction: "include",
		LabelIds:          labelIds,
		TopicName:         "projects/emailtocalandar-1546160165278/topics/watch",
		ForceSendFields:   nil,
		NullFields:        nil,
	}
	rsp, err := h.gml.Users.Watch(h.user.Name, req).Do()
	if err != nil {
		fmt.Println(err)
	}
	body, err := rsp.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", body)
}
