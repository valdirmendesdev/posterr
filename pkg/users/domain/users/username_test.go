package users_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/users/domain/users"
)

func Test_ValidUsername(t *testing.T) {
	u := users.Username("username")
	assert.Nil(t, u.Validate())
	assert.Equal(t, "username", u.String())
}

func Test_InvalidUsernameMinimumLength(t *testing.T) {
	u := users.Username("")
	err := u.Validate()
	assert.Error(t, err, users.ErrInvalidUsername)
}

func Test_InvalidUsernameMaximumLength(t *testing.T) {
	u := users.Username("usernameTooLongPlus14Characters")
	err := u.Validate()
	assert.Error(t, err, users.ErrInvalidUsername)
}

func Test_InvalidUsernameWithInvalidCharacters(t *testing.T) {
	u := users.Username("user!@#$%^&*()")
	err := u.Validate()
	assert.Error(t, err, users.ErrInvalidUsername)
}
