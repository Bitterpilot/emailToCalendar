package main

import (
	"emailToCal/decode"
	"emailToCal/readTable"
	"fmt"
)

func main() {
	// call decode pass in a msg and recieve a body
	msg := "casesTest/01.eml"

	body, err := decode.Decode(msg)
	if err != nil {
		fmt.Println(err)
	}

	table := readTable.ReadTable(body)

	for i, k := range table {
		fmt.Println(i, k)
	}

}
