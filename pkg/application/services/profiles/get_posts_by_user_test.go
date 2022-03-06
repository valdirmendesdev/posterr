package profiles_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/domain/users"
	"github.com/valdirmendesdev/posterr/pkg/application/services/profiles"
	posts_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/posts"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/types"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

var (
	postRepo *posts_infra.MemoryRepository
)

func setupGetPostsByUserRepos() {
	userRepo = users_infra.NewMemoryRepository()
	postRepo = posts_infra.NewMemoryRepository()
}

func setupTestGetPostsByUser() *profiles.GetPostsByUserService {
	setupGetPostsByUserRepos()
	return profiles.NewGetPostsByUserService(userRepo, postRepo)
}

func createPost(t *testing.T, user *users.User) *posts.Post {
	p, err := posts.NewPost(types.NewUUID(), user, "content", time.Now())
	assert.NoError(t, err)
	err = postRepo.Insert(p)
	assert.NoError(t, err)
	return p
}

func TestGetPostsByUserID(t *testing.T) {
	s := setupTestGetPostsByUser()
	lu := utils.GetLoggedUser()
	err := userRepo.Add(lu)

	createPost(t, lu)

	assert.NoError(t, err)

	result, err := s.Perform(profiles.GetPostsByUserRequest{
		Username: lu.Username.String(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Posts))
}
