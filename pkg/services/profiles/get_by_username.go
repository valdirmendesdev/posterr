package profiles

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/domain/users"
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
	Username users.Username
}

type GetByUsernameResponse struct {
	ID       types.UUID
	Username users.Username
	JoinedAt time.Time
}

func (s *getByUsernameService) Perform(request GetByUsernameRequest) (*GetByUsernameResponse, error) {
	u, err := s.userRepo.GetByUsername(request.Username)

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
