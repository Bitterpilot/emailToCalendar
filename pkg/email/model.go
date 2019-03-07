package email

// Msg defines a internal type for emails
type Msg struct {
	ID           int
	ExternalID   string
	ThreadID     string
	ReceivedTime int64
	Body         string
}

type External interface {
	ListEmails(query, label string) []*Msg
	GetEmail(msg *Msg) *Msg
	Watch(labelIds []string)
}

type Store interface {
	Store(msg Msg) error
	FindByID(ID int) Msg
	FindByExternalID(id int) Msg
	FindByExternalThreadID(id int) Msg
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}
