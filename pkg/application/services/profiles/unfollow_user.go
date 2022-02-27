package profiles

import (
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
)

type UnfollowUserService struct {
	userRepo       users.Repository
	friendshipRepo friendships.Repository
}

func NewUnfollowUserService(userRepo users.Repository, friendshipRepo friendships.Repository) *UnfollowUserService {
	return &UnfollowUserService{
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
	}
}

type UnfollowUserRequest struct {
	UserUsername     string
	FollowerUsername string
}

func (s *UnfollowUserService) Perform(request UnfollowUserRequest) error {
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

	f, err := s.friendshipRepo.Get(user.ID, follower.ID)
	if err != nil {
		return err
	}

	err = s.friendshipRepo.Delete(f)
	if err != nil {
		return err
	}

	return nil
}
