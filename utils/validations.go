package utils

import (
	"fmt"
	"regexp"
	"unicode"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func ValidateEmail(e string) error {
	if !rxEmail.MatchString(e) {
		err := fmt.Errorf("Email address %s is invalid", e)
		return err
	}
	return nil

}

func IsValidIdentifier(s string) bool {

	for _, r := range s {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '/' || r == '-') {
			return false
		}
	}
	return true
}
