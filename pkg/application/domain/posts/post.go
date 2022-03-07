package posts

import (
	"errors"
	"strconv"
	"time"

	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

const ContentMaxLength = 777

const (
	TypePost   = "post"
	TypeRepost = "repost"
	TypeQuote  = "quote"
)

var (
	ErrInvalidUser           = errors.New("invalid user")
	ErrInvalidCreatedAt      = errors.New("created at not specified")
	ErrContentOverMaxLength  = errors.New("over max length " + strconv.Itoa(ContentMaxLength) + " characters")
	ErrInvalidReferencedPost = errors.New("invalid referenced post")
	ErrInvalidContent        = errors.New("content not specified")
)

type Post struct {
	ID             types.UUID
	User           *users.User
	Content        string
	Type           string
	ReferencedPost *Post
	CreatedAt      time.Time
}

func NewPost(id types.UUID, user *users.User, content string, createdAt time.Time) (*Post, error) {
	p := &Post{
		ID:        id,
		User:      user,
		Content:   content,
		Type:      TypePost,
		CreatedAt: createdAt,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func NewRepost(id types.UUID, user *users.User, referencedPost *Post, createdAt time.Time) (*Post, error) {
	p, err := NewPost(id, user, "", createdAt)
	if err != nil {
		return nil, err
	}
	p.Type = TypeRepost
	p.ReferencedPost = referencedPost

	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func NewQuotePost(id types.UUID, user *users.User, referencedPost *Post, content string, createdAt time.Time) (*Post, error) {
	p, err := NewPost(id, user, content, createdAt)
	if err != nil {
		return nil, err
	}
	p.Type = TypeQuote
	p.ReferencedPost = referencedPost

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
		return ErrContentOverMaxLength
	}

	if p.CreatedAt.IsZero() {
		return ErrInvalidCreatedAt
	}

	switch p.Type {
	case TypeRepost, TypeQuote:
		if p.ReferencedPost == nil {
			return ErrInvalidReferencedPost
		}
	}

	if p.Type == TypeQuote && len(p.Content) == 0 {
		return ErrInvalidContent
	}

	return nil
}
