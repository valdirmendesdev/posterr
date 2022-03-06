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

func (r *MemoryRepository) ListByUserID(userID types.UUID) ([]*posts.Post, error) {
	posts := make([]*posts.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.User.ID == userID {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (r *MemoryRepository) FindByUsernameAndDate(username string, date time.Time) ([]*posts.Post, error) {
	un := users.Username(username)
	posts := make([]*posts.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.User.Username != un {
			continue
		}
		if post.CreatedAt.Day() == date.Day() &&
			post.CreatedAt.Month() == date.Month() &&
			post.CreatedAt.Year() == date.Year() {
			posts = append(posts, post)
		}
	}
	return posts, nil
}
