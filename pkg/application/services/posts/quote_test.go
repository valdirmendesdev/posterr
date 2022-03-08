package posts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	posts_domain "github.com/valdirmendesdev/posterr/pkg/application/domain/posts"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
	"github.com/valdirmendesdev/posterr/pkg/shared/utils"
)

func setupTestQuote(t *testing.T) *posts.QuoteService {
	setupRepos()
	return posts.NewQuoteService(userRepo, postRepo)
}

func TestQuote(t *testing.T) {
	s := setupTestQuote(t)
	lu := utils.GetLoggedUser()
	err := userRepo.Add(lu)
	assert.NoError(t, err)

	p := createPost(t, lu)
	err = postRepo.Insert(p)
	assert.NoError(t, err)

	result, err := s.Perform(posts.QuoteRequest{
		PostID:   p.ID.String(),
		Username: lu.Username.String(),
		Content:  "test",
	})

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, p, result.Post.ReferencedPost)
	assert.Equal(t, posts_domain.TypeQuote, result.Post.Type)
	assert.Equal(t, "test", result.Post.Content)
}
