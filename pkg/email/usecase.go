package email

type Usecase interface {
}

type External interface {
	ListEmails(user string) []*Msg
	GetEmail(user string, msg *Msg) *Msg
}
