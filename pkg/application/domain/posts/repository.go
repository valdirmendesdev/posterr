package posts

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type Repository interface {
	Insert(post *Post) error
	List() ([]*Post, error)
	ListByUserID(userID types.UUID) ([]*Post, error)
	FindByUsernameAndDate(username string, date time.Time) ([]*Post, error)
}
