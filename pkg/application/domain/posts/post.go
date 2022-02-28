package posts

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type Post struct {
	ID        types.UUID
	User      *users.User
	Content   string
	CreatedAt time.Time
}

func NewPost(id types.UUID, user *users.User, content string, createdAt time.Time) (*Post, error) {
	p := &Post{
		ID:        id,
		User:      user,
		Content:   content,
		CreatedAt: createdAt,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Post) Validate() error {
	return nil
}
