package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/profiles"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
)

type ProfileRoutesConfig struct {
	App      *fiber.App
	UserRepo users.Repository
}

func NewProfileRoutesConfigs(app *fiber.App, userRepo users.Repository) *ProfileRoutesConfig {
	return &ProfileRoutesConfig{
		App:      app,
		UserRepo: userRepo,
	}
}

func getProfile(cfg *ProfileRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		username := c.Params("username")
		s := profiles.NewGetByUsernameService(cfg.UserRepo)

		request := profiles.GetByUsernameRequest{
			Username: username,
		}
		user, err := s.Perform(request)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(err)
		}
		return c.JSON(presenters.Profile{
			ID:       user.ID,
			Username: user.Username.String(),
			JoinedAt: user.JoinedAt,
		})
	})
}

func MountProfilesRoutes(cfg *ProfileRoutesConfig) {
	g := cfg.App.Group("/profiles")
	g.Get("/:username", getProfile(cfg))
}
