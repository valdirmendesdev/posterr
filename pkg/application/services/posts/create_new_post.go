package posts

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type CreateNewPostService struct {
	postRepo posts.Repository
}

type CreateNewPostRequest struct {
	User    *users.User
	Content string
}

type CreateNewPostResponse struct {
	ID        types.UUID
	Username  string
	Content   string
	CreatedAt time.Time
}

func NewCreateNewPostService(postRepo posts.Repository) *CreateNewPostService {
	return &CreateNewPostService{
		postRepo: postRepo,
	}
}

func (s *CreateNewPostService) Perform(request CreateNewPostRequest) (*CreateNewPostResponse, error) {
	post, err := posts.NewPost(types.NewUUID(), request.User, request.Content, time.Now())
	if err != nil {
		return nil, err
	}
	err = s.postRepo.Insert(post)
	if err != nil {
		return nil, err
	}
	return &CreateNewPostResponse{
		ID:        post.ID,
		Username:  post.User.Username.String(),
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}, nil
}
