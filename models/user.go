package models

import (
	"errors"
	"log"
	"reflect"
)

// The User struct declares basic user information and the SQL requierments
type User struct {
	ID   int    `sql:"INTEGER PRIMARY KEY"`
	Name string `sql:"TEXT UNIQUE"`
}

func init() {
	// create connection
	db, err := NewDB("./foo.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	tx, err := db.NewTx()
	if err != nil {
		log.Panic(err)
	}

	// tables to check
	var user User
	tableName := reflect.TypeOf(user).Name()

	// Check
	if checkTableExists(tx, tableName) == false {
		initiateTable(tx, user)
		tx.Commit()
	}
	tx.Commit()

}

// CreateUser creates a new user.
// Returns an error if user is invalid or the tx fails.
func (tx *Tx) CreateUser(u *User) error {
	// Validate the input.
	if u == nil { // struct is empty
		return errors.New("user required")
	} else if u.Name == "" {
		return errors.New("name required")
	}

	// Perform the actual insert and return any errors.
	_, err := tx.Exec(`INSERT INTO user (name) VALUES(?)`, u.Name)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: User.Name" {
			return errors.New("User Name already taken")
		}
	}
	return err
}
