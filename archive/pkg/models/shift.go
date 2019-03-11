package models

// Shift is the concept of a shift
type Shift struct {
	Summary       string
	location      string
	Start         string
	End           string
	MsgID         string
	Processed     bool
	ProcessedTime int
	Removed       bool
	RemovedTime   int
}

// ShiftRepository is a data container from where the shift entities data comes
// from. This allows the low level database concept to be ignored until we are
// read to implement it.
type ShiftRepository interface {
	Find(id int) (*Shift, error)
	Store(u *Shift) error
}
