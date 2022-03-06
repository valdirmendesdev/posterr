package posts

import (
	"errors"
	"time"

	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var (
	ErrNotFound = errors.New("post not found")
)

type Repository interface {
	Insert(post *Post) error
	GetByID(id types.UUID) (*Post, error)
	List() ([]*Post, error)
	ListByUserID(userID types.UUID) ([]*Post, error)
	FindByUsernameAndDate(username string, date time.Time) ([]*Post, error)
}
