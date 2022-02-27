package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/profiles"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	shared_presenters "github.com/valdirmendesdev/posterr/pkg/shared/infrastructure/http/presenters"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

type ProfileRoutesConfig struct {
	App            *fiber.App
	UserRepo       users.Repository
	FriendshipRepo friendships.Repository
}

const (
	usernameParamName = "username"
)

func NewProfileRoutesConfigs(app *fiber.App, userRepo users.Repository, friendshipRepo friendships.Repository) *ProfileRoutesConfig {
	return &ProfileRoutesConfig{
		App:            app,
		UserRepo:       userRepo,
		FriendshipRepo: friendshipRepo,
	}
}

func getProfile(cfg *ProfileRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		username := c.Params(usernameParamName)
		s := profiles.NewGetByUsernameService(cfg.UserRepo)

		request := profiles.GetByUsernameRequest{
			Username: username,
		}
		user, err := s.Perform(request)
		if err != nil {
			var errorStatusCode int
			switch err {
			case users.ErrNotFound:
				errorStatusCode = http.StatusNotFound
			default:
				errorStatusCode = http.StatusBadRequest
			}

			return c.Status(errorStatusCode).JSON(&shared_presenters.Error{
				Message: err.Error(),
			})
		}

		return c.JSON(presenters.Profile{
			ID:       user.ID,
			Username: user.Username.String(),
			JoinedAt: user.JoinedAt,
		})
	})
}

func followUser(cfg *ProfileRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		username := c.Params(usernameParamName)
		s := profiles.NewFollowUserService(cfg.UserRepo, cfg.FriendshipRepo)

		lu := utils.GetLoggedUser()

		request := profiles.FollowUserRequest{
			UserUsername:     username,
			FollowerUsername: lu.Username.String(),
		}
		err := s.Perform(request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&shared_presenters.Error{
				Message: err.Error(),
			})
		}

		return c.SendStatus(http.StatusNoContent)
	})
}

func MountProfilesRoutes(cfg *ProfileRoutesConfig) {
	g := cfg.App.Group("/profiles")
	g.Get("/:username", getProfile(cfg))
	g.Put("/:username/follow", followUser(cfg))
}
