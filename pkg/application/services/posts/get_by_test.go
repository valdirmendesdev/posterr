package posts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

func setupTestGetById(t *testing.T) *posts.GetByIDService {
	setupPostsTest(t)
	return posts.NewGetByIDService(postRepo)
}

func TestGetPostByID(t *testing.T) {
	s := setupTestGetById(t)
	p := createPost(t, utils.GetLoggedUser())
	err := postRepo.Insert(p)
	assert.NoError(t, err)

	result, err := s.Perform(posts.GetByIDRequest{
		ID: p.ID.String(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, p, result.Post)
}
