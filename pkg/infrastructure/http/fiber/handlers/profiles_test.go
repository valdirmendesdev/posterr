package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/handlers"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	friendships_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/friendships"
	posts_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/posts"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	shared_presenters "github.com/valdirmendesdev/posterr/pkg/shared/infrastructure/http/presenters"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

var userRepo *users_infra.MemoryRepository
var friendshipRepo *friendships_infra.MemoryRepository

var user *users.User
var loggedUser *users.User

func setupRoutes() *handlers.ProfileRoutesConfig {
	userRepo = users_infra.NewMemoryRepository()
	friendshipRepo = friendships_infra.NewMemoryRepository()
	postsRepo = posts_infra.NewMemoryRepository()
	config := handlers.NewProfileRoutesConfigs(fiber.New(), userRepo, friendshipRepo, postsRepo)
	handlers.MountProfilesRoutes(config)
	return config
}

func createUser(t *testing.T) {
	var err error
	user, err = users.NewUser(types.NewUUID(), users.Username("anyuser"), time.Now())
	assert.NoError(t, err)
	err = userRepo.Add(user)
	assert.NoError(t, err)
}

func createLoggedUser(t *testing.T) {
	loggedUser = utils.GetLoggedUser()
	err := userRepo.Add(loggedUser)
	assert.NoError(t, err)
}

func createFriendship(t *testing.T, user, follower *users.User) *friendships.Friendship {
	f, err := friendships.NewFriendship(user, follower)
	assert.NoError(t, err)
	err = friendshipRepo.Insert(f)
	assert.NoError(t, err)
	return f
}

func decodeError(t *testing.T, response *http.Response) (error *shared_presenters.Error) {
	err := json.NewDecoder(response.Body).Decode(&error)
	assert.NoError(t, err)
	return
}

func createGetProfileRequest(username string) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/profiles/"+username, nil)
}

func checkProfile(t *testing.T, expected, actual *presenters.Profile) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.Equal(t, expected.JoinedAt.UTC().GoString(), actual.JoinedAt.UTC().GoString())
	assert.Equal(t, expected.Followers, actual.Followers)
	assert.Equal(t, expected.Following, actual.Following)
	assert.Equal(t, expected.IsFollowing, actual.IsFollowing)
}

func TestGetProfileRoute_without_friendships(t *testing.T) {
	routesConfig := setupRoutes()
	createUser(t)
	createLoggedUser(t)

	expected := &presenters.Profile{
		ID:          user.ID,
		Username:    user.Username.String(),
		JoinedAt:    user.JoinedAt,
		Followers:   0,
		Following:   0,
		IsFollowing: false,
	}

	request := createGetProfileRequest("anyuser")
	response, err := routesConfig.App.Test(request)
	var result *presenters.Profile
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.NoError(t, err)
	checkProfile(t, expected, result)
}

func TestGetProfileRoute_with_friendships(t *testing.T) {
	routesConfig := setupRoutes()
	createUser(t)
	createLoggedUser(t)

	//Logged user follows myusername
	createFriendship(t, user, loggedUser)

	//myusername follows Logged User
	createFriendship(t, loggedUser, user)

	expected := &presenters.Profile{
		ID:          user.ID,
		Username:    user.Username.String(),
		JoinedAt:    user.JoinedAt,
		Followers:   1,
		Following:   1,
		IsFollowing: true,
	}

	request := createGetProfileRequest("anyuser")
	response, err := routesConfig.App.Test(request)
	var result *presenters.Profile
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.NoError(t, err)
	checkProfile(t, expected, result)
}

func TestGetProfileRoute_invalid_username(t *testing.T) {
	routesConfig := setupRoutes()
	request := createGetProfileRequest("user@#$invalid")
	response, err := routesConfig.App.Test(request)

	expected := &shared_presenters.Error{
		Message: users.ErrInvalidUsername.Error(),
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	result := decodeError(t, response)
	assert.EqualValues(t, expected, result)
}

func TestGetProfileRoute_profile_not_found(t *testing.T) {
	routesConfig := setupRoutes()
	request := createGetProfileRequest("usernotfound")
	response, err := routesConfig.App.Test(request)

	expected := &shared_presenters.Error{
		Message: users.ErrNotFound.Error(),
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	result := decodeError(t, response)
	assert.EqualValues(t, expected, result)
}

func createFollowUserRequest(username string) *http.Request {
	return httptest.NewRequest(http.MethodPut, "/profiles/"+username+"/follow", nil)
}

func TestFollowUserRoute(t *testing.T) {
	rc := setupRoutes()
	createUser(t)
	createLoggedUser(t)

	request := createFollowUserRequest("anyuser")
	response, err := rc.App.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestFollowUserRoute_friendship_already_exist(t *testing.T) {
	rc := setupRoutes()
	createUser(t)
	createLoggedUser(t)
	createFriendship(t, user, loggedUser)

	request := createFollowUserRequest("anyuser")
	response, err := rc.App.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotModified, response.StatusCode)
}

func TestFollowUserRoute_bad_request(t *testing.T) {
	rc := setupRoutes()

	request := createFollowUserRequest("anyuser")
	response, err := rc.App.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	result := decodeError(t, response)
	assert.NotNil(t, result)
}

func createUnfollowUserRequest(username string) *http.Request {
	return httptest.NewRequest(http.MethodDelete, "/profiles/"+username+"/unfollow", nil)
}

func TestUnFollowUserRoute(t *testing.T) {
	rc := setupRoutes()
	createUser(t)
	createLoggedUser(t)
	createFriendship(t, user, loggedUser)

	request := createUnfollowUserRequest("anyuser")
	response, err := rc.App.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestUnFollowUserRoute_bad_request(t *testing.T) {
	rc := setupRoutes()

	request := createUnfollowUserRequest("anyuser")
	response, err := rc.App.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	result := decodeError(t, response)
	assert.NotNil(t, result)
}

func createGetPostsByUserRequest(username string) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/profiles/"+username+"/posts", nil)
}

func TestGetPostsByUser(t *testing.T) {
	rc := setupRoutes()
	createUser(t)

	request := createGetPostsByUserRequest("anyuser")
	response, err := rc.App.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
