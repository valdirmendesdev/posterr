package posts

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type QuoteService struct {
	userRepo users.Repository
	postRepo posts.Repository
}

type QuoteRequest struct {
	PostID   string
	Username string
	Content  string
}

type QuoteResponse struct {
	Post *posts.Post
}

func NewQuoteService(userRepo users.Repository, postRepo posts.Repository) *QuoteService {
	return &QuoteService{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *QuoteService) Perform(request QuoteRequest) (*QuoteResponse, error) {

	postID, err := types.ParseUUID(request.PostID)
	if err != nil {
		return nil, err
	}

	username := users.Username(request.Username)
	if err := username.Validate(); err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	u, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	quote, err := posts.NewQuotePost(types.NewUUID(), u, post, request.Content, time.Now())
	if err != nil {
		return nil, err
	}

	err = s.postRepo.Insert(quote)
	if err != nil {
		return nil, err
	}

	return &QuoteResponse{
		Post: quote,
	}, nil
}
