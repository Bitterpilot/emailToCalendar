package email

// Msg defines a internal type for emails
type Msg struct {
	ID           int
	ExternalID   string
	ThreadID     string
	ReceivedTime int64
	Body         string
}