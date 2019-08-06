package models

// Email is a basic type of an email
type Email struct {
	ID           int    // Database ID.
	MsgID        string // Unique for a message in the origin service.
	ThdID        string // ThreadID(gmail) or conversationID(outlook) are related messages.
	TimeReceived int64  // Unix epoch time.
	Body         []byte // Message Body.
}
