package email

import (
	"log"
	"testing"
)

func TestNewService(t *testing.T) {
	expect := *ServiceHandler{}
	got, err := NewService("me", "gmail", *log.Logger)
	if err != nil {
		t.Fatalf("Returned Error\n", err)
	}
	if got != expect {
		t.FailNow()
	}
}
