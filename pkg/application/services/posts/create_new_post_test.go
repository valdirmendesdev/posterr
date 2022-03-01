package posts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
)

func TestCreateNewPost(t *testing.T) {
	setupTest(t)
	u := createUser(t)

	s := posts.NewCreateNewPostService(postRepo)
	response, err := s.Perform(posts.CreateNewPostRequest{
		User:    u,
		Content: "content",
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.ID)
	assert.Equal(t, u.Username.String(), response.Username)
	assert.Equal(t, "content", response.Content)
	assert.NotNil(t, response.CreatedAt)
}
