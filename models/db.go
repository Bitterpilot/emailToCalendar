package models

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	// Database driver
	_ "github.com/mattn/go-sqlite3"
)

// DB holds a db connection
type DB struct {
	*sql.DB
}

// Tx holds a transaction
type Tx struct {
	*sql.Tx
}

// NewDB creates a new connection to a db and stores it in the DB struct
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// NewTx starts a new transaction and stores it in the Tx struct
func (db *DB) NewTx() (*Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func checkTableExists(tx *Tx, tableName string) bool {
	var exists bool
	query := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE TYPE='table' AND name='%s';", tableName)
	err := tx.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if table exists '%s' %v", tableName, err)
	}
	return exists
}

func initiateTable(tx *Tx, table interface{}) {
	var rows []string
	tableName := reflect.TypeOf(table).Name()

	val := reflect.ValueOf(table)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag

		// fieldType := field.Type    // get struct variable type
		fieldName := field.Name    //get struct variable's name
		fieldSQL := tag.Get("sql") // get struct tag's name

		row := fmt.Sprint(fieldName, " ", fieldSQL)
		rows = append(rows, row)
	}

	statement := fmt.Sprintf("CREATE TABLE %s (%s);", tableName, strings.Join(rows, ","))

	_, err := tx.Exec(statement)
	if err != nil {
		fmt.Println(err)
	}
}
