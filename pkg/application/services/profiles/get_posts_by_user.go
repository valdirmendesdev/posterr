package profiles

import (
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
)

type GetPostsByUserService struct {
	userRepo users.Repository
	postRepo posts.Repository
}

type GetPostsByUserRequest struct {
	Username string
}

type GetPostsByUserResponse struct {
	Posts []*posts.Post
}

func NewGetPostsByUserService(userRepo users.Repository, postRepo posts.Repository) *GetPostsByUserService {
	return &GetPostsByUserService{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *GetPostsByUserService) Perform(request GetPostsByUserRequest) (*GetPostsByUserResponse, error) {
	username := users.Username(request.Username)
	if err := username.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	posts, err := s.postRepo.ListByUserID(u.ID)
	if err != nil {
		return nil, err
	}

	return &GetPostsByUserResponse{
		Posts: posts,
	}, nil
}
