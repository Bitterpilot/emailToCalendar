package decode

import (
	"os"

	"github.com/jhillyerd/enmime"
)

// https://godoc.org/github.com/jhillyerd/enmime
// https://stackoverflow.com/questions/35038864/golang-global-variable-access

// Decode recieves a raw msg and returns the body text
func Decode(msg string) (body string, err error) {
	// Open a sample message file.
	r, err := os.Open(msg)
	if err != nil {
		return "", err
	}

	// Parse message body with enmime.
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return "", err
	}

	// fmt.Println(reflect.TypeOf(env.HTML)) = string

	return env.HTML, nil
}
