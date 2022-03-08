package posts

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type RepostService struct {
	userRepo users.Repository
	postRepo posts.Repository
}

type RepostRequest struct {
	PostID   string
	Username string
}

type RepostResponse struct {
	Post *posts.Post
}

func NewRepostService(userRepo users.Repository, postRepo posts.Repository) *RepostService {
	return &RepostService{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *RepostService) Perform(request RepostRequest) (*RepostResponse, error) {

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

	repost, err := posts.NewRepost(types.NewUUID(), u, post, time.Now())
	if err != nil {
		return nil, err
	}

	err = s.postRepo.Insert(repost)
	if err != nil {
		return nil, err
	}

	return &RepostResponse{
		Post: repost,
	}, nil
}
