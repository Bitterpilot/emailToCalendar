package gmail

import (
	"google.golang.org/api/gmail/v1"
)

func (srv *gmailSrv) watch(labelIds []string) {
	if len(labelIds) == 0 {
		labelIds[0] = "inbox"
	}

	req := &gmail.WatchRequest{
		LabelFilterAction: "include",
		LabelIds:          labelIds,
		TopicName:         "projects/emailtocalandar-1546160165278/topics/gmailAPI",
		ForceSendFields:   nil,
		NullFields:        nil,
	}
	srv.User.Watch(srv.Username, req)
}
