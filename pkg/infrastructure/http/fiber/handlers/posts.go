package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	posts_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	shared_presenters "github.com/valdirmendesdev/posterr/pkg/shared/infrastructure/http/presenters"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

type PostsRoutesConfig struct {
	App       *fiber.App
	UsersRepo users.Repository
	PostsRepo posts_domain.Repository
}

const (
	postIDParam = "postID"
)

func NewPostsRoutesConfig(app *fiber.App, userRepo users.Repository, postsRepo posts_domain.Repository) *PostsRoutesConfig {
	return &PostsRoutesConfig{
		App:       app,
		PostsRepo: postsRepo,
		UsersRepo: userRepo,
	}
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
			p, err := presenters.TransformPostToPresenter(post)
			if err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}
			posts = append(posts, p)
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

		p, err := presenters.TransformPostToPresenter(serviceResponse.Post)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Status(http.StatusCreated).JSON(p)
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

		p, err := presenters.TransformPostToPresenter(serviceResponse.Post)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.JSON(p)
	})
}

func repost(cfg *PostsRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		s := posts.NewRepostService(cfg.UsersRepo, cfg.PostsRepo)
		lu := utils.GetLoggedUser()

		_, err := s.Perform(posts.RepostRequest{
			PostID:   c.Params(postIDParam),
			Username: lu.Username.String(),
		})
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		return c.SendStatus(http.StatusNoContent)
	})
}

func quote(cfg *PostsRoutesConfig) fiber.Handler {
	return fiber.Handler(func(c *fiber.Ctx) error {
		postBody := new(presenters.CreatePost)

		if err := c.BodyParser(postBody); err != nil {
			return c.Status(http.StatusBadRequest).JSON(&shared_presenters.Error{
				Message: err.Error(),
			})
		}

		s := posts.NewQuoteService(cfg.UsersRepo, cfg.PostsRepo)
		lu := utils.GetLoggedUser()

		_, err := s.Perform(posts.QuoteRequest{
			PostID:   c.Params(postIDParam),
			Username: lu.Username.String(),
			Content:  postBody.Content,
		})
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		return c.SendStatus(http.StatusNoContent)
	})
}

func MountPostsRoutes(cfg *PostsRoutesConfig) {
	g := cfg.App.Group("/posts")
	g.Get("/", getPosts(cfg))
	g.Post("/", createPost(cfg))
	g.Get(fmt.Sprintf("/:%s", postIDParam), getPostID(cfg))
	g.Put(fmt.Sprintf("/:%s/repost", postIDParam), repost(cfg))
	g.Put(fmt.Sprintf("/:%s/quote", postIDParam), quote(cfg))
}
