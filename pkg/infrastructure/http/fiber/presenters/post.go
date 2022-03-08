package presenters

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type Post struct {
	ID             types.UUID `json:"id"`
	Username       string     `json:"username"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	CreatedAt      time.Time  `json:"created_at"`
	ReferencedPost *Post      `json:"referenced_post"`
}

type CreatePost struct {
	Content string `json:"content"`
}

func TransformPostToPresenter(post *posts.Post) (Post, error) {
	p := Post{
		ID:        post.ID,
		Username:  post.User.Username.String(),
		Content:   post.Content,
		Type:      post.Type,
		CreatedAt: post.CreatedAt,
	}

	if post.ReferencedPost != nil {
		referencedPost, err := TransformPostToPresenter(post.ReferencedPost)
		if err != nil {
			return p, err
		}
		p.ReferencedPost = &referencedPost
	}

	return p, nil
}
