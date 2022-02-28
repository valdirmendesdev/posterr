package posts

import "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"

type GetAllPostsService struct {
	postRepo posts.Repository
}

type GetAllPostResponse struct {
	Posts []*posts.Post
}

func NewGetAllPostsService(postRepo posts.Repository) *GetAllPostsService {
	return &GetAllPostsService{
		postRepo: postRepo,
	}
}

func (s *GetAllPostsService) Perform() (*GetAllPostResponse, error) {
	posts, err := s.postRepo.List()
	if err != nil {
		return nil, err
	}
	return &GetAllPostResponse{
		Posts: posts,
	}, nil
}
