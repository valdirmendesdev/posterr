package users_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
	"github.com/valdirmendesdev/posterr/pkg/users/domain/users"
)

func TestNewUser(t *testing.T) {

	now := time.Now()
	u, err := users.NewUser(types.NewUUID(), users.Username("username"), now)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.Equal(t, "username", u.Username.String())
	assert.NotNil(t, u.JoinedAt)
	assert.Equal(t, now, u.JoinedAt)
}

func TestNewUser_invalid_username(t *testing.T) {
	_, err := users.NewUser(types.NewUUID(), users.Username("user!@#$%^&*()"), time.Now())
	assert.ErrorIs(t, err, users.ErrInvalidUsername)
}

func TestNewUser_invalid_joinedAt(t *testing.T) {
	var now time.Time
	_, err := users.NewUser(types.NewUUID(), users.Username("username"), now)
	assert.ErrorIs(t, err, users.ErrInvalidJoinedAt)
}
