package friendships

import (
	"errors"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
)

type Friendship struct {
	User     *users.User
	Follower *users.User
}

var (
	ErrFriendshipInvalidUser     = errors.New("invalid user")
	ErrFriendshipInvalidFollower = errors.New("invalid follower")
	ErrFriendshipWithYourself    = errors.New("user and follower are the same")
)

func NewFriendship(user, follower *users.User) (*Friendship, error) {
	f := &Friendship{
		User:     user,
		Follower: follower,
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Friendship) Validate() error {
	if f.User == nil {
		return ErrFriendshipInvalidUser
	}
	if f.Follower == nil {
		return ErrFriendshipInvalidFollower
	}

	if f.User.ID == f.Follower.ID {
		return ErrFriendshipWithYourself
	}
	return nil
}
