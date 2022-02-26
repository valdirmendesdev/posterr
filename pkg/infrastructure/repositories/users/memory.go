package users

import (
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type MemoryRepository struct {
	users map[types.UUID]*users.User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[types.UUID]*users.User),
	}
}

func (r *MemoryRepository) GetByUsername(un users.Username) (*users.User, error) {
	for _, u := range r.users {
		if u.Username == un {
			return u, nil
		}
	}

	return nil, users.ErrNotFound
}

func (r *MemoryRepository) Add(u *users.User) error {
	r.users[u.ID] = u
	return nil
}
