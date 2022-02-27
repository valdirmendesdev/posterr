package profiles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	users_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/profiles"
	friendships_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/friendships"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

func TestUnfollowUser(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewUnfollowUserService(uR, uF)

	user, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("myuser"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(user)
	assert.NoError(t, err)

	follower, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("follower"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(follower)
	assert.NoError(t, err)

	fs, err := friendships.NewFriendship(user, follower)
	assert.NoError(t, err)
	err = uF.Insert(fs)
	assert.NoError(t, err)

	f, err := uF.Get(user.ID, follower.ID)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	err = s.Perform(profiles.UnfollowUserRequest{
		UserUsername:     user.Username.String(),
		FollowerUsername: follower.Username.String(),
	})
	assert.NoError(t, err)

	f, err = uF.Get(user.ID, follower.ID)
	assert.Nil(t, f)
	assert.ErrorIs(t, friendships.ErrNotFound, err)
}

func TestUnfollowUser_FriendshipNotExists(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewUnfollowUserService(uR, uF)

	user, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("myuser"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(user)
	assert.NoError(t, err)

	follower, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("follower"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(follower)
	assert.NoError(t, err)

	unfollowUserRequest := profiles.UnfollowUserRequest{
		UserUsername:     user.Username.String(),
		FollowerUsername: follower.Username.String(),
	}

	err = s.Perform(unfollowUserRequest)
	assert.ErrorIs(t, friendships.ErrNotFound, err)
}

func TestUnfollowUser_UserNotExist(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewUnfollowUserService(uR, uF)

	unfollowUserRequest := profiles.UnfollowUserRequest{
		UserUsername:     "myusernotfound",
		FollowerUsername: "follower",
	}

	err := s.Perform(unfollowUserRequest)
	assert.ErrorIs(t, profiles.ErrUserNotExist, err)
}

func TestUnfollowUser_FollowerNotExist(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewUnfollowUserService(uR, uF)

	user, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("myuser"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(user)
	assert.NoError(t, err)

	unfollowUserRequest := profiles.UnfollowUserRequest{
		UserUsername:     user.Username.String(),
		FollowerUsername: "follower",
	}

	err = s.Perform(unfollowUserRequest)
	assert.ErrorIs(t, profiles.ErrFollowerNotExist, err)
}
