package posts

import (
	"errors"
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var (
	ErrDailyPostsLimitReached = errors.New("daily posts limit reached")
)

const PostsLimitByDay = 5

type CreateNewPostService struct {
	postRepo posts.Repository
}

type CreateNewPostRequest struct {
	User    *users.User
	Content string
}

type CreateNewPostResponse struct {
	Post *posts.Post
}

func NewCreateNewPostService(postRepo posts.Repository) *CreateNewPostService {
	return &CreateNewPostService{
		postRepo: postRepo,
	}
}

func (s *CreateNewPostService) Perform(request CreateNewPostRequest) (*CreateNewPostResponse, error) {

	dayPosts, err := s.postRepo.FindByUsernameAndDate(request.User.Username.String(), time.Now())
	if err != nil {
		return nil, err
	}

	if len(dayPosts) == PostsLimitByDay {
		return nil, ErrDailyPostsLimitReached
	}

	post, err := posts.NewPost(types.NewUUID(), request.User, request.Content, time.Now())
	if err != nil {
		return nil, err
	}
	err = s.postRepo.Insert(post)
	if err != nil {
		return nil, err
	}
	return &CreateNewPostResponse{
		Post: post,
	}, nil
}
