package email

import (
	"google.golang.org/api/gmail/v1"
)

// This won't work, it is gmail Specific
func (h *ServiceHandler) GetEmail(MsgID string) *gmail.Message {
	getMsg, err := h.Srv./*^^^Users*/.Messages.Get(h.User, MsgID).Do()
	if err != nil {
		h.Logger.Printf("Could not retrieve email: %v", err)
	}
	return getMsg
}
