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
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var memoryRepository *users_infra.MemoryRepository

func setupRoutes() *handlers.ProfileRoutesConfig {
	memoryRepository = users_infra.NewMemoryRepository()
	config := &handlers.ProfileRoutesConfig{
		App:      fiber.New(),
		UserRepo: memoryRepository,
	}
	handlers.MountProfilesRoutes(config)
	return config
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

	request := httptest.NewRequest(http.MethodGet, "/profiles/myusername", nil)
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
