package friendships

import (
	"errors"

	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var (
	ErrNotFound      = errors.New("friendship not found")
	ErrAlreadyExists = errors.New("friendship already exists")
)

type Repository interface {
	Insert(f *Friendship) error
	Delete(f *Friendship) error
	Get(userID, followerID types.UUID) (*Friendship, error)
	GetFollowersNumber(userID types.UUID) (int, error)
	GetFollowingNumber(userID types.UUID) (int, error)
}
