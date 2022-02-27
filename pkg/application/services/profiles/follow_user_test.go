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

func TestFollowUser(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewFollowUserService(uR, uF)

	user, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("myuser"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(user)
	assert.NoError(t, err)

	follower, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("follower"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(follower)
	assert.NoError(t, err)

	err = s.Perform(profiles.FollowUserRequest{
		UserUsername:     user.Username.String(),
		FollowerUsername: follower.Username.String(),
	})
	assert.NoError(t, err)

	f, err := uF.Get(user.ID, follower.ID)
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestFollowUser_FriendshipAlreadyExists(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewFollowUserService(uR, uF)

	user, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("myuser"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(user)
	assert.NoError(t, err)

	follower, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("follower"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(follower)
	assert.NoError(t, err)

	followUserRequest := profiles.FollowUserRequest{
		UserUsername:     user.Username.String(),
		FollowerUsername: follower.Username.String(),
	}

	err = s.Perform(followUserRequest)
	assert.NoError(t, err)

	err = s.Perform(followUserRequest)
	assert.ErrorIs(t, friendships.ErrAlreadyExists, err)
}

func TestFollowUser_UserNotExist(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewFollowUserService(uR, uF)

	err := s.Perform(profiles.FollowUserRequest{
		UserUsername:     "myusernotfound",
		FollowerUsername: "follower",
	})
	assert.ErrorIs(t, profiles.ErrUserNotExist, err)
}

func TestFollowUser_FollowerNotExist(t *testing.T) {
	uR := users.NewMemoryRepository()
	uF := friendships_infra.NewMemoryRepository()
	s := profiles.NewFollowUserService(uR, uF)

	user, err := users_domain.NewUser(types.NewUUID(), users_domain.Username("myuser"), time.Now())
	assert.NoError(t, err)
	err = uR.Add(user)
	assert.NoError(t, err)

	err = s.Perform(profiles.FollowUserRequest{
		UserUsername:     "myuser",
		FollowerUsername: "followernotfound",
	})
	assert.ErrorIs(t, profiles.ErrFollowerNotExist, err)
}
