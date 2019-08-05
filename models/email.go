package models

// Email is a basic type of an email
type Email struct {
	DBID         int
	MsgID        string
	ThdID        string // ThreadID(gmail) or conversationID(outlook) are related messages
	TimeReceived int64
	Body         []byte
}
