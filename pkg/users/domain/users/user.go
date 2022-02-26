package users

import (
	"errors"
	"time"

	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var (
	ErrInvalidJoinedAt = errors.New("invalid joined at")
)

type User struct {
	ID       types.UUID
	Username Username
	JoinedAt time.Time
}

func NewUser(id types.UUID, username Username, joinedAt time.Time) (*User, error) {
	u := &User{
		ID:       id,
		Username: username,
		JoinedAt: joinedAt,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	if err := u.Username.Validate(); err != nil {
		return err
	}

	if u.JoinedAt.IsZero() {
		return ErrInvalidJoinedAt
	}
	return nil
}
