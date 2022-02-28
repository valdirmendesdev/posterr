package presenters

import (
	"time"

	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type Profile struct {
	ID        types.UUID `json:"id"`
	Username  string     `json:"username"`
	JoinedAt  time.Time  `json:"joined_at"`
	Followers int        `json:"followers"`
	Following int        `json:"following"`
}
