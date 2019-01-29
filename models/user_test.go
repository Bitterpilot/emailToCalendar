package models

import (
	"math/rand"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	testSet := []struct {
		name   string
		expect string
	}{
		// empty string case
		{name: "",
			expect: "Invalid username"},
		// empty string case
		{name: "0",
			expect: "Invalid username"},
		// over 250 chars case
		{name: RandString(251),
			expect: "Invalid username"},
		// exactly 250 chars case
		{name: RandString(250),
			expect: "Invalid username"},
		// over 249 chars case
		{name: RandString(249),
			expect: "nil"},
		// illegal character " "
		{name: "Susie Q.",
			expect: "Invalid username"},
		// illegal character " "
		{name: "Susie@Q.",
			expect: "Invalid username"},
		// good name
		{name: "Susie_Q.",
			expect: "nil"},
	}
	for key, value := range testSet {
		user := &User{Name: value.name}
		err := user.Validate()
		if err != nil {
			if err.Error() != value.expect {
				if value.name == "" {
					value.name = "empty sting"
				}
				t.Errorf("Item %d failed. Expected %s for %s",
					key-1, value.expect, value.name)
			}
		}
	}
}

// Random string generator
// https://stackoverflow.com/a/31832326
// LANGISSUE: this only tests latin alphabet characters
var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
