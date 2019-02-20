package email

type Usecase interface {
	GetTable(msg *Msg) []string
}

type External interface {
	ListEmails(user, query, label string) []*Msg
	GetEmail(user string, msg *Msg) *Msg
}
