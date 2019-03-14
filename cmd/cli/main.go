package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitterpilot/emailToCalendar/pkg/email"
)

func main() {
	user := "me"
	provider := "gmail"
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	s, err := email.NewService(user, provider, logger)
	if err != nil {
		fmt.Println(err)
	}
	m := s.GetEmail("167c40b0acd40d44")
	fmt.Println(m.Snippet)
}
