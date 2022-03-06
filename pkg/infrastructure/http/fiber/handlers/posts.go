package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	posts_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	shared_presenters "github.com/valdirmendesdev/posterr/pkg/shared/infrastructure/http/presenters"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

type PostsRoutesConfig struct {
	App       *fiber.App
	PostsRepo posts_domain.Repository
}

const (
	postIDParam = "postID"
)

func NewPostsRoutesConfig(app *fiber.App, postsRepo posts_domain.Repository) *PostsRoutesConfig {
	return &PostsRoutesConfig{
		App:       app,
		PostsRepo: postsRepo,
	}
}

func MountPostsRoutes(cfg *PostsRoutesConfig) {
	g := cfg.App.Group("/posts")
	g.Get("/", getPosts(cfg))
	g.Post("/", createPost(cfg))
	g.Get(fmt.Sprintf("/:%s", postIDParam), getPostID(cfg))
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

func createPost(cfg *PostsRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		postBody := new(presenters.CreatePost)

		if err := c.BodyParser(postBody); err != nil {
			return c.Status(http.StatusBadRequest).JSON(&shared_presenters.Error{
				Message: err.Error(),
			})
		}

		s := posts.NewCreateNewPostService(cfg.PostsRepo)
		lu := utils.GetLoggedUser()
		serviceResponse, err := s.Perform(posts.CreateNewPostRequest{
			User:    lu,
			Content: postBody.Content,
		})
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(&shared_presenters.Error{
				Message: err.Error(),
			})
		}

		return c.Status(http.StatusCreated).JSON(presenters.Post{
			ID:        serviceResponse.ID,
			Username:  lu.Username.String(),
			Content:   serviceResponse.Content,
			CreatedAt: serviceResponse.CreatedAt,
		})
	})
}

func getPostID(cfg *PostsRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		s := posts.NewGetByIDService(cfg.PostsRepo)

		serviceResponse, err := s.Perform(posts.GetByIDRequest{
			ID: c.Params(postIDParam),
		})
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		return c.JSON(presenters.Post{
			ID:        serviceResponse.Post.ID,
			Username:  serviceResponse.Post.User.Username.String(),
			Content:   serviceResponse.Post.Content,
			CreatedAt: serviceResponse.Post.CreatedAt,
		})
	})
}
