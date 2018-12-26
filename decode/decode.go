package main

import (
	"fmt"
	"github.com/jhillyerd/enmime"
	"os"
)

// https://godoc.org/github.com/jhillyerd/enmime
// https://stackoverflow.com/questions/35038864/golang-global-variable-access

func main() {
	// Open a sample message file.
	r, err := os.Open("./test.1.raw")
	if err != nil {
		fmt.Print(err)
		return
	}

	// Parse message body with enmime.
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		fmt.Print(err)
		return
	}

	// Headers can be retrieved via Envelope.GetHeader(name).
	fmt.Printf("From: %v\n", env.GetHeader("From"))

	// Address-type headers can be parsed into a list of decoded mail.Address structs.
	alist, _ := env.AddressList("To")
	for _, addr := range alist {
		fmt.Printf("To: %s <%s>\n", addr.Name, addr.Address)
	}

	// enmime can decode quoted-printable headers.
	fmt.Printf("Subject: %v\n", env.GetHeader("Subject"))

	// The plain text body is available as mime.Text.
	fmt.Printf("Text Body: %v chars\n", len(env.Text))

	// The HTML body is stored in mime.HTML.
	fmt.Printf("HTML Body: %v chars\n", len(env.HTML))

	// mime.Inlines is a slice of inlined attacments.
	fmt.Printf("Inlines: %v\n", len(env.Inlines))

	// mime.Attachments contains the non-inline attachments.
	fmt.Printf("Attachments: %v\n", len(env.Attachments))

	fmt.Println(env.HTML)

}
