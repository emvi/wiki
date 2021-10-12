package mail

import (
	"regexp"
)

const (
	emailReg = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

// EmailValid checks if an email is syntactically valid.
func EmailValid(email string) bool {
	if len(email) == 0 || len(email) > 255 {
		return false
	}

	reg := regexp.MustCompile(emailReg)

	if !reg.Match([]byte(email)) {
		return false
	}

	return true
}
