package users

import (
	"errors"
	"regexp"
)

type Username string

var (
	ErrInvalidUsername = errors.New("invalid username")
)

func (u Username) Validate() error {
	const minLength = 1
	const maxLength = 14
	const alphanumericCharactersPattern = "[a-zA-Z0-9]"

	if len(u) < minLength || len(u) > maxLength {
		return ErrInvalidUsername
	}

	acceptableUsernameRegex := regexp.MustCompile(`^` + alphanumericCharactersPattern + `+$`)
	if !acceptableUsernameRegex.MatchString(string(u)) {
		return ErrInvalidUsername
	}

	return nil
}

func (u Username) String() string {
	return string(u)
}
