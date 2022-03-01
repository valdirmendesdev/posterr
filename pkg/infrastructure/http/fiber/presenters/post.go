package presenters

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type Post struct {
	ID        types.UUID `json:"id"`
	Username  string     `json:"username"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
}
