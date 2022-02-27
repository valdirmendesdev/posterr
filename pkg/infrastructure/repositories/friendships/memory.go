package friendships

import (
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

type MemoryRepository struct {
	friendships []*friendships.Friendship
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		friendships: []*friendships.Friendship{},
	}
}

func (r *MemoryRepository) Insert(f *friendships.Friendship) error {
	ff, err := r.Get(f.User.ID, f.Follower.ID)
	if err != nil {

		switch err {
		case friendships.ErrNotFound:
			r.friendships = append(r.friendships, f)
		default:
			return err
		}
	}
	if ff != nil {
		return friendships.ErrAlreadyExists
	}
	return nil
}

func (r *MemoryRepository) Delete(f *friendships.Friendship) error {
	for i, ff := range r.friendships {
		if ff.User.ID == f.User.ID && ff.Follower.ID == f.Follower.ID {
			r.friendships[i] = r.friendships[len(r.friendships)-1]
			r.friendships = r.friendships[:len(r.friendships)-1]
			return nil
		}
	}
	return friendships.ErrNotFound
}

func (r *MemoryRepository) Get(userID, followerID types.UUID) (*friendships.Friendship, error) {
	for _, f := range r.friendships {
		if f.User.ID == userID && f.Follower.ID == followerID {
			return f, nil
		}
	}

	return nil, friendships.ErrNotFound
}

func (r *MemoryRepository) GetFollowersNumber(userID types.UUID) (int, error) {
	followers := 0
	for _, uf := range r.friendships {
		if uf.User.ID == userID {
			followers += 1
		}
	}
	return followers, nil
}

func (r *MemoryRepository) GetFollowingNumber(userID types.UUID) (int, error) {
	following := 0
	for _, uf := range r.friendships {
		if uf.Follower.ID == userID {
			following += 1
		}
	}
	return following, nil
}
