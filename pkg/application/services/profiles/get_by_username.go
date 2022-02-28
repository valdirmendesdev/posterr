package profiles

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type GetByUsernameService struct {
	userRepo       users.Repository
	friendshipRepo friendships.Repository
}

func NewGetByUsernameService(userRepo users.Repository, friendshipRepo friendships.Repository) *GetByUsernameService {
	return &GetByUsernameService{
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
	}
}

type GetByUsernameRequest struct {
	Username string
}

type GetByUsernameResponse struct {
	ID        types.UUID
	Username  users.Username
	JoinedAt  time.Time
	Followers int
	Following int
}

func (s *GetByUsernameService) Perform(request GetByUsernameRequest) (*GetByUsernameResponse, error) {
	username := users.Username(request.Username)
	if err := username.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	followers, err := s.friendshipRepo.GetFollowersNumber(u.ID)
	if err != nil {
		return nil, err
	}

	following, err := s.friendshipRepo.GetFollowingNumber(u.ID)
	if err != nil {
		return nil, err
	}

	uD := &GetByUsernameResponse{
		ID:        u.ID,
		Username:  u.Username,
		JoinedAt:  u.JoinedAt,
		Followers: followers,
		Following: following,
	}

	return uD, nil
}
