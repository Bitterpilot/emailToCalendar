package models

import (
	"errors"
	"regexp"
)

// The User struct declares basic user information and the SQL column requirements
type User struct {
	ID   int    `sql:"INTEGER PRIMARY KEY"`
	Name string `sql:"TEXT UNIQUE"`
	// Label    string `sql:"TEXT"`                // The label the message is under
	// Sender   string `sql:"TEXT"`                // an email address
	// Subject  string `sql:"TEXT"`                // the expected subject line
	// Calendar string `sql:"TEXT"`                // the calendar to place the events
}

type UserRepository interface {
	Find(id int) (*User, error)
	Store(u *User) error
}

// reUsername is a regular expresion checking for illegal characters
// legal characters are
//		a-z	A-Z	0-9	- _ .
// LANG: this only tests latin alphabet characters
var reUsername = regexp.MustCompile("^[a-zA-Z0-9-_.]+$")

// ValidateUser ensures the username is valid
func (u *User) Validate() error {
	switch {
	case len(u.Name) == 0:
		return errors.New("Invalid username")
	case len(u.Name) > 250:
		return errors.New("Invalid username")
	case !reUsername.MatchString(u.Name):
		return errors.New("Invalid username")
	default:
		return nil
	}
}
