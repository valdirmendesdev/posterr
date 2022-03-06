package posts

import (
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type GetByIDService struct {
	postRepo posts.Repository
}

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Post *posts.Post
}

func NewGetByIDService(postRepo posts.Repository) *GetByIDService {
	return &GetByIDService{
		postRepo: postRepo,
	}
}

func (s *GetByIDService) Perform(request GetByIDRequest) (*GetByIDResponse, error) {

	uuid, err := types.ParseUUID(request.ID)
	if err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetByID(uuid)
	if err != nil {
		return nil, err
	}

	return &GetByIDResponse{
		Post: post,
	}, nil
}
