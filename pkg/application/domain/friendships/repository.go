package friendships

import "github.com/valdirmendesdev/posterr/pkg/shared/types"

type Repository interface {
	Insert(f *Friendship) error
	Delete(f *Friendship) error
	Get(userID, followerID types.UUID) (*Friendship, error)
	GetFollowersNumber(userID types.UUID) (int, error)
	GetFollowingNumber(userID types.UUID) (int, error)
}
