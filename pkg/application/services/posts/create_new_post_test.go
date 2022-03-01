package posts_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valdirmendesdev/posterr/pkg/application/services/posts"
)

func TestCreateNewPost(t *testing.T) {
	setupPostsTest(t)
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

func TestCreatePostOverTheDayLimit(t *testing.T) {
	setupPostsTest(t)
	u := createUser(t)
	s := posts.NewCreateNewPostService(postRepo)

	for i := 0; i < posts.PostsLimitByDay; i++ {
		_, err := s.Perform(posts.CreateNewPostRequest{
			User:    u,
			Content: fmt.Sprintf("content %d", i),
		})
		assert.NoError(t, err)
	}
	response, err := s.Perform(posts.CreateNewPostRequest{
		User:    u,
		Content: "content",
	})
	assert.Nil(t, response)
	assert.ErrorIs(t, err, posts.ErrDailyPostsLimitReached)
}