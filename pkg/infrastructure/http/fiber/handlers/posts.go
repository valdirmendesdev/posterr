package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	posts_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
)

type PostsRoutesConfig struct {
	App       *fiber.App
	PostsRepo posts_domain.Repository
}

func NewPostsRoutesConfig(app *fiber.App, postsRepo posts_domain.Repository) *PostsRoutesConfig {
	return &PostsRoutesConfig{
		App:       app,
		PostsRepo: postsRepo,
	}
}

func MountPostsRoutes(cfg *PostsRoutesConfig) {
	g := cfg.App.Group("/posts")
	g.Get("/", getPosts(cfg))
}

func getPosts(cfg *PostsRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		s := posts.NewGetAllPostsService(cfg.PostsRepo)

		serviceResponse, err := s.Perform()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		posts := []presenters.Post{}
		for _, post := range serviceResponse.Posts {
			posts = append(posts, presenters.Post{
				ID:        post.ID,
				Username:  post.User.Username.String(),
				Content:   post.Content,
				CreatedAt: post.CreatedAt,
			})
		}

		return c.JSON(posts)
	})
}
