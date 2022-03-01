package profiles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/profiles"
	friendships_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/friendships"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

var userRepo *users_infra.MemoryRepository
var friendshipRepo *friendships_infra.MemoryRepository

func createUser(t *testing.T) *users.User {
	u, err := users.NewUser(
		types.NewUUID(),
		users.Username("username"),
		time.Now(),
	)
	assert.NoError(t, err)
	return u
}

func setupRepos() {
	userRepo = users_infra.NewMemoryRepository()
	friendshipRepo = friendships_infra.NewMemoryRepository()
}

func setupTest() *profiles.GetByUsernameService {
	setupRepos()
	return profiles.NewGetByUsernameService(userRepo, friendshipRepo)
}

func checkProfile(t *testing.T, u *users.User, p *profiles.GetByUsernameResponse, followersNumber, followingNumber int) {
	assert.Equal(t, u.ID, p.ID)
	assert.Equal(t, u.Username, p.Username)
	assert.Equal(t, u.JoinedAt, p.JoinedAt)
	assert.Equal(t, followersNumber, p.Followers)
	assert.Equal(t, followingNumber, p.Following)
}

func TestGetByUsername(t *testing.T) {
	s := setupTest()
	u := createUser(t)
	userRepo.Add(u)

	p, err := s.Perform(profiles.GetByUsernameRequest{
		Username: "username",
	})

	assert.NotNil(t, p)
	assert.NoError(t, err)
	checkProfile(t, u, p, 0, 0)
}

func TestGetByUsername_user_not_found(t *testing.T) {
	s := setupTest()

	_, err := s.Perform(profiles.GetByUsernameRequest{
		Username: "notfound",
	})

	assert.ErrorIs(t, err, users.ErrNotFound)
}

func TestGetByUsername_invalid_username(t *testing.T) {
	s := setupTest()

	_, err := s.Perform(profiles.GetByUsernameRequest{
		Username: "user#$%invalid",
	})

	assert.ErrorIs(t, err, users.ErrInvalidUsername)
}

func TestGetByUsername_without_friendships(t *testing.T) {
	s := setupTest()
	u := createUser(t)
	userRepo.Add(u)

	p, err := s.Perform(profiles.GetByUsernameRequest{
		Username: "username",
	})

	assert.NotNil(t, p)
	assert.NoError(t, err)
	checkProfile(t, u, p, 0, 0)
}

func TestGetByUsername_with_friendships(t *testing.T) {
	s := setupTest()

	u := createUser(t)
	userRepo.Add(u)

	lu := utils.GetLoggedUser()

	f, err := friendships.NewFriendship(u, lu)
	assert.NoError(t, err)
	err = friendshipRepo.Insert(f)
	assert.NoError(t, err)

	f, err = friendships.NewFriendship(lu, u)
	assert.NoError(t, err)
	err = friendshipRepo.Insert(f)
	assert.NoError(t, err)

	p, err := s.Perform(profiles.GetByUsernameRequest{
		Username: "username",
	})

	assert.NotNil(t, p)
	assert.NoError(t, err)
	checkProfile(t, u, p, 1, 1)
}
