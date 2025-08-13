package utils

import (
	"log"
)

func MustNotError(err error, msg string) {
	if err != nil {
		log.Fatalf("Message: %s\nError: %s\n", msg, err.Error())
	}
}

func IsValid(ok bool, msg string) bool {
	if !ok {
		log.Println(msg)
		return false
	}
	return true
}