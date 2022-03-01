package posts

import "time"

type Repository interface {
	Insert(post *Post) error
	List() ([]*Post, error)
	ListByUsername(username string) ([]*Post, error)
	FindByUsernameAndDate(username string, date time.Time) ([]*Post, error)
}
