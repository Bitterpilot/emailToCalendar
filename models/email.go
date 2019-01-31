package models

type Email struct {
	ID    int
	MsgID string
	ThdID string
}

// EmailRepository is a data container from where the email entities data comes
// from. This allows the low level database concept to be ignored until we are
// read to implement it.
type EmailRepository interface {
	Find(id int) (*Email, error)
	Store(u *Email) error
}
