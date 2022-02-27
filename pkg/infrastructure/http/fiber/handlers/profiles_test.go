package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/handlers"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	shared_presenters "github.com/valdirmendesdev/posterr/pkg/shared/infrastructure/http/presenters"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var memoryRepository *users_infra.MemoryRepository

func setupRoutes() *handlers.ProfileRoutesConfig {
	memoryRepository = users_infra.NewMemoryRepository()
	config := handlers.NewProfileRoutesConfigs(fiber.New(), memoryRepository)
	handlers.MountProfilesRoutes(config)
	return config
}

func createGetProfileRequest(username string) *http.Request {
	return httptest.NewRequest(http.MethodGet, "/profiles/"+username, nil)
}

func TestGetProfileRoute(t *testing.T) {
	routesConfig := setupRoutes()

	user, err := users.NewUser(
		types.NewUUID(),
		users.Username("myusername"),
		time.Now(),
	)
	assert.NoError(t, err)
	memoryRepository.Add(user)

	expected := &presenters.Profile{
		ID:       user.ID,
		Username: user.Username.String(),
		JoinedAt: user.JoinedAt,
	}

	request := createGetProfileRequest("myusername")
	response, err := routesConfig.App.Test(request)
	var result *presenters.Profile
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Username, result.Username)
	assert.Equal(t, expected.JoinedAt.UTC().GoString(), result.JoinedAt.UTC().GoString())
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
	var result *shared_presenters.Error
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.NoError(t, err)
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
	var result *shared_presenters.Error
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, result)
}
