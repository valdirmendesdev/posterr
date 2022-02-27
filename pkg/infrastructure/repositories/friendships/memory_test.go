package friendships_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	friendships_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/friendships"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var user, follower *users.User

func createUserFollower(t *testing.T) (*users.User, *users.User) {
	user, err := users.NewUser(types.NewUUID(), users.Username("myuser"), time.Now())
	assert.NoError(t, err)
	follower, err := users.NewUser(types.NewUUID(), users.Username("follower"), time.Now())
	assert.NoError(t, err)
	return user, follower
}

func createFriendship(t *testing.T, u, f *users.User) *friendships_domain.Friendship {
	friendship, err := friendships_domain.NewFriendship(u, f)
	assert.NoError(t, err)
	return friendship
}

func setupTest(t *testing.T) *friendships.MemoryRepository {
	user, follower = createUserFollower(t)
	r := friendships.NewMemoryRepository()
	return r
}

func TestCreateNewFriendship(t *testing.T) {
	r := setupTest(t)
	f := createFriendship(t, user, follower)
	err := r.Insert(f)
	assert.NoError(t, err)

	f2, err := r.Get(user.ID, follower.ID)
	assert.NoError(t, err)
	assert.Equal(t, f, f2)
}

func TestFriendshipAlreadyExists(t *testing.T) {
	r := setupTest(t)
	f := createFriendship(t, user, follower)
	err := r.Insert(f)
	assert.NoError(t, err)

	err = r.Insert(f)
	assert.ErrorIs(t, friendships_domain.ErrAlreadyExists, err)
}

func TestDeleteFriendship(t *testing.T) {
	r := setupTest(t)
	f := createFriendship(t, user, follower)
	err := r.Insert(f)
	assert.NoError(t, err)

	err = r.Delete(f)
	assert.NoError(t, err)

	f2, err := r.Get(user.ID, follower.ID)
	assert.Nil(t, f2)
	assert.ErrorIs(t, friendships_domain.ErrNotFound, err)
}

func TestDeleteFriendship_NotFound(t *testing.T) {
	r := setupTest(t)
	f := createFriendship(t, user, follower)

	err := r.Delete(f)
	assert.ErrorIs(t, friendships_domain.ErrNotFound, err)
}

func TestGetFollowersNumber(t *testing.T) {
	r := setupTest(t)

	n, err := r.GetFollowersNumber(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	f := createFriendship(t, user, follower)
	err = r.Insert(f)
	assert.NoError(t, err)

	n, err = r.GetFollowersNumber(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestGetFollowingNumber(t *testing.T) {
	r := setupTest(t)

	n, err := r.GetFollowingNumber(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	f := createFriendship(t, follower, user)
	err = r.Insert(f)
	assert.NoError(t, err)

	n, err = r.GetFollowingNumber(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}
