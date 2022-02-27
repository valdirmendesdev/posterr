package profiles

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type getByUsernameService struct {
	userRepo users.Repository
}

func NewGetByUsernameService(userRepo users.Repository) *getByUsernameService {
	return &getByUsernameService{
		userRepo: userRepo,
	}
}

type GetByUsernameRequest struct {
	Username string
}

type GetByUsernameResponse struct {
	ID       types.UUID
	Username users.Username
	JoinedAt time.Time
}

func (s *getByUsernameService) Perform(request GetByUsernameRequest) (*GetByUsernameResponse, error) {
	username := users.Username(request.Username)
	if err := username.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUsername(username)

	if err != nil {
		return nil, err
	}

	uD := &GetByUsernameResponse{
		ID:       u.ID,
		Username: u.Username,
		JoinedAt: u.JoinedAt,
	}

	return uD, nil
}
