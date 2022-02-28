package posts

import (
	"errors"
	"strconv"
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

const ContentMaxLength = 777

var (
	ErrInvalidUser      = errors.New("invalid user")
	ErrInvalidCreatedAt = errors.New("created at not specified")
	ErrOverMaxLength    = errors.New("over max length " + strconv.Itoa(ContentMaxLength))
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
	if p.User == nil {
		return ErrInvalidUser
	}

	if len(p.Content) > ContentMaxLength {
		return ErrOverMaxLength
	}

	if p.CreatedAt.IsZero() {
		return ErrInvalidCreatedAt
	}
	return nil
}
