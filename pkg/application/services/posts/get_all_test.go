package posts_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	posts_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	posts_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
)

var (
	postRepo *posts_infra.MemoryRepository
)

func createUser(t *testing.T) *users.User {
	u, err := users.NewUser(types.NewUUID(), "myuser", time.Now())
	assert.NoError(t, err)
	return u
}

func createPost(t *testing.T, user *users.User) *posts_domain.Post {
	p, err := posts_domain.NewPost(types.NewUUID(), user, "content", time.Now())
	assert.NoError(t, err)
	return p
}

func setupTest(t *testing.T) {
	postRepo = posts_infra.NewMemoryRepository()
}

func TestGetAllPosts(t *testing.T) {
	setupTest(t)
	u := createUser(t)
	postRepo.Insert(createPost(t, u))
	s := posts.NewGetAllPostsService(postRepo)

	response, err := s.Perform()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Posts))
}
