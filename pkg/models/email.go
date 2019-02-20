package models

type Email struct {
	ID       int
	MsgID    string
	ThdID    string
	Received int64
	Body     []byte
}

// EmailExternal is a data container from where the email entities data comes
// from. This allows the low level database concept to be ignored until we are
// read to implement it.
type EmailExternal interface {
	Find(id int) (*Email, error)
	Store(u *Email) error
}
