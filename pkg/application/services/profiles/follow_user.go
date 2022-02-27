package profiles

import (
	"errors"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
)

type FollowUserService struct {
	userRepo       users.Repository
	friendshipRepo friendships.Repository
}

func NewFollowUserService(userRepo users.Repository, friendshipRepo friendships.Repository) *FollowUserService {
	return &FollowUserService{
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
	}
}

var (
	ErrUserNotExist     = errors.New("user not exist")
	ErrFollowerNotExist = errors.New("follower not exist")
)

type FollowUserRequest struct {
	UserUsername     string
	FollowerUsername string
}

func (s *FollowUserService) Perform(request FollowUserRequest) error {
	user, err := s.userRepo.GetByUsername(users.Username(request.UserUsername))
	if err != nil {
		switch err {
		case users.ErrNotFound:
			return ErrUserNotExist
		default:
			return err
		}
	}

	follower, err := s.userRepo.GetByUsername(users.Username(request.FollowerUsername))
	if err != nil {
		switch err {
		case users.ErrNotFound:
			return ErrFollowerNotExist
		default:
			return err
		}
	}

	f, err := friendships.NewFriendship(user, follower)
	if err != nil {
		return err
	}

	err = s.friendshipRepo.Insert(f)
	if err != nil {
		return err
	}

	return nil
}
