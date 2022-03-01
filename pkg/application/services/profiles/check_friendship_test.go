package profiles_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/profiles"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

func createFriendship(t *testing.T, user, follower *users.User) {
	f, err := friendships.NewFriendship(user, follower)
	assert.NoError(t, err)
	err = friendshipRepo.Insert(f)
	assert.NoError(t, err)
}

func TestCheckFriendship_friendship_not_found(t *testing.T) {
	setupRepos()
	u := createUser(t)
	lu := utils.GetLoggedUser()
	err := userRepo.Add(u)
	assert.NoError(t, err)
	err = userRepo.Add(lu)
	assert.NoError(t, err)
	s := profiles.NewCheckFriendshipService(userRepo, friendshipRepo)
	response, err := s.Perform(profiles.CheckFriendshipRequest{
		Username:         "username",
		FollowerUsername: lu.Username.String(),
	})
	assert.NoError(t, err)
	assert.Nil(t, response.Friendship)
}

func TestCheckFriendship(t *testing.T) {
	setupRepos()
	u := createUser(t)
	lu := utils.GetLoggedUser()

	err := userRepo.Add(u)
	assert.NoError(t, err)

	err = userRepo.Add(lu)
	assert.NoError(t, err)
	createFriendship(t, u, lu)

	s := profiles.NewCheckFriendshipService(userRepo, friendshipRepo)
	response, err := s.Perform(profiles.CheckFriendshipRequest{
		Username:         "username",
		FollowerUsername: lu.Username.String(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, response.Friendship)
	assert.Equal(t, u.ID, response.Friendship.User.ID)
	assert.Equal(t, lu.ID, response.Friendship.Follower.ID)
}
