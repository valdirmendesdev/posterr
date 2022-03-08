package posts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	posts_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	posts_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/posts"
	users_infra "github.com/valdirmendesdev/posterr/pkg/infrastructure/repositories/users"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

func setupRepos() {
	userRepo = users_infra.NewMemoryRepository()
	postRepo = posts_infra.NewMemoryRepository()
}

func setupTestRepost(t *testing.T) *posts.RepostService {
	setupRepos()
	return posts.NewRepostService(userRepo, postRepo)
}

func TestRepost(t *testing.T) {
	s := setupTestRepost(t)
	lu := utils.GetLoggedUser()
	err := userRepo.Add(lu)
	assert.NoError(t, err)

	p := createPost(t, lu)
	err = postRepo.Insert(p)
	assert.NoError(t, err)

	result, err := s.Perform(posts.RepostRequest{
		PostID:   p.ID.String(),
		Username: lu.Username.String(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, p, result.Post.ReferencedPost)
	assert.Equal(t, posts_domain.TypeRepost, result.Post.Type)
}
