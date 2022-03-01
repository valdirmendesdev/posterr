package profiles

import (
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
)

type CheckFriendshipService struct {
	userRepo       users.Repository
	friendshipRepo friendships.Repository
}

type CheckFriendshipRequest struct {
	Username         string
	FollowerUsername string
}

type CheckFriendshipResponse struct {
	Friendship *friendships.Friendship
}

func NewCheckFriendshipService(userRepo users.Repository, friendshipRepo friendships.Repository) *CheckFriendshipService {
	return &CheckFriendshipService{
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
	}
}

func (s *CheckFriendshipService) Perform(request CheckFriendshipRequest) (*CheckFriendshipResponse, error) {
	username := users.Username(request.Username)
	if err := username.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	followerUsername := users.Username(request.FollowerUsername)
	if err := followerUsername.Validate(); err != nil {
		return nil, err
	}

	follower, err := s.userRepo.GetByUsername(followerUsername)
	if err != nil {
		return nil, err
	}

	f, err := s.friendshipRepo.Get(u.ID, follower.ID)
	if err != nil && err != friendships.ErrNotFound {
		return nil, err
	}

	return &CheckFriendshipResponse{
		Friendship: f,
	}, nil
}
