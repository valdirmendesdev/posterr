package profiles

import (
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

type FollowUserRequest struct {
	UserUsername     string
	FollowerUsername string
}

func (s *FollowUserService) Perform(request FollowUserRequest) error {
	user, err := s.userRepo.GetByUsername(users.Username(request.UserUsername))
	if err != nil {
		return err
	}

	follower, err := s.userRepo.GetByUsername(users.Username(request.FollowerUsername))
	if err != nil {
		return err
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
