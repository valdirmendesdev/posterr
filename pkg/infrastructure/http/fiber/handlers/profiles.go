package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/friendships"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
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
	PostsRepo      posts.Repository
}

const (
	usernameParamName = "username"
)

func NewProfileRoutesConfigs(app *fiber.App, userRepo users.Repository, friendshipRepo friendships.Repository, postRepo posts.Repository) *ProfileRoutesConfig {
	return &ProfileRoutesConfig{
		App:            app,
		UserRepo:       userRepo,
		FriendshipRepo: friendshipRepo,
		PostsRepo:      postRepo,
	}
}

func getProfile(cfg *ProfileRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		username := c.Params(usernameParamName)
		s := profiles.NewGetByUsernameService(cfg.UserRepo, cfg.FriendshipRepo)

		request := profiles.GetByUsernameRequest{
			Username: username,
		}
		profile, err := s.Perform(request)
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

		cfs := profiles.NewCheckFriendshipService(cfg.UserRepo, cfg.FriendshipRepo)
		lu := utils.GetLoggedUser()

		friendship, _ := cfs.Perform(profiles.CheckFriendshipRequest{
			Username:         username,
			FollowerUsername: lu.Username.String(),
		})

		isFollowing := false
		if friendship.Friendship != nil {
			isFollowing = true
		}

		return c.JSON(presenters.Profile{
			ID:          profile.ID,
			Username:    profile.Username.String(),
			JoinedAt:    profile.JoinedAt,
			Followers:   profile.Followers,
			Following:   profile.Following,
			IsFollowing: isFollowing,
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
			switch err {
			case friendships.ErrAlreadyExists:
				return c.SendStatus(http.StatusNotModified)

			default:
				return c.Status(http.StatusBadRequest).JSON(&shared_presenters.Error{
					Message: err.Error(),
				})
			}
		}

		return c.SendStatus(http.StatusNoContent)
	})
}

func unfollowUser(cfg *ProfileRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		username := c.Params(usernameParamName)
		s := profiles.NewUnfollowUserService(cfg.UserRepo, cfg.FriendshipRepo)

		lu := utils.GetLoggedUser()

		request := profiles.UnfollowUserRequest{
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

func getPostsByUser(cfg *ProfileRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		username := c.Params(usernameParamName)
		s := profiles.NewGetPostsByUserService(cfg.UserRepo, cfg.PostsRepo)

		request := profiles.GetPostsByUserRequest{
			Username: username,
		}
		result, err := s.Perform(request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&shared_presenters.Error{
				Message: err.Error(),
			})
		}

		posts := []*presenters.Post{}
		for _, post := range result.Posts {
			posts = append(posts, &presenters.Post{
				ID:        post.ID,
				Username:  post.User.Username.String(),
				Content:   post.Content,
				CreatedAt: post.CreatedAt,
			})
		}

		return c.JSON(posts)
	})
}

func MountProfilesRoutes(cfg *ProfileRoutesConfig) {
	g := cfg.App.Group("/profiles")
	g.Get("/:username", getProfile(cfg))
	g.Put("/:username/follow", followUser(cfg))
	g.Delete("/:username/unfollow", unfollowUser(cfg))
	g.Get("/:username/posts", getPostsByUser(cfg))
}
