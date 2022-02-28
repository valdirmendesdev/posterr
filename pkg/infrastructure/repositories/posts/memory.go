package posts

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type MemoryRepository struct {
	posts map[types.UUID]*posts.Post
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		posts: make(map[types.UUID]*posts.Post),
	}
}

func (r *MemoryRepository) Insert(post *posts.Post) error {
	r.posts[post.ID] = post
	return nil
}

func (r *MemoryRepository) List() ([]*posts.Post, error) {
	posts := make([]*posts.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *MemoryRepository) ListByUsername(username string) ([]*posts.Post, error) {
	un := users.Username(username)
	posts := make([]*posts.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.User.Username == un {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *MemoryRepository) ListByDate(date time.Time) ([]*posts.Post, error) {
	posts := make([]*posts.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.CreatedAt.Day() == date.Day() &&
			post.CreatedAt.Month() == date.Month() &&
			post.CreatedAt.Year() == date.Year() {
			posts = append(posts, post)
		}
	}
	return posts, nil
}
