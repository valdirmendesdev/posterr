package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	posts_service "github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/handlers"
	"github.com/valdirmendesdev/posterr/pkg/infrastructure/http/fiber/presenters"
	posts_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

var (
	postsRepo posts.Repository
)

func setupPostsRoutes() *handlers.PostsRoutesConfig {
	postsRepo = posts_infra.NewMemoryRepository()
	config := handlers.NewPostsRoutesConfig(fiber.New(), postsRepo)
	handlers.MountPostsRoutes(config)
	return config
}

func createGetPostsRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/posts/", nil)
}

func createPost(t *testing.T) *posts.Post {
	u := utils.GetLoggedUser()
	post, err := posts.NewPost(types.NewUUID(), u, "my post", time.Now())
	assert.NoError(t, err)
	err = postsRepo.Insert(post)
	assert.NoError(t, err)
	return post
}

func decodeAllPostsResponse(t *testing.T, response *http.Response) *[]presenters.Post {
	posts := []presenters.Post{}
	err := json.NewDecoder(response.Body).Decode(&posts)
	assert.NoError(t, err)
	return &posts
}

func TestGetAllPosts(t *testing.T) {
	rc := setupPostsRoutes()
	p := createPost(t)
	request := createGetPostsRequest()
	response, err := rc.App.Test(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	posts := decodeAllPostsResponse(t, response)
	assert.Equal(t, 1, len(*posts))
	assert.Equal(t, p.ID, (*posts)[0].ID)
	assert.Equal(t, p.User.Username.String(), (*posts)[0].Username)
	assert.Equal(t, p.Content, (*posts)[0].Content)
	assert.Equal(t, p.CreatedAt.UTC(), (*posts)[0].CreatedAt.UTC())
}

func createNewPostRequest(t *testing.T, user *users.User, content string) *http.Request {

	bodyBytes, err := json.Marshal(presenters.CreatePost{
		Content: content,
	})
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/posts/", bytes.NewReader(bodyBytes))
	request.Header.Add("Content-Type", "application/json")
	return request
}

func TestCreateNewPost(t *testing.T) {
	rc := setupPostsRoutes()
	lu := utils.GetLoggedUser()
	request := createNewPostRequest(t, lu, "content")
	response, err := rc.App.Test(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestCreatePostOverTheDayLimit(t *testing.T) {
	rc := setupPostsRoutes()
	lu := utils.GetLoggedUser()
	for i := 0; i < posts_service.PostsLimitByDay; i++ {
		createPost(t)
	}
	request := createNewPostRequest(t, lu, "content")
	response, err := rc.App.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	responseError := decodeError(t, response)
	assert.Equal(t, posts_service.ErrDailyPostsLimitReached.Error(), responseError.Message)
}
